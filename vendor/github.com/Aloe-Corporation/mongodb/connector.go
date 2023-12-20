package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Conf contains all information to connect to a MongoDB server.
type Conf struct {
	DB         string `yaml:"db"`             // Name of the database.
	Host       string `yaml:"host"`           // URL to reach the mongoDB server.
	Port       int    `yaml:"port,omitempty"` // Optionnal port, if set to 0 it won't be processed.
	Username   string `yaml:"username"`       // Credential to authenticate to the db.
	Password   string `yaml:"password"`       // Credential to authenticate to the db.
	AuthSource string `yaml:"auth_source"`    // Database to check authentication
	Timeout    int    `yaml:"timeout"`        // Connection timeout in seconds
}

// Connector is the connector used to communicate with MongoDB database server.
// It embeds a native mongo.Client so it can be used as is and is supercharged with
// additionnal methods.
type Connector struct {
	*mongo.Client
	DB          string
	Collections map[string]*mongo.Collection
}

// Collection returns the  *mongo.Collection identified its name.
// If the specified collections doesn't exists on con.Collections map
// then add it.
func (con *Connector) Collection(collectionName string) *mongo.Collection {
	if con.Collections == nil {
		con.Collections = make(map[string]*mongo.Collection)
	}

	if _, ok := con.Collections[collectionName]; !ok {
		con.Collections[collectionName] = con.Database(con.DB).Collection(collectionName)
	}

	return con.Collections[collectionName]
}

// TryConnection tests ping, it end if the ping is a success or timeout.
func (con *Connector) TryConnection() error {
	if err := con.Ping(context.Background(), nil); err != nil {
		return fmt.Errorf("fail to ping mongo: %w", err)
	}

	return nil
}

// FactoryConnector instanciates a new *Connector with the given params.
func FactoryConnector(c Conf) (*Connector, error) {
	connectionURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", c.Username, c.Password, c.Host, c.AuthSource)
	if c.Port != 0 {
		connectionURI = fmt.Sprintf("mongodb://%s:%s@%s/%s?retryWrites=true&w=majority", c.Username, c.Password, c.Host, c.AuthSource)

	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	clientOptions := options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverAPI)
	err := clientOptions.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid client options: %w", err)
	}

	timeout := time.Second * time.Duration(c.Timeout)

	clientOptions.ServerSelectionTimeout = &timeout

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("fail to create mongo client: %w", err)
	}

	con := &Connector{
		DB:          c.DB,
		Client:      client,
		Collections: make(map[string]*mongo.Collection),
	}
	return con, nil
}

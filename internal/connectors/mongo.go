package connectors

import (
	"context"
	"fmt"

	"github.com/Aloe-Corporation/mongodb"
	"go.uber.org/zap"
)

const (
	TypeConnectorMongo = "MongoDB"
)

var Mongo = make(map[string]*mongodb.Connector)

func initAllConnectorMongo() error {
	for key, config := range Config.Mongo {
		if err := initConnectorMongo(key, config); err != nil {
			return fmt.Errorf("fail to init all MongoDB connectors: %w", err)
		}
	}

	return nil
}

// initConnectorMongo init a mongo connector in Mongo map
func initConnectorMongo(key string, config mongodb.Conf) error {
	var err error

	log.Info("Init MongoDB connector " + key + "...")
	connector, err := mongodb.FactoryConnector(config)
	if err != nil {
		return fmt.Errorf("fail to init MongoDB connector %s: %w", key, err)
	}

	log.Info("Try connection MongoDB " + key + "...")
	err = connector.TryConnection()
	if err != nil {
		return fmt.Errorf("fail to ping MongoDB %s: %w", key, err)
	}

	log.Info("MongoDB connector " + key + " is ready to use")
	Mongo[key] = connector

	return err
}

func closeMongo() {
	for key, connector := range Mongo {
		if err := connector.Disconnect(context.Background()); err != nil {
			log.Error("fail to disconnect mongo connector",
				zap.String("key", key),
				zap.Error(err),
			)
		}
	}
}

func GetConnectorMongo(connectorName string) (*mongodb.Connector, error) {
	connector, exist := Mongo[connectorName]
	if !exist {
		return nil, &ConnectorNotFoundError{
			ConnectorType: TypeConnectorMongo,
			ConnectorName: connectorName,
		}
	}

	return connector, nil
}

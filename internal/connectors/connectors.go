package connectors

import (
	"github.com/Aloe-Corporation/logs"
	"github.com/Aloe-Corporation/mongodb"
	"github.com/Aloe-Corporation/sqldb"
)

var (
	log = logs.Get()

	// Config of the repositories package.
	Config Conf
)

// Conf for the repositories package.
type Conf struct {
	Mongo    map[string]mongodb.Conf `mapstructure:"mongo"`
	Postgres map[string]sqldb.Conf   `mapstructure:"postgres"`
}

// Init the data sources connectors.
func Init() error {
	// Mongo connectors
	log.Info("Init all Mongo connector...")
	if err := initAllConnectorMongo(); err != nil {
		return err
	}
	log.Info("All Mongo  connector is ready to use")

	// Postgres connectors
	log.Info("Init all Postgres connector...")
	if err := initAllConnectorPostgres(); err != nil {
		return err
	}
	log.Info("All Postgres  connector is ready to use")

	return nil
}

// Close DAO connectors.
func Close() error {
	closeMongo()
	closePostgres()

	return nil
}

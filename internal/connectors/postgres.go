package connectors

import (
	"fmt"

	"github.com/Aloe-Corporation/sqldb"
	"go.uber.org/zap"
)

const (
	TypeConnectorPostgres = "Postgres"
)

var Postgres = make(map[string]*sqldb.Connector)

func initAllConnectorPostgres() error {
	for key, config := range Config.Postgres {
		if err := initConnectorPostgres(key, config); err != nil {
			return fmt.Errorf("fail to init all Postgres connectors: %w", err)
		}
	}

	return nil
}

// initConnectorPostgres init a postgres connector in Postgres map
func initConnectorPostgres(key string, config sqldb.Conf) error {
	var err error

	log.Info("Init Postgres connector " + key + "...")
	connector, err := sqldb.FactoryConnector(config)
	if err != nil {
		return fmt.Errorf("fail to init Postgres connector %s: %w", key, err)
	}

	log.Info("Try connection Postgres " + key + "...")
	err = connector.TryConnection(10)
	if err != nil {
		return fmt.Errorf("fail to ping Postgres %s: %w", key, err)
	}

	log.Info("Postgres connector " + key + " is ready to use")
	Postgres[key] = connector

	return err
}

func closePostgres() {
	for key, connector := range Postgres {
		err := connector.Close()
		if err != nil {
			log.Error("fail to disconnect postgres connector",
				zap.String("key", key),
				zap.Error(err),
			)
		}
	}
}

func GetConnectorPostgres(connectorName string) (*sqldb.Connector, error) {
	connector, exist := Postgres[connectorName]
	if !exist {
		return nil, &ConnectorNotFoundError{
			ConnectorType: TypeConnectorPostgres,
			ConnectorName: connectorName,
		}
	}

	return connector, nil
}

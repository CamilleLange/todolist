package connectors

import (
	"github.com/Aloe-Corporation/sqldb"

	"github.com/Aloe-Corporation/logs"
)

var (
	log = logs.Get()

	// Config of the repositories package.
	Config Conf
)

// Conf for the repositories package.
type Conf struct {
	MySQL map[string]sqldb.Conf `mapstructure:"mysql"`
}

// Init the data sources connectors.
func Init() error {

	return nil
}

// Close DAO connectors.
func Close() error {

	return nil
}

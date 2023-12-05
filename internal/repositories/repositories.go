package repositories

import (
	"github.com/Aloe-Corporation/logs"
)

var (
	log = logs.Get()
)

// DAOFactoryOptions is the generic struct used by all DAO factory method patterns to build specific DAO.
type DAOFactoryOptions struct {
	Type      string `mapstructure:"type"`
	Connector string `mapstructure:"connector"`
}

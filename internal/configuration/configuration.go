package configuration

import (
	"fmt"
	"strings"

	"github.com/Aloe-Corporation/logs"
	"github.com/CamilleLange/todolist/internal/connectors"
	"github.com/CamilleLange/todolist/internal/controllers"
	"github.com/CamilleLange/todolist/internal/ginrouters"
	"github.com/spf13/viper"
)

var (
	log = logs.Get()
	// Config hold all config data for the all API
	Config *Conf
)

// Conf hold all config data
type Conf struct {
	// Modules
	Logger *logs.Conf `mapstructure:"logger"`

	// Internal Packages
	Connectors  *connectors.Conf  `mapstructure:"connectors"`
	Controllers *controllers.Conf `mapstructure:"controllers"`
	GinRouters  *ginrouters.Conf  `mapstructure:"ginrouters"`
}

// LoadConf load the configuration from the file at the given path.
func LoadConf(path, prefix string) error {
	Config = new(Conf)

	Config.Logger = &logs.Config

	Config.Connectors = &connectors.Config
	Config.Controllers = &controllers.Config
	Config.GinRouters = &ginrouters.Config

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("can't load API configuration : %w", err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("can't unmarshall API configuration : %w", err)
	}

	return nil

}

// InitAllModules is use for Init all modules.
func InitAllModules() error {
	log.Info("init logs modules...")
	err := logs.Init()
	if err != nil {
		return fmt.Errorf("fail to init logs modules: %w", err)
	}
	log.Info("logs ready")
	return nil
}

// InitAllPkg is use for Init all internal package.
func InitAllPkg() error {
	log.Info("init connectors package...")
	err := connectors.Init()
	if err != nil {
		return fmt.Errorf("fail to init connectors package: %w", err)
	}
	log.Info("connectors ready")

	log.Info("init controllers package...")
	err = controllers.Init()
	if err != nil {
		return fmt.Errorf("fail to init controllers package: %w", err)
	}
	log.Info("controllers ready")

	log.Info("init routers package...")
	ginrouters.Init()
	log.Info("routers ready")

	return nil
}

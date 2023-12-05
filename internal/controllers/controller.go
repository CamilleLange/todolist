package controllers

import (
	"fmt"

	"github.com/Aloe-Corporation/logs"
)

var (
	log = logs.Get()
	// Config of the controllers package
	Config Conf

	// Add your intance below

	// TaskInstance is an instance of ITaskController.
	TaskInstance ITaskController
)

// Conf for the controllers package
type Conf struct {
	TaskController TaskControllerConf `mapstructure:"task_controller"`
}

// Init the controllerss
func Init() error {
	var err error

	log.Info("init TaskController...")
	TaskInstance, err = factoryTaskController(Config.TaskController)
	if err != nil {
		return fmt.Errorf("fail to build TaskController: %w", err)
	}
	log.Info("TaskController is ready to use")

	log.Info("controllers package ready")
	return err
}

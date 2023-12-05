package controllers

import (
	"context"
	"fmt"

	model "github.com/CamilleLange/todolist/pkg/structs"

	"github.com/CamilleLange/todolist/internal/repositories"
)

var (
	_ ITaskController = (*TaskController)(nil)
)

// ITaskController is an interface for TaskController and TaskControllerMocking.
type ITaskController interface {
	Create(ctx context.Context) (*model.TaskPublicDTO, error)
	Get(ctx context.Context) (*model.Task, error)
	GetAll(ctx context.Context) ([]model.TaskPublicDTO, error)
	Update(ctx context.Context) error
	Delete(ctx context.Context) error
}

// TaskControllerConf is a configuration structure for TaskController.
type TaskControllerConf struct {
	TaskDAO repositories.DAOFactoryOptions `mapstructure:"task_dao"`
}

// TaskController is an controllers to manage business logic of Task.
type TaskController struct {
	daoTask repositories.ITaskDAO
}

func (c *TaskController) Create(ctx context.Context) (*model.TaskPublicDTO, error) {
	task, err := c.daoTask.Create(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to create task: %w", err)
	}

	return model.FactoryTaskPublicDTO(task), nil
}

func (c *TaskController) Get(ctx context.Context) (*model.Task, error) {
	task, err := c.daoTask.ReadByUUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get task: %w", err)
	}

	return task, nil
}

func (c *TaskController) GetAll(ctx context.Context) ([]model.TaskPublicDTO, error) {
	tasks, err := c.daoTask.ReadAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get tasks: %w", err)
	}

	publicTasks := []model.TaskPublicDTO{}
	for _, f := range tasks {
		publicTasks = append(publicTasks, *model.FactoryTaskPublicDTO(f))
	}

	return publicTasks, nil
}

func (c *TaskController) Update(ctx context.Context) error {
	err := c.daoTask.Update(ctx)
	if err != nil {
		return fmt.Errorf("fail to update task: %w", err)
	}

	return nil
}

func (c *TaskController) Delete(ctx context.Context) error {
	err := c.daoTask.Delete(ctx)
	if err != nil {
		return fmt.Errorf("fail to delete task: %w", err)
	}

	return nil
}

// factoryTaskController is use to build an TaskController according to the conf.
func factoryTaskController(c TaskControllerConf) (*TaskController, error) {
	log.Info("loading TaskDAO...")
	daoTask, err := repositories.ProxyFactoryTaskDAO(c.TaskDAO)
	if err != nil {
		return nil, fmt.Errorf("fail to load TaskDAO: %w", err)
	}
	log.Info("TaskDAO loaded")

	controllers := &TaskController{
		daoTask: daoTask,
	}
	return controllers, nil
}

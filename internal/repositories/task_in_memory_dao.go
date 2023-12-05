package repositories

import (
	"context"
	"fmt"
	"time"

	model "github.com/CamilleLange/todolist/pkg/structs"
	"github.com/google/uuid"
)

const (
	// TypeTaskInMemoryDAO is an identifier to build TaskInMemoryDAO.
	TypeTaskInMemoryDAO = "TaskInMemoryDAO"
)

var _ ITaskDAO = (*TaskInMemoryDAO)(nil)

// TaskInMemoryDAO is a TaskDAO with not implemented features.
type TaskInMemoryDAO struct {
	connectorName string

	tasks map[uuid.UUID]*model.Task
}

func (dao *TaskInMemoryDAO) Create(ctx context.Context) (*model.Task, error) {
	strTask := ctx.Value("create_task")
	taskToCreate, castable := strTask.(*model.TaskCreateDTO)
	if !castable {
		return nil, fmt.Errorf("can't cast the context value to *model.TaskCreateDTO")
	}

	task := taskToCreate.ReverseCreateDTO()
	task.UUID = uuid.New()
	task.CreatedAt = time.Now()
	task.LastUpdated = task.CreatedAt

	dao.tasks[task.UUID] = task
	return task, nil
}

func (dao *TaskInMemoryDAO) ReadByUUID(ctx context.Context) (*model.Task, error) {
	strUUID := ctx.Value("task_uuid")
	taskUUID, castable := strUUID.(*uuid.UUID)
	if !castable {
		return nil, fmt.Errorf("can't cast the context value to *uuid.UUID")
	}

	task, exist := dao.tasks[*taskUUID]
	if !exist {
		return nil, fmt.Errorf("no task with this UUID (%s) exist", taskUUID.String())
	}
	return task, nil
}

func (dao *TaskInMemoryDAO) ReadAll(ctx context.Context) ([]*model.Task, error) {
	tasks := make([]*model.Task, 0)
	for _, task := range dao.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (dao *TaskInMemoryDAO) Update(ctx context.Context) error {
	strFieldsName := ctx.Value("task_fields_name")
	fieldsName, castable := strFieldsName.([]string)
	if !castable {
		return fmt.Errorf("can't cast the context value to *[]string")
	}

	strValuesToUpdate := ctx.Value("task_values_to_update")
	valuesToUpdate, castable := strValuesToUpdate.(map[string]any)
	if !castable {
		return fmt.Errorf("can't cast the context value to *map[string]any")
	}

	strUUID := ctx.Value("task_uuid")
	taskUUID, castable := strUUID.(*uuid.UUID)
	if !castable {
		return fmt.Errorf("can't cast the context value to *uuid.UUID")
	}

	taskToUpdate, exist := dao.tasks[*taskUUID]
	if !exist {
		return fmt.Errorf("no task with this UUID (%s) exist", taskUUID.String())
	}

	for _, field := range fieldsName {
		value, castable := valuesToUpdate[field].(string)
		if !castable {
			return fmt.Errorf("value for field %s not castable to string", field)
		}

		switch field {
		case "description":
			taskToUpdate.WhatToDo = value

		case "status":
			taskToUpdate.Status = value

		default:
			return fmt.Errorf("unknown field name %s", field)
		}
	}

	taskToUpdate.LastUpdated = time.Now()
	return nil
}

func (dao *TaskInMemoryDAO) Delete(ctx context.Context) error {
	strUUID := ctx.Value("task_uuid")
	taskUUID, castable := strUUID.(*uuid.UUID)
	if !castable {
		return fmt.Errorf("can't cast the context value to *uuid.UUID")
	}

	if _, exist := dao.tasks[*taskUUID]; !exist {
		return fmt.Errorf("no task with this UUID (%s) exist", taskUUID.String())
	}

	delete(dao.tasks, *taskUUID)
	return nil
}

// factoryTaskInMemoryDAO build TaskInMemoryDAO.
func factoryTaskInMemoryDAO(opt DAOFactoryOptions) (*TaskInMemoryDAO, error) {
	return &TaskInMemoryDAO{
		connectorName: opt.Connector,
		tasks:         make(map[uuid.UUID]*model.Task),
	}, nil
}

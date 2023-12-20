package repositories

import (
	"context"
	"fmt"

	model "github.com/CamilleLange/todolist/pkg/structs"
	"github.com/google/uuid"
)

// mapTaskDAO is used by ProxyFactoryTaskDAO to store TaskDAO.
var mapTaskDAO = make(map[string]map[string]ITaskDAO)

// ITaskDAO is a DAO interface to manage Task.
type ITaskDAO interface {
	Create(ctx context.Context) (*model.Task, error)
	ReadByUUID(ctx context.Context) (*model.Task, error)
	ReadAll(ctx context.Context) ([]*model.Task, error)
	Update(ctx context.Context) error
	Delete(ctx context.Context) error
}

// ProxyFactoryTaskDAO uses FactoryTaskDAO if the TaskDAO don't exist, and returns TaskDAO.
func ProxyFactoryTaskDAO(opt DAOFactoryOptions) (ITaskDAO, error) {
	// Test if exist
	mapConnector, mapExist := mapTaskDAO[opt.Type]
	if mapExist {
		daoTask, present := mapConnector[opt.Connector]
		if present {
			return daoTask, nil
		}
	}

	// Build new TaskDAO
	daoTask, err := FactoryTaskDAO(opt)
	if err != nil {
		return nil, fmt.Errorf("fail to build new TaskDAO: %w", err)
	}

	// Save new TaskDAO
	if !mapExist {
		mapTaskDAO[opt.Type] = make(map[string]ITaskDAO)
	}
	mapTaskDAO[opt.Type][opt.Connector] = daoTask

	return daoTask, nil
}

// FactoryTaskDAO builds a new TaskDAO according to the typename.
func FactoryTaskDAO(opt DAOFactoryOptions) (ITaskDAO, error) {
	var dao ITaskDAO
	var err error

	switch opt.Type {
	case TypeTaskVoidDAO:
		dao, err = factoryTaskVoidDAO(opt)
	case TypeTaskInMemoryDAO:
		dao, err = factoryTaskInMemoryDAO(opt)
	case TypeTaskPostgresDAO:
		dao, err = factoryTaskPostgresDAO(opt)
	case TypeTaskMongoDAO:
		dao, err = factoryTaskMongoDAO(opt)
	default:
		return nil, &DAOTypeNotFoundError{Type: opt.Type}
	}

	if err != nil {
		return nil, fmt.Errorf("fail to build %v: %w", opt.Type, err)
	}

	return dao, nil
}

// getUpdateDataFrom search in the context and convert to the correct type all data needed to execute the update query
func getUpdateDataFrom(ctx context.Context) ([]string, map[string]any, *uuid.UUID, error) {
	strFieldsName := ctx.Value("task_fields_name")
	fieldsName, castable := strFieldsName.([]string)
	if !castable {
		return nil, nil, nil, fmt.Errorf("can't cast the context value to *[]string")
	}

	strValuesToUpdate := ctx.Value("task_values_to_update")
	valuesToUpdate, castable := strValuesToUpdate.(map[string]any)
	if !castable {
		return nil, nil, nil, fmt.Errorf("can't cast the context value to *map[string]any")
	}

	strUUID := ctx.Value("task_uuid")
	taskUUID, castable := strUUID.(*uuid.UUID)
	if !castable {
		return nil, nil, nil, fmt.Errorf("can't cast the context value to *uuid.UUID")
	}
	return fieldsName, valuesToUpdate, taskUUID, nil
}

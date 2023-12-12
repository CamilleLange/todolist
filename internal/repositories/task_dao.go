package repositories

import (
	"context"
	"fmt"

	model "github.com/CamilleLange/todolist/pkg/structs"
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
	default:
		return nil, &DAOTypeNotFoundError{Type: opt.Type}
	}

	if err != nil {
		return nil, fmt.Errorf("fail to build %v: %w", opt.Type, err)
	}

	return dao, nil
}

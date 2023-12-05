package repositories

import (
	"context"
	model "github.com/CamilleLange/todolist/pkg/structs"
)

const (
	// TypeTaskVoidDAO is an identifier to build TaskVoidDAO.
	TypeTaskVoidDAO = "TaskVoidDAO"
)

var _ ITaskDAO = (*TaskVoidDAO)(nil)

// TaskVoidDAO is a TaskDAO with not implemented features.
type TaskVoidDAO struct {
	connectorName string
}


func (dao *TaskVoidDAO) Create(ctx context.Context) (*model.Task, error) {
	return nil, ErrFeatureNotImplemented
}

func (dao *TaskVoidDAO) ReadByUUID(ctx context.Context) (*model.Task, error) {
	return nil, ErrFeatureNotImplemented
}

func (dao *TaskVoidDAO) ReadAll(ctx context.Context) ([]*model.Task, error) {
	return nil, ErrFeatureNotImplemented
}

func (dao *TaskVoidDAO) Update(ctx context.Context) error {
	return ErrFeatureNotImplemented
}

func (dao *TaskVoidDAO) Delete(ctx context.Context) error {
	return ErrFeatureNotImplemented
}

// factoryTaskVoidDAO build TaskVoidDAO.
func factoryTaskVoidDAO(opt DAOFactoryOptions) (*TaskVoidDAO, error) {
	return &TaskVoidDAO{
		connectorName: opt.Connector,
	}, nil
}

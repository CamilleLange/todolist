package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Aloe-Corporation/sqldb"
	"github.com/CamilleLange/todolist/internal/connectors"
	model "github.com/CamilleLange/todolist/pkg/structs"
	"github.com/google/uuid"
)

const (
	// TypeTaskPostgresDAO is an identifier to build TaskPostgresDAO.
	TypeTaskPostgresDAO = "TaskPostgresDAO"
)

var _ ITaskDAO = (*TaskPostgresDAO)(nil)

// TaskPostgresDAO is a TaskDAO with not implemented features.
type TaskPostgresDAO struct {
	connector     *sqldb.Connector
	connectorName string
}

func (dao *TaskPostgresDAO) Create(ctx context.Context) (*model.Task, error) {
	// Find the task data.
	strTask := ctx.Value("create_task")
	taskToCreate, castable := strTask.(*model.TaskCreateDTO)
	if !castable {
		return nil, fmt.Errorf("can't cast the context value to *model.TaskCreateDTO")
	}

	var (
		taskUUID               uuid.UUID
		createdAt, lastUpdated time.Time
	)

	// Insert the task data and query the default value generated by the database.
	query := "INSERT INTO tasks (description, status) VALUES ($1, $2) RETURNING task_uuid, created_at, last_updated"
	if err := dao.connector.QueryRow(query, taskToCreate.WhatToDo, taskToCreate.Status).Scan(
		&taskUUID,
		&createdAt,
		&lastUpdated,
	); err != nil {
		return nil, fmt.Errorf("can't scan the auto generated data : %w", err)
	}

	// Set the task last data and return it.
	task := taskToCreate.ReverseCreateDTO()
	task.UUID = taskUUID
	task.CreatedAt = createdAt
	task.LastUpdated = lastUpdated

	return task, nil
}

func (dao *TaskPostgresDAO) ReadByUUID(ctx context.Context) (*model.Task, error) {
	// Find task UUID.
	strUUID := ctx.Value("task_uuid")
	taskUUID, castable := strUUID.(*uuid.UUID)
	if !castable {
		return nil, fmt.Errorf("can't cast the context value to *uuid.UUID")
	}

	// Query the database with the task UUID.
	task := new(model.Task)
	query := "SELECT task_uuid, description, status, created_at, last_updated FROM tasks WHERE task_uuid = $1;"
	if err := dao.connector.QueryRow(query, taskUUID).Scan(
		&task.UUID,
		&task.WhatToDo,
		&task.Status,
		&task.CreatedAt,
		&task.LastUpdated,
	); err != nil {
		return nil, fmt.Errorf("can't query by UUID : %w", err)
	}
	return task, nil
}

func (dao *TaskPostgresDAO) ReadAll(ctx context.Context) ([]*model.Task, error) {
	// Query the database for all tasks.
	query := "SELECT task_uuid, description, status, created_at, last_updated FROM tasks;"
	rows, err := dao.connector.Query(query)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("can't query all tasks : %w", err)
	}

	// Scan the result to retrieve each task.
	tasks := make([]*model.Task, 0)
	for rows.Next() {
		task := new(model.Task)

		if err := rows.Scan(
			&task.UUID,
			&task.WhatToDo,
			&task.Status,
			&task.CreatedAt,
			&task.LastUpdated,
		); err != nil {
			return nil, fmt.Errorf("can't scan row : %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (dao *TaskPostgresDAO) Update(ctx context.Context) error {
	// Find the update data.
	fieldsName, valuesToUpdate, taskUUID, err := getUpdateDataFrom(ctx)
	if err != nil {
		return fmt.Errorf("can't get data for the update request from the context : %w", err)
	}

	// Create the query string.
	var (
		setStatements []string
		params        []any
	)

	for _, field := range fieldsName {
		setStatements = append(setStatements, fmt.Sprintf("%s = $%d", field, len(params)+1))
		params = append(params, valuesToUpdate[field])
	}

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE task_uuid = $%d", strings.Join(setStatements, ", "), len(params)+1)
	params = append(params, taskUUID)

	// Open a transaction.
	tx, err := dao.connector.Begin()
	if err != nil {
		return fmt.Errorf("can't begin the transaction :%w", err)
	}

	// Execute the update query.
	result, err := dao.connector.Exec(tx, query, params...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("can't rollback the tx : %w", err)
		}
		return fmt.Errorf("can't update the task : %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("can't rollback the tx : %w", err)
		}
		return fmt.Errorf("can't get the number of affected rows : %w", err)
	}
	if rowsAffected != 1 {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("can't rollback the tx : %w", err)
		}
		return fmt.Errorf("we update more than one row")
	}

	// Commit the transaction.
	if err := dao.connector.Commit(tx); err != nil {
		return fmt.Errorf("can't commit the transaction : %w", err)
	}
	return nil
}

func (dao *TaskPostgresDAO) Delete(ctx context.Context) error {
	// Find task UUID.
	strUUID := ctx.Value("task_uuid")
	taskUUID, castable := strUUID.(*uuid.UUID)
	if !castable {
		return fmt.Errorf("can't cast the context value to *uuid.UUID")
	}

	// Open a transaction.
	tx, err := dao.connector.Begin()
	if err != nil {
		return fmt.Errorf("can't begin the transaction :%w", err)
	}

	// Query the database to delete the task.
	query := "DELETE FROM tasks WHERE task_uuid = $1;"
	result, err := dao.connector.Exec(tx, query, taskUUID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("can't rollback the tx : %w", err)
		}
		return fmt.Errorf("can't delete the task : %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("can't rollback the tx : %w", err)
		}
		return fmt.Errorf("can't get the number of affected rows : %w", err)
	}
	if rowsAffected != 1 {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("can't rollback the tx : %w", err)
		}
		return fmt.Errorf("we delete more than one row")
	}

	// Commit the transaction.
	if err := dao.connector.Commit(tx); err != nil {
		return fmt.Errorf("can't commit the transaction : %w", err)
	}
	return nil
}

// factoryTaskPostgresDAO build TaskPostgresDAO.
func factoryTaskPostgresDAO(opt DAOFactoryOptions) (*TaskPostgresDAO, error) {
	connector, err := connectors.GetConnectorPostgres(opt.Connector)
	if err != nil {
		return nil, fmt.Errorf("fail to get connector: %w", err)
	}

	return &TaskPostgresDAO{
		connector:     connector,
		connectorName: opt.Connector,
	}, nil
}

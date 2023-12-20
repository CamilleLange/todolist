package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/Aloe-Corporation/mongodb"
	"github.com/CamilleLange/todolist/internal/connectors"
	model "github.com/CamilleLange/todolist/pkg/structs"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// TypeTaskMongoDAO is an identifier to build TaskMongoDAO.
	TypeTaskMongoDAO = "TaskMongoDAO"

	collectionName = "tasks"
)

var _ ITaskDAO = (*TaskMongoDAO)(nil)

// TaskMongoDAO is a TaskDAO with not implemented features.
type TaskMongoDAO struct {
	connector     *mongodb.Connector
	connectorName string
}

func (dao *TaskMongoDAO) Create(ctx context.Context) (*model.Task, error) {
	// Find the task data.
	strTask := ctx.Value("create_task")
	taskToCreate, castable := strTask.(*model.TaskCreateDTO)
	if !castable {
		return nil, fmt.Errorf("can't cast the context value to *model.TaskCreateDTO")
	}

	// generate values to insert the new task.
	var (
		taskUUID    uuid.UUID = uuid.New()
		createdAt   time.Time = time.Now().UTC()
		lastUpdated           = createdAt
	)
	task := taskToCreate.ReverseCreateDTO()
	task.UUID = taskUUID
	task.CreatedAt = createdAt
	task.LastUpdated = lastUpdated

	if _, err := dao.connector.Collection(collectionName).InsertOne(ctx, task); err != nil {
		return nil, fmt.Errorf("can't insert the new task : %w", err)
	}

	return task, nil
}

func (dao *TaskMongoDAO) ReadByUUID(ctx context.Context) (*model.Task, error) {
	// Find task UUID.
	strUUID := ctx.Value("task_uuid")
	taskUUID, castable := strUUID.(*uuid.UUID)
	if !castable {
		return nil, fmt.Errorf("can't cast the context value to *uuid.UUID")
	}

	task := new(model.Task)
	filter := bson.M{"task_uuid": taskUUID}
	err := dao.connector.Collection(collectionName).FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		return nil, fmt.Errorf("can't decode task : %w", err)
	}
	return task, nil
}

func (dao *TaskMongoDAO) ReadAll(ctx context.Context) ([]*model.Task, error) {
	cursor, err := dao.connector.Collection(collectionName).Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("can't read all tasks : %w", err)
	}
	defer cursor.Close(context.Background())

	// Scan the result to retrieve each task.
	tasks := make([]*model.Task, 0)
	for cursor.Next(ctx) {
		task := new(model.Task)

		if err := cursor.Decode(&task); err != nil {
			return nil, fmt.Errorf("can't decode task : %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (dao *TaskMongoDAO) Update(ctx context.Context) error {
	// Find the update data.
	fieldsName, valuesToUpdate, taskUUID, err := getUpdateDataFrom(ctx)
	if err != nil {
		return fmt.Errorf("can't get data for the update request from the context : %w", err)
	}

	filter := bson.M{"task_uuid": taskUUID}
	updateFields := bson.M{}
	for _, fieldName := range fieldsName {
		if value, ok := valuesToUpdate[fieldName]; ok {
			updateFields[fieldName] = value
		}
	}
	update := bson.M{"$set": updateFields}

	if _, err := dao.connector.Collection(collectionName).UpdateOne(context.Background(), filter, update); err != nil {
		return fmt.Errorf("can't update the task : %w", err)
	}
	return nil
}

func (dao *TaskMongoDAO) Delete(ctx context.Context) error {
	// Find task UUID.
	strUUID := ctx.Value("task_uuid")
	taskUUID, castable := strUUID.(*uuid.UUID)
	if !castable {
		return fmt.Errorf("can't cast the context value to *uuid.UUID")
	}

	filter := bson.M{"task_uuid": taskUUID}
	if _, err := dao.connector.Collection(collectionName).DeleteOne(context.Background(), filter); err != nil {
		return fmt.Errorf("can't delete the task : %w", err)
	}
	return nil
}

// factoryTaskMongoDAO build TaskMongoDAO.
func factoryTaskMongoDAO(opt DAOFactoryOptions) (*TaskMongoDAO, error) {
	connector, err := connectors.GetConnectorMongo(opt.Connector)
	if err != nil {
		return nil, fmt.Errorf("fail to get connector: %w", err)
	}

	return &TaskMongoDAO{
		connector:     connector,
		connectorName: opt.Connector,
	}, nil
}

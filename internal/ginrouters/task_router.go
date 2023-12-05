package ginrouters

import (
	"context"
	"net/http"
	"sync"

	"github.com/CamilleLange/todolist/internal/controllers"
	model "github.com/CamilleLange/todolist/pkg/structs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	ginparamsmapper "gitlab.com/Zandraz/gin-params-mapper"
	"go.uber.org/zap"
)

var (
	onceInitTaskRouter sync.Once
	// singletonTaskRouter is a singleton instance of TaskRouter.
	singletonTaskRouter *TaskRouter
)

// TaskRouter groups a set of handlers to manage entrypoints of Task.
type TaskRouter struct {
	ctlTask controllers.ITaskController
}

func (r *TaskRouter) Post(c *gin.Context) {
	task := new(model.TaskCreateDTO)
	if err := c.ShouldBindJSON(task); err != nil {
		log.Error("TaskRouter.Post fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := ValidateInstance.Struct(task); err != nil {
		log.Error("TaskRouter.Post fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.WithValue(c, "create_task", task)

	createdTask, err := r.ctlTask.Create(ctx)
	if err != nil {
		log.Error("TaskRouter.Post fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusOK, createdTask)
}

func (r *TaskRouter) GetAll(c *gin.Context) {
	ctx := context.WithValue(c, "", "")

	tasks, err := r.ctlTask.GetAll(ctx)
	if err != nil {
		log.Error("TaskRouter.GetAll fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (r *TaskRouter) Get(c *gin.Context) {
	var taskUUID uuid.UUID
	if err := ginparamsmapper.GetPathParamFromContext("task_uuid", c, &taskUUID); err != nil {
		log.Error("TaskRouter.Get fail to get path param task_uuid: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	ctx := context.WithValue(c, "task_uuid", &taskUUID)

	task, err := r.ctlTask.Get(ctx)
	if err != nil {
		log.Error("TaskController.GetByUUID fail",
			zap.Any("task_uuid", taskUUID),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusOK, task)
}

func (r *TaskRouter) Put(c *gin.Context) {
	var taskUUID uuid.UUID
	if err := ginparamsmapper.GetPathParamFromContext("task_uuid", c, &taskUUID); err != nil {
		log.Error("TaskRouter.Get fail to get path param task_uuid: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	taskFieldsToUpdate := map[string]any{}
	if err := c.ShouldBindJSON(&taskFieldsToUpdate); err != nil {
		log.Error("TaskRouter.Put fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "bad request")
		return
	}

	if len(taskFieldsToUpdate) == 0 {
		c.JSON(http.StatusNoContent, "no fields to update")
		return
	}

	taskUpdateDTO := new(model.TaskUpdateDTO)
	config := &mapstructure.DecoderConfig{
		ErrorUnused: true, // Extra fields throw an error.
		Result:      &taskUpdateDTO,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		log.Error("TaskRouter.Put fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Internal server error")
		return
	}

	if err := decoder.Decode(taskFieldsToUpdate); err != nil {
		log.Error("TaskRouter.Put fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Can't decode request body")
		return
	}

	fieldsNameToUpdate := make([]string, 0)
	for fieldName := range taskFieldsToUpdate {
		fieldsNameToUpdate = append(fieldsNameToUpdate, fieldName)
	}

	// Only validate the fields that appears in request body
	if err := ValidateInstance.StructPartial(taskUpdateDTO, fieldsNameToUpdate...); err != nil {
		log.Error("TaskRouter.Put fail : invalid request body fields: %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Can't decode request body")
		return
	}

	ctx := context.WithValue(c, "task_fields_name", fieldsNameToUpdate)
	ctx = context.WithValue(ctx, "task_values_to_update", taskFieldsToUpdate)
	ctx = context.WithValue(ctx, "task_uuid", &taskUUID)

	err = r.ctlTask.Update(ctx)
	if err != nil {
		log.Error("TaskRouter.Put fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "can't update the resource")
		return
	}

	c.JSON(http.StatusNoContent, "Task updated.")
}

func (r *TaskRouter) Delete(c *gin.Context) {
	var taskUUID uuid.UUID

	if err := ginparamsmapper.GetPathParamFromContext("task_uuid", c, &taskUUID); err != nil {
		log.Error("TaskRouter.Delete fail to get path param task_uuid: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	ctx := context.WithValue(c, "task_uuid", &taskUUID)

	err := r.ctlTask.Delete(ctx)
	if err != nil {
		log.Error("TaskRouter.Delete fail",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	c.JSON(http.StatusNoContent, "Task deleted.")
}

// GetInstanceTaskRouter get singleton instance of TaskRouter.
func GetInstanceTaskRouter() *TaskRouter {
	if singletonTaskRouter == nil {
		onceInitTaskRouter.Do(
			func() {
				singletonTaskRouter = &TaskRouter{
					ctlTask: controllers.TaskInstance,
				}
			},
		)
	}

	return singletonTaskRouter
}

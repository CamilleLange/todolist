package structs

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	UUID        uuid.UUID
	WhatToDo    string
	Status      string
	CreatedAt   time.Time
	LastUpdated time.Time
}

type TaskPublicDTO struct {
	UUID        uuid.UUID `json:"task_uuid" mapstructure:"task_uuid" binding:"required"`
	WhatToDo    string    `json:"description" mapstructure:"description" binding:"required"`
	Status      string    `json:"status" mapstructure:"status" binding:"required"`
	CreatedAt   time.Time `json:"created_at" mapstructure:"created_at" binding:"required"`
	LastUpdated time.Time `json:"last_updated" mapstructure:"last_updated" binding:"required"`
}

func (dto *TaskPublicDTO) ReversePublicDTO() *Task {
	return &Task{
		UUID:        dto.UUID,
		WhatToDo:    dto.WhatToDo,
		Status:      dto.Status,
		CreatedAt:   dto.CreatedAt,
		LastUpdated: dto.LastUpdated,
	}
}

func FactoryTaskPublicDTO(task *Task) *TaskPublicDTO {
	return &TaskPublicDTO{
		UUID:        task.UUID,
		WhatToDo:    task.WhatToDo,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		LastUpdated: task.LastUpdated,
	}
}

type TaskCreateDTO struct {
	WhatToDo string `json:"description" mapstructure:"description" binding:"required"`
	Status   string `json:"status" mapstructure:"status" binding:"required"`
}

func (dto *TaskCreateDTO) ReverseCreateDTO() *Task {
	return &Task{
		UUID:     uuid.New(),
		WhatToDo: dto.WhatToDo,
		Status:   dto.Status,
	}
}

func FactoryTaskCreateDTO(task *Task) *TaskCreateDTO {
	return &TaskCreateDTO{
		WhatToDo: task.WhatToDo,
		Status:   task.Status,
	}
}

type TaskUpdateDTO struct {
	WhatToDo string `json:"description" mapstructure:"description" binding:"required"`
	Status   string `json:"status" mapstructure:"status" binding:"required"`
}

func (dto *TaskUpdateDTO) ReverseUpdateDTO() *Task {
	return &Task{
		WhatToDo: dto.WhatToDo,
		Status:   dto.Status,
	}
}

func FactoryTaskUpdateDTO(task *Task) *TaskUpdateDTO {
	return &TaskUpdateDTO{
		WhatToDo: task.WhatToDo,
		Status:   task.Status,
	}
}

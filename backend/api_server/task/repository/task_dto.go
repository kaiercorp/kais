package repository

import (
	"api_server/ent"
	"time"
)

type TaskDTO struct {
	ID           int                    `json:"id,omitempty"`
	ProjectID    *int                   `json:"project_id,omitempty"`
	DatasetID    *int                   `json:"dataset_id,omitempty"`
	Title        string                 `json:"title,omitempty"`
	Description  string                 `json:"description,omitempty"`
	EngineType   string                 `json:"engine_type,omitempty"`
	TargetMetric string                 `json:"target_metric,omitempty"`
	UserParams   map[string]interface{} `json:"user_params,omitempty"`
	Params       []string               `json:"params,omitempty"`
	CreatedAt    time.Time              `json:"created_at,omitempty"`
	UpdatedAt    time.Time              `json:"updated_at,omitempty"`

	Modelings []*ModelingDTO `json:"modelings,omitempty"`
}

type TaskPages struct {
	Tasks     []TaskDTO `json:"tasks"`
	TotalPage int       `json:"totalPage"`
}

func ConvertTaskEntToDTO(entity *ent.Task) *TaskDTO {
	return &TaskDTO{
		ID:           entity.ID,
		ProjectID:    &entity.ProjectID,
		DatasetID:    &entity.DatasetID,
		Title:        entity.Title,
		Description:  entity.Description,
		EngineType:   entity.EngineType,
		TargetMetric: entity.TargetMetric,
		Params:       entity.Params,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
		Modelings:    ConvertModelingEntsToDTOs(entity.Edges.Modelings),
	}
}

func ConvertTaskEntsToDTOs(ents []*ent.Task) []*TaskDTO {
	dtos := []*TaskDTO{}

	for _, entity := range ents {
		dtos = append(dtos, ConvertTaskEntToDTO(entity))
	}

	return dtos
}

func convertTaskEntToTaskListItem(entity *ent.Task) TaskDTO {
	return TaskDTO{
		ID:           entity.ID,
		ProjectID:    &entity.ProjectID,
		DatasetID:    &entity.DatasetID,
		Title:        entity.Title,
		Description:  entity.Description,
		EngineType:   entity.EngineType,
		TargetMetric: entity.TargetMetric,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
		Params:       entity.Params,
	}
}

func ConvertTaskEntsToTaskPages(ents []*ent.Task, pageCount int) *TaskPages {
	dtos := []TaskDTO{}

	for _, entity := range ents {
		dtos = append(dtos, convertTaskEntToTaskListItem(entity))
	}

	return &TaskPages{
		Tasks:     dtos,
		TotalPage: pageCount,
	}
}

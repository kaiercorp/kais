package repository

import (
	"api_server/ent"
	"database/sql"
	"time"
)

type ModelingDTO struct {
	ID            int       `json:"id"`
	LocalID       int       `json:"local_id"`
	TaskID        int       `json:"task_id"`
	ParentID      int       `json:"parent_id"`
	ParentLocalID int       `json:"parent_local_id"`
	DatasetID     int       `json:"datset_id"`
	Params        []string  `json:"params"`
	DatasetStat   []string  `json:"dataset_stat"`
	ModelingType  string    `json:"modeling_type"`
	ModelingStep  string    `json:"modeling_step"`
	Performance   []string  `json:"performance"`
	Progress      float64   `json:"progress"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	StartedAt     time.Time `json:"started_at"`

	ModelingModels []*ModelingModels  `json:"modeling_models"`
	Scores         map[string]float64 `json:"scores"`
	InfTime        float64            `json:"inf_time"`
}

type ModelingModels struct {
	ID         int       `json:"id"`
	DataType   string    `json:"data_type"`
	Data       string    `json:"data"`
	ModelingID int       `json:"modeling_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type ModelingDB struct {
	ID            int                `json:"id"`
	LocalID       int                `json:"local_id"`
	TaskID        int                `json:"task_id"`
	ParentID      int                `json:"parent_id"`
	ParentLocalID int                `json:"parent_local_id"`
	DatasetID     int                `json:"datset_id"`
	Params        sql.NullString     `json:"params"`
	DatasetStat   sql.NullString     `json:"dataset_stat"`
	ModelingType  string             `json:"modeling_type"`
	ModelingStep  string             `json:"modeling_step"`
	Performance   sql.NullString     `json:"performance"`
	Progress      float64            `json:"progress"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	StartedAt     time.Time          `json:"started_at"`
	Scores        map[string]float64 `json:"scores"`
	InfTime       float64            `json:"inf_time"`
}

func ConvertModelingEntToDTO(entity *ent.Modeling) *ModelingDTO {
	return &ModelingDTO{
		ID:             entity.ID,
		LocalID:        entity.LocalID,
		TaskID:         entity.TaskID,
		ParentID:       entity.ParentID,
		ParentLocalID:  entity.ParentLocalID,
		DatasetID:      entity.DatasetID,
		Params:         entity.Params,
		DatasetStat:    entity.DatasetStat,
		ModelingType:   entity.ModelingType,
		ModelingStep:   entity.ModelingStep,
		Performance:    entity.Performance,
		Progress:       entity.Progress,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		ModelingModels: ConvertModelingModelsEntsToDTOs(entity.Edges.ModelingModels),
	}
}

func ConvertModelingEntsToDTOs(ents []*ent.Modeling) []*ModelingDTO {
	dtos := []*ModelingDTO{}

	for _, entity := range ents {
		dtos = append(dtos, ConvertModelingEntToDTO(entity))
	}

	return dtos
}

func ConvertModelingModelsEntsToDTOs(ents []*ent.ModelingModels) []*ModelingModels {
	dtos := []*ModelingModels{}

	for _, entity := range ents {
		dtos = append(dtos, ConvertModelingModelEntToDTO(entity))
	}

	return dtos
}

func ConvertModelingModelEntToDTO(entity *ent.ModelingModels) *ModelingModels {
	return &ModelingModels{
		ID:         entity.ID,
		DataType:   entity.DataType,
		Data:       entity.Data,
		ModelingID: entity.ModelingID,
		CreatedAt:  entity.CreatedAt,
	}
}

type EngineParams struct {
	ModelingID    int      `json:"modeling_id"`
	MultiLabel    bool     `json:"multi_label"`
	MultiNode     bool     `json:"multi_node"`
	DataPath      string   `json:"data_path"`
	SavePath      string   `json:"save_path"`
	ImgHeght      string   `json:"img_height"`
	ImgWidth      string   `json:"img_width"`
	TargetMetric  string   `json:"target_metric"`
	DeviceIDs     []int    `json:"device_ids"`
	GPUAuto       bool     `json:"gpu_auto"`
	EngineType    string   `json:"engine_type"`
	OriginID      int      `json:"origin_id"`
	OriginPath    string   `json:"origin_path"`
	IndexColumn   string   `json:"index_column"`
	OutputColumns []string `json:"output_columns"`
	InputColumns  []string `json:"input_columns"`
	Mode          string   `json:"mode"`
}

type EvaluationDTO struct {
	TaskID       int                    `json:"task_id,omitempty"`
	ParentID     int                    `json:"parent_id,omitempty"`
	DatasetID    int                    `json:"dataset_id,omitempty"`
	EngineType   string                 `json:"engine_type"`
	TargetMetric string                 `json:"target_metric,omitempty"`
	IsPath       bool                   `json:"is_path,omitempty"` // is true, blind test
	UserParams   map[string]interface{} `json:"user_params,omitempty"`
}

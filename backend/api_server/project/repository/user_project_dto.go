package repository

import (
	"time"

	"api_server/ent"
)

type UserProjectDTO struct {
	ID        int        `json:"id,omitempty"`
	Username  string     `json:"username,omitempty"`
	ProjectId int        `json:"project_id,omitempty"`
	IsUse     *bool      `json:"is_use,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func ConvertUserProjectEntToDTO(entity *ent.UserProject) *UserProjectDTO {
	return &UserProjectDTO{
		ID:        entity.ID,
		Username:  entity.Username,
		ProjectId: entity.ProjectID,
		IsUse:     &entity.IsUse,
		CreatedAt: &entity.CreatedAt,
		UpdatedAt: &entity.UpdatedAt,
	}
}

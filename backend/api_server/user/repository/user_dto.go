package repository

import (
	"api_server/ent"
	"time"
)

type UserDTO struct {
	ID       int     `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Group    *int    `json:"group,omitempty"`
	IsUse    *bool   `json:"is_use,omitempty"`
	//Token     *string    `json:"token,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	LoginAt   *time.Time `json:"login_at,omitempty"`
}

func ConvertUserEntToDTO(entity *ent.User) *UserDTO {
	return &UserDTO{
		ID:       entity.ID,
		Name:     &entity.Name,
		Username: &entity.Username,
		Password: &entity.Password,
		Group:    &entity.Group,
		IsUse:    &entity.IsUse,
		//Token:     &entity.Token,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: &entity.UpdatedAt,
	}
}

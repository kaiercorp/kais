package repository

import (
	"time"

	"api_server/ent"
)

type ProjectDTO struct {
	ID          int        `json:"id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	IsUse       *bool      `json:"is_use,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type ProjectListItem struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
}

type ProjectPages struct {
	Projects  []ProjectListItem `json:"projects"`
	TotalPage int               `json:"totalPage"`
	HasMore   bool              `json:"hasMore"`
	NextPage  int               `json:"nextPage"`
}

func ConvertEntToDTO(entity *ent.Project) *ProjectDTO {
	return &ProjectDTO{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		IsUse:       &entity.IsUse,
		CreatedAt:   &entity.CreatedAt,
		UpdatedAt:   &entity.UpdatedAt,
	}
}

func ConvertEntsToProjectPages(ents []*ent.Project, pageCount int, hasMore bool, nextPage int) *ProjectPages {
	projectListItem := []ProjectListItem{}

	for _, ent := range ents {
		projectListItem = append(projectListItem, convertEntToProjectListItem(ent))
	}

	return &ProjectPages{
		Projects:  projectListItem,
		TotalPage: pageCount,
		HasMore:   hasMore,
		NextPage:  nextPage,
	}
}

func convertEntToProjectListItem(entity *ent.Project) ProjectListItem {
	return ProjectListItem{
		Id:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		CreatedAt:   &entity.CreatedAt,
	}
}

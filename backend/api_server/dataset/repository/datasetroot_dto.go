package repository

import "api_server/ent"

type DatasetRootDTO struct {
	ID       int           `json:"id,omitempty"`
	Name     string        `json:"name,omitempty"`
	Path     string        `json:"path,omitempty"`
	IsUse    bool          `json:"is_use,omitempty"`
	Datasets []*DatasetDTO `json:"datasets,omitempty"`
}

func ConvertDatasetrootEntToDTO(entity *ent.DatasetRoot) *DatasetRootDTO {
	return &DatasetRootDTO{
		ID:       entity.ID,
		Name:     entity.Name,
		Path:     entity.Path,
		IsUse:    entity.IsUse,
		Datasets: ConvertDatasetEntsToDTOs(entity.Edges.Datasets),
	}
}

func ConvertDatasetrootEntsToDTOs(ents []*ent.DatasetRoot) []*DatasetRootDTO {
	drDTOs := []*DatasetRootDTO{}

	for _, dr := range ents {
		drDTOs = append(drDTOs, ConvertDatasetrootEntToDTO(dr))
	}

	return drDTOs
}

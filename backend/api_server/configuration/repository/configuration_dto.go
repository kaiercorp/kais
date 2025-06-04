package repository

import (
	"api_server/ent"
	"time"
)

type ConfigDTO struct {
	ID         int       `json:"id"`
	ConfigType string    `json:"config_type"`
	ConfigKey  string    `json:"config_key"`
	ConfigVal  string    `json:"config_val"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ConvertEntToDTO(entConfig *ent.Configuration) *ConfigDTO {
	return &ConfigDTO{
		ID:         entConfig.ID,
		ConfigType: entConfig.ConfigType,
		ConfigKey:  entConfig.ConfigKey,
		ConfigVal:  entConfig.ConfigVal,
		CreatedAt:  entConfig.CreatedAt,
		UpdatedAt:  entConfig.UpdatedAt,
	}
}

func ConvertEntsToDTOs(ents []*ent.Configuration) []*ConfigDTO {
	configDTO := []*ConfigDTO{}

	for _, config := range ents {
		configDTO = append(configDTO, &ConfigDTO{
			ID:         config.ID,
			ConfigType: config.ConfigType,
			ConfigKey:  config.ConfigKey,
			ConfigVal:  config.ConfigVal,
			CreatedAt:  config.CreatedAt,
			UpdatedAt:  config.UpdatedAt,
		})
	}

	return configDTO
}

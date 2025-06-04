package repository

import "api_server/ent"

type DeviceDTO struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	IP         string `json:"ip,omitempty"`
	Port       *int   `json:"port,omitempty"`
	IsUse      *bool  `json:"is_use,omitempty"`
	Type       string `json:"type,omitempty"`
	Connection string `json:"connection,omitempty"`
}

type DeviceRemoveDTO struct {
	IDs []int `json:"ids"`
}

func ConvertEntsToDTOs(ents []*ent.Device) []*DeviceDTO {
	dDTOs := []*DeviceDTO{}

	for _, d := range ents {
		dDTOs = append(dDTOs, ConvertEntToDTO(d))
	}

	return dDTOs
}

func ConvertEntToDTO(entity *ent.Device) *DeviceDTO {
	return &DeviceDTO{
		ID:         entity.ID,
		Name:       entity.Name,
		IP:         entity.IP,
		Port:       &entity.Port,
		IsUse:      &entity.IsUse,
		Type:       entity.Type,
		Connection: entity.Connection,
	}
}

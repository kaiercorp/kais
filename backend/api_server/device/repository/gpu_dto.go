package repository

import "api_server/ent"

type DongleGPUStatus struct {
	DongleID        string   `json:"dongle_id"`
	MaxGPUCount     int      `json:"max_gpu_count"`
	CurrentGPUCount int      `json:"current_gpu_count"`
	CurrentGPUIds   []string `json:"current_gpu_ids"`
}
type GPUDTO struct {
	ID       int    `json:"id"`
	UUID     string `json:"uuid,omitempty"`
	Index    int    `json:"index"`
	Name     string `json:"name,omitempty"`
	State    string `json:"state,omitempty"`
	IsUse    bool   `json:"is_use,omitempty"`
	DeviceID int    `json:"device_id"`
}

func ConvertGPUEntsToGPUDTOs(ents []*ent.Gpu) []*GPUDTO {
	dDTOs := make([]*GPUDTO, len(ents))

	for idx, d := range ents {
		dDTOs[idx] = ConvertGPUEntToGPUDTO(d)
	}

	return dDTOs
}

func ConvertGPUEntToGPUDTO(entity *ent.Gpu) *GPUDTO {
	return &GPUDTO{
		ID:       entity.ID,
		UUID:     entity.UUID,
		Index:    entity.Index,
		Name:     entity.Name,
		State:    entity.State,
		IsUse:    entity.IsUse,
		DeviceID: entity.DeviceID,
	}
}

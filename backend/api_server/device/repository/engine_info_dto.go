package repository

import "api_server/ent"

type CPUInfo struct {
	Uilization float64 `json:"utilization_cpu,omitempty"`
}

type DiskInfo struct {
	Free    uint64  `json:"free,omitempty"`
	Used    uint64  `json:"used,omitempty"`
	Total   uint64  `json:"total,omitempty"`
	Percent float64 `json:"percent,omitempty"`
}

type GPUInfo struct {
	ID                int     `json:"id"`
	Index             int     `json:"index"`
	Name              string  `json:"name,omitempty"`
	UUID              string  `json:"uuid,omitempty"`
	MemoryTotal       float64 `json:"memory_total,omitempty"`
	MemoryUsed        float64 `json:"memory_used,omitempty"`
	UtilizationGPU    float64 `json:"utilization_gpu,omitempty"`
	UtilizationMemory float64 `json:"utilization_memory,omitempty"`
}

type EngineInfoDTO struct {
	DeviceID   int       `json:"device_id,omitempty"`
	DeviceName string    `json:"device_name,omitempty"`
	CPU        CPUInfo   `json:"cpu,omitempty"`
	DISK       DiskInfo  `json:"disk,omitempty"`
	GPUs       []GPUInfo `json:"gpu,omitempty"`
}

func ConvertGPUInfoEntToDTO(entity *ent.Gpu) GPUInfo {
	return GPUInfo{
		ID:    entity.ID,
		Index: entity.Index,
		Name:  entity.Name,
		UUID:  entity.UUID,
	}
}

func ConvertGPUInfoEntsToDTOs(ents []*ent.Gpu) []GPUInfo {
	dtos := []GPUInfo{}

	for _, entity := range ents {
		dtos = append(dtos, ConvertGPUInfoEntToDTO(entity))
	}

	return dtos
}

func ConvertEnginInfoEntToDTO(entity *ent.Device) *EngineInfoDTO {
	return &EngineInfoDTO{
		DeviceID:   entity.ID,
		DeviceName: entity.Name,
		GPUs:       ConvertGPUInfoEntsToDTOs(entity.Edges.Gpu),
	}
}

func ConvertEngineInfoEntsToDTOs(ents []*ent.Device) []*EngineInfoDTO {
	dtos := []*EngineInfoDTO{}

	for _, entity := range ents {
		dtos = append(dtos, ConvertEnginInfoEntToDTO(entity))
	}

	return dtos
}

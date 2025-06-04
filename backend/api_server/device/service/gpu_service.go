package service

import (
	"context"
	"fmt"
	"sync"

	repo "api_server/device/repository"
	"api_server/ent"
	"api_server/logger"
)

type IGPUService interface {
	// ReadActive는 사용 중인(활성화된) GPU 디바이스 정보를 조회하여 반환합니다.
	//
	// 매개변수:
	//   - 없음
	//
	// 반환값:
	//   - []*repo.EngineInfoDTO: 활성화된 GPU 엔진 정보 DTO 리스트
	//   - *logger.Report: 오류 발생 시 리포트, 없으면 nil
	ReadActive() ([]*repo.EngineInfoDTO, *logger.Report)

	// ReadManyByDeviceID는 주어진 디바이스 ID에 연결된 GPU 리스트를 조회합니다.
	//
	// 매개변수:
	//   - deviceID: 조회할 GPU 디바이스의 고유 ID
	//
	// 반환값:
	//   - []*repo.GPUDTO: 해당 디바이스에 연결된 GPU DTO 리스트
	//   - *logger.Report: 오류 발생 시 리포트, 없으면 nil
	ReadManyByDeviceID(deviceID int) ([]*repo.GPUDTO, *logger.Report)
	ViewGpuByIndex(index string) (*repo.GPUDTO, error)
}

type GPUService struct {
	ctx context.Context
	dao repo.IGPUDAO
}

var onceGPUService sync.Once
var instanceGPUService *GPUService

func NewGPUService(dao repo.IGPUDAO) *GPUService {
	onceGPUService.Do(func() {
		instanceGPUService = &GPUService{
			ctx: context.Background(),
			dao: dao,
		}
	})

	return instanceGPUService
}

func NewStatic() *GPUService {
	if instanceGPUService == nil {
		return NewGPUService(repo.NewGPUDAO())
	}
	return instanceGPUService
}

// ReadActive는 사용 중인(활성화된) GPU 디바이스 정보를 조회하여 반환합니다.
func (svc *GPUService) ReadActive() ([]*repo.EngineInfoDTO, *logger.Report) {
	if deviceList, err := svc.dao.SelectIsUse(svc.ctx); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertEngineInfoEntsToDTOs(deviceList), nil
	}
}

// ReadManyByDeviceID는 주어진 디바이스 ID에 연결된 GPU 리스트를 조회합니다.
func (svc *GPUService) ReadManyByDeviceID(deviceID int) ([]*repo.GPUDTO, *logger.Report) {
	gpuList, err := svc.dao.SelectManyByDeviceID(svc.ctx, deviceID)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}
	return repo.ConvertGPUEntsToGPUDTOs(gpuList), nil
}

func (svc *GPUService) ViewGpuByIndex(gpuIndex string) (*repo.GPUDTO, error) {
	gpu, err := svc.dao.SelectGpuByIndex(svc.ctx, gpuIndex)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("failed to find GPU with gpu_index=%s and is_use=true", gpuIndex)
		}
		return nil, err
	} else {
		return repo.ConvertGPUEntToGPUDTO(gpu), nil
	}
}

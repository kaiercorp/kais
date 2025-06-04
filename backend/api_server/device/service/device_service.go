package service

import (
	"context"
	"fmt"
	"sync"

	repo "api_server/device/repository"
	"api_server/logger"
)

type IDeviceService interface {
	Create(req repo.DeviceDTO) ([]*repo.DeviceDTO, *logger.Report)
	ReadAll() ([]*repo.DeviceDTO, *logger.Report)
	ReadActive() ([]*repo.DeviceDTO, *logger.Report)
	ReadOne(id int) (*repo.DeviceDTO, *logger.Report)
	Edit(req repo.DeviceDTO) (*repo.DeviceDTO, *logger.Report)
	DeleteMany(ids []int) ([]*repo.DeviceDTO, *logger.Report)
	DeleteOne(device_id int) ([]*repo.DeviceDTO, *logger.Report)
}

type DeviceService struct {
	ctx context.Context
	dao repo.IDeviceDAO
}

var once sync.Once
var instance *DeviceService

func New(dao repo.IDeviceDAO) *DeviceService {
	once.Do(func() {
		instance = &DeviceService{
			ctx: context.Background(),
			dao: dao,
		}
	})

	return instance
}

func (svc *DeviceService) Create(req repo.DeviceDTO) ([]*repo.DeviceDTO, *logger.Report) {
	device, err := svc.dao.SelectByIPAndPort(svc.ctx, req.IP, *req.Port)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}
	if len(device) > 0 {
		return nil, logger.CreateReport(&logger.CODE_DB_DUPLICATE, fmt.Errorf("device already exists: %+v", device))
	}
	_, err = svc.dao.InsertOne(svc.ctx, req)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}
	return svc.ReadAll()
}

func (svc *DeviceService) ReadAll() ([]*repo.DeviceDTO, *logger.Report) {
	if deviceList, err := svc.dao.SelectAll(svc.ctx); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertEntsToDTOs(deviceList), nil
	}
}

func (svc *DeviceService) ReadActive() ([]*repo.DeviceDTO, *logger.Report) {
	if deviceList, err := svc.dao.SelectActive(svc.ctx); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertEntsToDTOs(deviceList), nil
	}
}

func (svc *DeviceService) ReadOne(id int) (*repo.DeviceDTO, *logger.Report) {
	if device, err := svc.dao.SelectOne(svc.ctx, id); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertEntToDTO(device), nil
	}
}

func (svc *DeviceService) Edit(req repo.DeviceDTO) (*repo.DeviceDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("%+v", req))

	exist, err := svc.dao.SelectByIPAndPort(svc.ctx, req.IP, *req.Port)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}
	if len(exist) > 0 {
		if exist[0].ID != req.ID {
			return nil, logger.CreateReport(&logger.CODE_DB_DUPLICATE, fmt.Errorf("device already exists: %+v", exist))
		}
	}

	device, err := svc.dao.UpdateOne(svc.ctx, req)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}
	return repo.ConvertEntToDTO(device), nil
}

func (svc *DeviceService) DeleteMany(ids []int) ([]*repo.DeviceDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf(`%+v`, ids))
	logger.Error("delete many...", ids)
	if _, err := svc.dao.DeleteMany(svc.ctx, ids); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_DELETE, err)
	} else {
		return svc.ReadAll()
	}
}

func (svc *DeviceService) DeleteOne(device_id int) ([]*repo.DeviceDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"device_id": %d}`, device_id))
	if err := svc.dao.DeleteOne(svc.ctx, device_id); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_DELETE, err)
	} else {
		return svc.ReadAll()
	}
}

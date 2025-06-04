package service

import (
	"context"
	"fmt"
	"sync"

	repo "api_server/configuration/repository"
	"api_server/logger"
	"api_server/utils"
)

type IConfigService interface {
	IsIgnoreLicense() bool
	Get(key string) string
	Read(group int) ([]*repo.ConfigDTO, *logger.Report)
	EditMany(configs []repo.ConfigDTO, group int) ([]*repo.ConfigDTO, *logger.Report)
	EditOne(req repo.ConfigDTO) (*repo.ConfigDTO, *logger.Report)
}

type ConfigService struct {
	ctx context.Context
	dao repo.IConfigDAO
}

var once sync.Once
var instance *ConfigService

func New(dao repo.IConfigDAO) *ConfigService {
	once.Do(func() {
		println("Config Service instance")
		instance = &ConfigService{
			ctx: context.Background(),
			dao: dao,
		}
	})

	return instance
}

func NewStatic() *ConfigService {
	if instance == nil {
		return New(repo.New())
	}

	return instance
}

func (svc *ConfigService) IsIgnoreLicense() bool {
	if config, err := svc.dao.SelectOneByKey(svc.ctx, "IGNORE_LICENSE"); err != nil {
		return false
	} else if config == nil {
		return false
	} else if config.ConfigVal == "true" {
		return true
	}

	return false
}

func (svc *ConfigService) Get(key string) string {
	// fmt.Printf("{'key': %s}\n", key)
	if config, err := svc.dao.SelectOneByKey(svc.ctx, key); err != nil {
		return ""
	} else if config == nil {
		return ""
	} else {
		return config.ConfigVal
	}
}

func (svc *ConfigService) Read(group int) ([]*repo.ConfigDTO, *logger.Report) {
	fmt.Printf("{'group': %d}\n", group)
	configs, err := svc.dao.SelectAllByType(svc.ctx, utils.CONFIG_TYPE_USER)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	if group == 0 {
		if configsInit, err := svc.dao.SelectAllByType(svc.ctx, utils.CONFIG_TYPE_SYSTEM); err == nil {
			configs = append(configs, configsInit...)
		}
	}

	return repo.ConvertEntsToDTOs(configs), nil
}

func (svc *ConfigService) EditMany(configs []repo.ConfigDTO, group int) ([]*repo.ConfigDTO, *logger.Report) {
	fmt.Printf("%+v\n", configs)
	if err := svc.dao.UpdateMany(svc.ctx, configs); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	return svc.Read(group)
}

func (svc *ConfigService) EditOne(req repo.ConfigDTO) (*repo.ConfigDTO, *logger.Report) {
	fmt.Printf("%+v\n", req)
	if config, err := svc.dao.UpdateOne(svc.ctx, req); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	} else {
		return repo.ConvertEntToDTO(config), nil
	}
}

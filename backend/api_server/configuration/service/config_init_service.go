package service

import (
	"context"
	"sync"

	repo "api_server/configuration/repository"
)

type IConfigInitService interface {
}

type ConfigInitService struct {
	ctx context.Context
	dao repo.IConfigInitDAO
}

var onceInit sync.Once
var instanceInit *ConfigInitService

func NewInit() *ConfigInitService {
	onceInit.Do(func() {
		println("Config Init Service instance")
		instanceInit = &ConfigInitService{
			ctx: context.Background(),
			dao: repo.NewConfigInit(),
		}
	})

	return instanceInit
}

func (svc *ConfigInitService) Init(version string) {
	svc.dao.InsertDefaultValues(svc.ctx, version)
}

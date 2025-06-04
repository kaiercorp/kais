package router

import (
	"sync"

	"github.com/gin-gonic/gin"

	auth_service "api_server/auth/service"
	repo "api_server/configuration/repository"
	"api_server/configuration/service"
	"api_server/logger"
)

type ConfigController struct {
	svc      service.IConfigService
	auth_svc auth_service.IUserService
}

var once sync.Once
var instance *ConfigController

func New(svc service.IConfigService, auth_svc auth_service.IUserService) *ConfigController {
	once.Do(func() {
		logger.Debug("Config Controller instance")
		instance = &ConfigController{
			svc:      svc,
			auth_svc: auth_svc,
		}
	})

	return instance
}

func (ctlr *ConfigController) GetAll(c *gin.Context) {
	logger.ApiRequest(c)

	group := ctlr.auth_svc.ReadUserGroup(c.GetHeader("Authorization"))

	data, report := ctlr.svc.Read(group)
	logger.ApiResponse(c, report, data)
}

func (ctlr *ConfigController) UpdateMany(c *gin.Context) {
	logger.ApiRequest(c)

	group := ctlr.auth_svc.ReadUserGroup(c.GetHeader("Authorization"))

	configs := []repo.ConfigDTO{}
	if err := c.ShouldBindJSON(&configs); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.EditMany(configs, group)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *ConfigController) UpdateLanguage(c *gin.Context) {
	logger.ApiRequest(c)

	config := repo.ConfigDTO{}
	if err := c.ShouldBindJSON(&config); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.EditOne(config)
		logger.ApiResponse(c, report, data)
	}
}

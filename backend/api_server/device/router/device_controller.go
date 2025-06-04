package router

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	repo "api_server/device/repository"
	"api_server/device/service"
	"api_server/logger"
)

type DeviceController struct {
	svc service.IDeviceService
}

var once sync.Once
var instance *DeviceController

func New(svc service.IDeviceService) *DeviceController {
	once.Do(func() {
		logger.Debug("Device Controller instance")
		instance = &DeviceController{
			svc: svc,
		}
	})

	return instance
}

func (ctlr *DeviceController) CreateOne(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.DeviceDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.Create(reqDTO)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *DeviceController) GetAll(c *gin.Context) {
	logger.ApiRequest(c)

	data, report := ctlr.svc.ReadAll()
	logger.ApiResponse(c, report, data)
}

func (ctlr *DeviceController) UpdateById(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.DeviceDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.Edit(reqDTO)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *DeviceController) DeleteByIds(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.DeviceRemoveDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.DeleteMany(reqDTO.IDs)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *DeviceController) DeleteById(c *gin.Context) {
	logger.ApiRequest(c)

	if device_id, err := strconv.Atoi(c.Param("device_id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.DeleteOne(device_id)
		logger.ApiResponse(c, report, data)
	}
}

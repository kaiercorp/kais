package router

import (
	"sync"

	"github.com/gin-gonic/gin"

	"api_server/device/service"
	"api_server/logger"
)

type GPUController struct {
	svc service.IGPUService
}

var onceGPU sync.Once
var instanceGPU *GPUController

func NewGPU(svc service.IGPUService) *GPUController {
	onceGPU.Do(func() {
		logger.Debug("GPU Controller instance")
		instanceGPU = &GPUController{
			svc: svc,
		}
	})

	return instanceGPU
}

func (ctlr *GPUController) GetAll(c *gin.Context) {
	logger.ApiRequest(c)

	data, report := ctlr.svc.ReadActive()

	logger.ApiResponse(c, report, data)
}

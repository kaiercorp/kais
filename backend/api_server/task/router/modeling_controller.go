package router

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"api_server/logger"
	repo "api_server/task/repository"
	"api_server/task/service"
)

type ModelingController struct {
	svc service.IModelingService
}

var onceModeling sync.Once
var instanceModeling *ModelingController

func NewModelingController(svc service.IModelingService) *ModelingController {
	onceModeling.Do(func() {
		logger.Debug("Modeling Controller instance")
		instanceModeling = &ModelingController{
			svc: svc,
		}
	})

	return instanceModeling
}

func (ctlr *ModelingController) GetByTask(c *gin.Context) {
	logger.ApiRequest(c)

	if task_id, err := strconv.Atoi(c.Param("task_id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.ReadByTask(task_id)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *ModelingController) GetById(c *gin.Context) {
	logger.ApiRequest(c)

	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.ReadFull(id)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *ModelingController) GetModelingType(c *gin.Context) {
	logger.ApiRequest(c)

	if task_id, err := strconv.Atoi(c.Param("task_id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.ReadModelingType(task_id)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *ModelingController) AddEvaluation(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.EvaluationDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.CreateEvaluation(reqDTO)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *ModelingController) StopModeling(c *gin.Context) {
	logger.ApiRequest(c)

	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		report := ctlr.svc.StopModelingTask(id)
		logger.ApiResponse(c, report, nil)
	}
}

func (ctlr *ModelingController) DeleteById(c *gin.Context) {
	logger.ApiRequest(c)

	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		report := ctlr.svc.DeleteOne(id)
		logger.ApiResponse(c, report, nil)
	}
}

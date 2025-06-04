package router

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"api_server/logger"
	repo "api_server/task/repository"
	"api_server/task/service"
)

type TaskController struct {
	svc service.ITaskService
}

var onceTask sync.Once
var instanceTask *TaskController

func NewTaskController(svc service.ITaskService) *TaskController {
	onceTask.Do(func() {
		logger.Debug("Task Controller instance")
		instanceTask = &TaskController{
			svc: svc,
		}
	})

	return instanceTask
}

func (ctlr *TaskController) CreateOne(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.TaskDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		_, _, report := ctlr.svc.Create(reqDTO)
		logger.ApiResponse(c, report, nil)
	}
}

func (ctlr *TaskController) GetByProject(c *gin.Context) {
	logger.ApiRequest(c)

	if project_id, err := strconv.Atoi(c.Param("project_id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.ReadByProject(project_id)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *TaskController) GetOne(c *gin.Context) {
	logger.ApiRequest(c)

	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.ReadOne(id)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *TaskController) UpdateById(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.TaskDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.Edit(reqDTO)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *TaskController) DeleteById(c *gin.Context) {
	logger.ApiRequest(c)

	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		report := ctlr.svc.DeleteOne(id)
		logger.ApiResponse(c, report, nil)
	}
}

func (ctlr *TaskController) DeleteByProject(c *gin.Context) {
	logger.ApiRequest(c)

	if project_id, err := strconv.Atoi(c.Param("project_id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		report := ctlr.svc.DeleteByProject(project_id)
		logger.ApiResponse(c, report, nil)
	}
}

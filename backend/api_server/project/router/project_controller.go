package router

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"api_server/logger"
	repo "api_server/project/repository"
	"api_server/project/service"
	"api_server/utils"
)

type ProjectController struct {
	svc service.IProjectService
}

var once sync.Once
var instance *ProjectController

func New(svc service.IProjectService) *ProjectController {
	once.Do(func() { // atomic, does not allow repeating
		logger.Debug("Project Controller instance")
		instance = &ProjectController{
			svc: svc,
		}
	})

	return instance
}

func (ctlr *ProjectController) CreateOne(c *gin.Context) {
	logger.ApiRequest(c)

	ctxData, err := utils.GetDataFromToken(c)
	if err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}
	reqDTO := repo.ProjectDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}
	data, report := ctlr.svc.Create(reqDTO, ctxData.Username)
	logger.ApiResponse(c, report, data)
}

func (ctlr *ProjectController) GetByPages(c *gin.Context) {
	logger.ApiRequest(c)

	if page, err := strconv.Atoi(c.Query("page")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.Read(page)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *ProjectController) UpdateById(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.ProjectDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.Edit(reqDTO)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *ProjectController) DeleteById(c *gin.Context) {
	logger.ApiRequest(c)

	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.Delete(id)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *ProjectController) GetByPagesWithUser(c *gin.Context) {
	logger.ApiRequest(c)
	ctxData, err := utils.GetDataFromToken(c)
	if err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	if page, err := strconv.Atoi(c.Query("page")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.ReadByUsername(page, ctxData.Username)
		logger.ApiResponse(c, report, data)
	}
}

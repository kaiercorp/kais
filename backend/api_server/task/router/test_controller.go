package router

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"api_server/logger"
	repo "api_server/task/repository"
	"api_server/task/service"
)

type TestController struct {
	svc service.ITestService
}

var onceTest sync.Once
var instanceTest *TestController

func NewTestController(svc service.ITestService) *TestController {
	onceTest.Do(func() {
		logger.Debug("Test Controller instance")
		instanceTest = &TestController{
			svc: svc,
		}
	})

	return instanceTest
}

func (ctlr *TestController) LoadOneModel(c *gin.Context) {
	logger.ApiRequest((c))

	reqDTO := repo.TestDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.LoadModel(reqDTO)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *TestController) UnloadModel(c *gin.Context) {
	logger.ApiRequest(c)

	if test_id, err := strconv.Atoi(c.Param("test_id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		report := ctlr.svc.UnloadModel(test_id)
		logger.ApiResponse(c, report, nil)
	}
}

func (ctlr *TestController) InferenceVCLS(c *gin.Context) {
	logger.ApiRequest(c)

	if test_id, err := strconv.Atoi(c.PostForm("test_id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		// 이미지 파일 가져오기
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			r := logger.CreateReport(&logger.CODE_REQUEST, err)
			logger.ApiResponse(c, r, nil)
		}
		defer file.Close()

		data, report := ctlr.svc.InferenceVCLS(test_id, header.Filename, file)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *TestController) InferenceTabular(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.TabularTestRequest{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.InferenceTabular(reqDTO.TestID, reqDTO.XInputs)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *TestController) GetDatasetColumns(c *gin.Context) {
	logger.ApiRequest(c)

	if modeling_id, err := strconv.Atoi(c.Param("modeling_id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.GetDatasetColumns(modeling_id)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *TestController) PredictXFeatures(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.TabularTestRequest{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.PredictTabular(reqDTO)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *TestController) FeatureImportanceLIME(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.TabularTestRequest{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.FeatureImportanceLIME(reqDTO)
		logger.ApiResponse(c, report, data)
	}
}

func (ctrl *TestController) GetLoaddedModels(c *gin.Context) {
	logger.ApiRequest(c)

	result, err := ctrl.svc.LoadedModels()

	logger.ApiResponse(c, err, result)
}

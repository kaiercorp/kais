package router

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	repo "api_server/dataset/repository"
	"api_server/dataset/service"
	"api_server/logger"
)

type DatasetController struct {
	svc service.DatasetServiceInterface
}

var onceDataset sync.Once
var datasetControllerInstance *DatasetController

func NewDatasetController(datasetService service.DatasetServiceInterface) *DatasetController {
	onceDataset.Do(func() {
		logger.Debug("Dataset Controller instance")
		datasetControllerInstance = &DatasetController{
			svc: datasetService,
		}
	})

	return datasetControllerInstance
}

func (ctlr *DatasetController) GetDatasets(c *gin.Context) {
	logger.ApiRequest(c)

	if page, err := strconv.Atoi(c.Query("page")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		datasetType := c.QueryArray("datasetType[]")
		data, report := ctlr.svc.ViewDatasets(datasetType, page)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *DatasetController) GetTableColumn(c *gin.Context) {
	logger.ApiRequest(c)

	if parent_id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.ViewTableColumn(parent_id)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *DatasetController) GetClasses(c *gin.Context) {
	logger.ApiRequest(c)

	if parent_id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		engine_type := c.Param("engine_type")
		data, report := ctlr.svc.ViewClasses(parent_id, engine_type)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *DatasetController) DeleteDsataset(c *gin.Context) {
	logger.ApiRequest(c)

	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		report := ctlr.svc.RemoveDataset(id)
		logger.ApiResponse(c, report, nil)
	}
}

func (ctlr *DatasetController) FetchDatasetStatistics(c *gin.Context) {
	logger.ApiRequest(c)

	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.GetDataStatistics(id)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *DatasetController) FetchTabularDatasetCompareNumerical(c *gin.Context) {
	logger.ApiRequest(c)

	selectedFeatures := repo.CompareFeaturesStatics{}
	if errParam := c.ShouldBindJSON(&selectedFeatures); errParam != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, errParam)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.SelectTabularDatasetCompareNumerical(selectedFeatures)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *DatasetController) FetchTabularDatasetCompareCategorical(c *gin.Context) {
	logger.ApiRequest(c)

	selectedFeatures := repo.CompareFeaturesStatics{}
	if errParams := c.ShouldBindJSON(&selectedFeatures); errParams != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, errParams)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.SelectTabularDatasetCompareCategorical(selectedFeatures)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *DatasetController) FetchTabularDatasetCompareCategoricalNumerical(c *gin.Context) {
	logger.ApiRequest(c)

	selectedFeatures := repo.CompareFeaturesStatics{}
	if errParams := c.ShouldBindJSON(&selectedFeatures); errParams != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, errParams)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.SelectTabularDatasetCompareCategoricalNumerical(selectedFeatures)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *DatasetController) FetchDatasetStatByType(c *gin.Context) {
	logger.ApiRequest(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}
	statType := c.Param("stat_type")
	filePath, report := ctlr.svc.GetDataStatByTypeFromJSON(id, statType)
	if report != nil {
		logger.ApiResponse(c, report, nil)
		return
	}
	logger.ApiResponseWithJsonFile(c, filePath)
}

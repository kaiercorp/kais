package router

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	repo "api_server/dataset/repository"
	"api_server/dataset/service"
	"api_server/logger"
)

type DatasetRootController struct {
	svc service.IDatasetRootService
}

var onceDatasetRoot sync.Once
var datasetRootControllerInstance *DatasetRootController

func NewDatasetRootController(datasetRootService service.IDatasetRootService) *DatasetRootController {
	onceDatasetRoot.Do(func() {
		logger.Debug("Dataset Root Controller instance")
		datasetRootControllerInstance = &DatasetRootController{
			svc: datasetRootService,
		}
	})

	return datasetRootControllerInstance
}

func (ctlr *DatasetRootController) GetDatasetroots(c *gin.Context) {
	logger.ApiRequest(c)

	data, report := ctlr.svc.ViewDatasetrootAll()
	logger.ApiResponse(c, report, data)
}

func (ctlr *DatasetRootController) DeleteDatasetroot(c *gin.Context) {
	logger.ApiRequest(c)

	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		report := ctlr.svc.RemoveDatasetroot(id)
		logger.ApiResponse(c, report, nil)
	}
}

// @Description Get dataset root paths
// @Tags system
// @Success 200
// @Router /data [get]
func (ctlr *DatasetRootController) GetDatasetrootsForTAPI(c *gin.Context) {
	logger.ApiRequest(c)

	data, report := ctlr.svc.ViewDatasetrootAllForAPI()
	logger.ApiResponse(c, report, data)
}

// @Description Post dataset root path
// @Tags system
// @Success 200
// @Router /data [post]
// @Accept json
// @Param Body body dto.DatasetRootRequest true "dataset root information"
func (ctlr *DatasetRootController) PostDatasetrootForTAPI(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.DatasetRootDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.AddDatasetroot(reqDTO)
		logger.ApiResponse(c, report, data)
	}
}

// @Description Put dataset root path
// @Tags system
// @Success 200
// @Router /data [put]
// @Accept json
// @Param Body body dto.DatasetRootForAPI true "dataset root information"
func (ctlr *DatasetRootController) PutDatasetrootForTAPI(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.DatasetRootDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		reqDTO.IsUse = true
		data, report := ctlr.svc.EditDatasetroot(reqDTO)
		logger.ApiResponse(c, report, data)
	}
}

// @Description Delete dataset root path
// @Tags system
// @Success 200
// @Router /data/{datasetroot_id} [delete]
// @Param datasetroot_id path int true "datasetroot_id"
func (ctlr *DatasetRootController) DeleteDatasetrootForTAPI(c *gin.Context) {
	logger.ApiRequest(c)

	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		report := ctlr.svc.RemoveDatasetroot(id)
		logger.ApiResponse(c, report, nil)
	}
}

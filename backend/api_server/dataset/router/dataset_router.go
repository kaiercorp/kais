package router

import (
	"github.com/gin-gonic/gin"

	"api_server/dataset/modules"
	"api_server/dataset/repository"
	"api_server/dataset/service"
	"api_server/utils"
)

func InitRouter(r *gin.Engine) {
	datasetDAO := repository.NewDatasetDAO()
	datasetRootDAO := repository.NewDatasetRootDAO(datasetDAO)
	datasetWatcher := modules.NewDatasetWatcher(modules.NewDatasetValidator(datasetDAO), modules.NewDatasetAnalyzer(datasetDAO), datasetDAO, datasetRootDAO)
	datasetController := NewDatasetController(service.NewDatasetService(datasetWatcher, modules.NewDatasetAnalyzer(datasetDAO), datasetDAO))
	datasetRootController := NewDatasetRootController(service.NewDatasetRootService(datasetWatcher, datasetRootDAO, datasetDAO))

	apiRouter := r.Group(utils.API_BASE_URL_V1 + "/dataset")
	{
		apiRouter.GET("", datasetController.GetDatasets)
		apiRouter.GET("/column/:id", datasetController.GetTableColumn)
		apiRouter.GET("/classes/:id/:engine_type", datasetController.GetClasses)
		apiRouter.GET("/stat/:id", datasetController.FetchDatasetStatistics)
		apiRouter.GET("/stat/json/:id/:stat_type", datasetController.FetchDatasetStatByType)
		apiRouter.DELETE(":id", datasetController.DeleteDsataset)
		apiRouter.POST("/analyze/tabular/compare/numerical", datasetController.FetchTabularDatasetCompareNumerical)
		apiRouter.POST("/analyze/tabular/compare/categorical", datasetController.FetchTabularDatasetCompareCategorical)
		apiRouter.POST("/analyze/tabular/compare/categoricalNumerical", datasetController.FetchTabularDatasetCompareCategoricalNumerical)
	}

	apiRouterDR := r.Group(utils.API_BASE_URL_V1 + "/dataroot")
	{
		apiRouterDR.GET("", datasetRootController.GetDatasetroots)
		apiRouterDR.DELETE(":id", datasetRootController.DeleteDatasetroot)
	}
}

package controllers

import (
	"github.com/gin-gonic/gin"

	modules_dataset "api_server/dataset/modules"
	repo_dataset "api_server/dataset/repository"
	service_dataset "api_server/dataset/service"
	repo_devce "api_server/device/repository"
	service_device "api_server/device/service"
	service_tapi "api_server/tapi/service"
	repo_task "api_server/task/repository"
	service_task "api_server/task/service"
	"api_server/utils"
)

func InitRouter(r *gin.Engine) {
	device_dao := repo_devce.New()
	gpu_dao := repo_devce.NewGPUDAO()
	task_dao := repo_task.NewTaskDAO()
	modeling_dao := repo_task.NewModelingDAO()
	modeling_details_dao := repo_task.NewModelingDetailDAO()
	dataset_dao := repo_dataset.NewDatasetDAO()
	datasetroot_dao := repo_dataset.NewDatasetRootDAO(dataset_dao)
	device_svc := service_device.New(device_dao)
	gpu_svc := service_device.NewGPUService(gpu_dao)
	task_svc := service_task.NewTaskService(task_dao, modeling_dao, dataset_dao, device_dao)
	modeling_svc := service_task.NewModelingService(modeling_dao, modeling_details_dao, dataset_dao, device_dao)
	datasetWatcher := modules_dataset.NewDatasetWatcher(modules_dataset.NewDatasetValidator(dataset_dao),
		modules_dataset.NewDatasetAnalyzer(dataset_dao), dataset_dao, datasetroot_dao)
	dataset_svc := service_dataset.NewDatasetService(datasetWatcher, modules_dataset.NewDatasetAnalyzer(dataset_dao), dataset_dao)
	datasetroot_svc := service_dataset.NewDatasetRootService(datasetWatcher, datasetroot_dao, dataset_dao)
	tapi_svc := service_tapi.New(device_svc, gpu_svc, task_svc, modeling_svc, dataset_svc, datasetroot_svc, datasetWatcher, modeling_dao)
	controllers := New(tapi_svc)
	apiRouter := r.Group(utils.API_BASE_URL_V1)
	{
		apiRouter.POST("/:engine/modeling", controllers.StartModeling) // Auto Train
		apiRouter.DELETE("/:engine/modeling/:modeling_id", controllers.StopModeling)
		apiRouter.GET("/:engine/modeling/:modeling_id/:threshold", controllers.ListModelingModel) // Modeling 별 Model 목록 조회
		apiRouter.GET("/:engine/list", controllers.ListModeling)                                  // 목록 조회
		//apiRouter.GET("/:engine/logs/:modeling_id", controllers.ListModelingLogs)
	}

	{
		apiRouter.GET("/:engine/load/list", controllers.GetLoaddedModels)
		apiRouter.POST("/:engine/load", controllers.PostLoadModel)
		apiRouter.DELETE("/:engine/load/:gpu_id/:model_name", controllers.DeleteUnloadModel)
		apiRouter.POST("/:engine/inference", controllers.PostStartFileTest)
	}

	apiSysRouter := r.Group(utils.API_BASE_URL_V1 + "/sys")
	{
		apiSysRouter.GET("", controllers.GetSystemStatus)
	}
	apiDatasetRouter := r.Group(utils.API_BASE_URL_V1 + "/dataset")
	{
		apiDatasetRouter.POST("", controllers.UploadDataset)
	}

	//apiTestRouter := r.Group(utils.API_BASE_URL_V1 + "/test")
	{
		//apiTestRouter.GET("/01", controllers.Test01)
	}
}

package router

import (
	"github.com/gin-gonic/gin"

	repo_dataset "api_server/dataset/repository"
	repo_device "api_server/device/repository"
	repo_task "api_server/task/repository"
	service_task "api_server/task/service"
	"api_server/utils"
)

func InitRouter(r *gin.Engine) {
	task_dao := repo_task.NewTaskDAO()
	modeling_dao := repo_task.NewModelingDAO()
	modelingDetails_dao := repo_task.NewModelingDetailDAO()
	dataset_dao := repo_dataset.NewDatasetDAO()
	device_dao := repo_device.New()

	task_service := service_task.NewTaskService(task_dao, modeling_dao, dataset_dao, device_dao)
	taskController := NewTaskController(task_service)

	apiRouter := r.Group(utils.API_BASE_URL_V1 + "/task")
	{
		apiRouter.POST("", taskController.CreateOne)
		apiRouter.GET("/list/:project_id", taskController.GetByProject)
		apiRouter.GET("/:id", taskController.GetOne)
		apiRouter.PUT("", taskController.UpdateById)
		apiRouter.DELETE("/:id", taskController.DeleteById)
	}

	modeling_service := service_task.NewModelingService(modeling_dao, modelingDetails_dao, dataset_dao, device_dao)
	modelingController := NewModelingController(modeling_service)

	apiModelingRouter := r.Group(utils.API_BASE_URL_V1 + "/modeling")
	{
		apiModelingRouter.GET("/list/:task_id", modelingController.GetByTask)
		apiModelingRouter.GET("/testable/:task_id", instanceModeling.GetModelingType)
		apiModelingRouter.GET("/:id", modelingController.GetById)
		apiModelingRouter.POST("/evaluation", modelingController.AddEvaluation)
		apiModelingRouter.DELETE("/stop/:id", modelingController.StopModeling)
		apiModelingRouter.DELETE("/:id", modelingController.DeleteById)
	}

	testService := service_task.NewTestService(task_dao, modeling_dao)
	testController := NewTestController(testService)

	apiTestRouter := r.Group(utils.API_BASE_URL_V1 + "/test")
	{
		apiTestRouter.POST("/load_one", testController.LoadOneModel)
		apiTestRouter.DELETE("/unload/:test_id", testController.UnloadModel)
		apiTestRouter.POST("/vcls", testController.InferenceVCLS)
		apiTestRouter.POST("/tabular", testController.InferenceTabular)
		apiTestRouter.GET("/tabular/columns/:modeling_id", testController.GetDatasetColumns)
		apiTestRouter.POST("/tabular/predict", testController.PredictXFeatures)
		apiTestRouter.POST("/tabular/lime", testController.FeatureImportanceLIME)
	}
}

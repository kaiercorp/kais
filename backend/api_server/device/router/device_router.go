package router

import (
	"github.com/gin-gonic/gin"

	"api_server/device/repository"
	"api_server/device/service"
	"api_server/utils"
)

func InitRouter(r *gin.Engine) {
	repo := repository.New()
	svc := service.New(repo)
	controller := New(svc)

	apiRouter := r.Group(utils.API_BASE_URL_V1 + "/device")
	{
		apiRouter.POST("", controller.CreateOne)
		apiRouter.GET("", controller.GetAll)
		apiRouter.PUT("", controller.UpdateById)
		apiRouter.DELETE("", controller.DeleteByIds)
		apiRouter.DELETE("/:device_id", controller.DeleteById)
	}

	repo_gpu := repository.NewGPUDAO()
	svc_gpu := service.NewGPUService(repo_gpu)
	controller_gpu := NewGPU(svc_gpu)

	apiRouterGPU := r.Group(utils.API_BASE_URL_V1 + "/gpu")
	{
		apiRouterGPU.GET("", controller_gpu.GetAll)
	}
}

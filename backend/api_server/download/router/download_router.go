package router

import (
	"github.com/gin-gonic/gin"

	"api_server/download/service"
	"api_server/utils"
)

func InitRouter(r *gin.Engine) {
	svc := service.New()
	controller := New(svc)

	apiRouter := r.Group(utils.API_BASE_URL_V1 + "/download")
	{
		apiRouter.GET("/:modeling_id/:model_name", controller.DownloadModel)
	}
}

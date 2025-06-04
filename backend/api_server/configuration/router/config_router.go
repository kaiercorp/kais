package router

import (
	"github.com/gin-gonic/gin"

	auth_repo "api_server/auth/repository"
	auth_service "api_server/auth/service"
	"api_server/configuration/repository"
	"api_server/configuration/service"
	"api_server/utils"
)

func InitRouter(r *gin.Engine) {
	repo := repository.New()
	svc := service.New(repo)
	auth_repo := auth_repo.New()
	auth_svc := auth_service.NewUser(auth_repo)
	controller := New(svc, auth_svc)

	apiRouter := r.Group(utils.API_BASE_URL_V1 + "/configuration")
	{
		apiRouter.GET("", controller.GetAll)
		apiRouter.POST("", controller.UpdateMany)
		apiRouter.POST("/lang", controller.UpdateLanguage)
	}
}

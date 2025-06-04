package router

import (
	"github.com/gin-gonic/gin"

	"api_server/auth/repository"
	"api_server/auth/service"
	"api_server/utils"
)

func InitAuthRouter(r *gin.Engine) {
	repo := repository.New()
	svc := service.NewAuth(repo)
	controller := New(svc)

	apiRouter := r.Group(utils.API_BASE_URL_V1 + "/auth")
	{
		apiRouter.POST("/login", controller.Login)
		apiRouter.POST("/logout", controller.Logout)
	}
}

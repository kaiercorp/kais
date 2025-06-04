package router

import (
	"github.com/gin-gonic/gin"

	"api_server/user/repository"
	"api_server/user/service"
	"api_server/utils"
)

func InitRouter(r *gin.Engine) {

	usr_dao := repository.NewUserDAO()
	usr_svc := service.NewUserService(usr_dao)
	controller := New(usr_svc)
	apiRouter := r.Group(utils.API_BASE_URL_V1 + "/user")
	{
		apiRouter.Use(utils.JWTAuthMiddleware())
		apiRouter.GET("/:username", utils.GroupMiddleware(0, 1, 2), controller.GetUser)
		apiRouter.GET("/group/:group", utils.GroupMiddleware(0, 1), controller.GetUsersByGroupGT)
		apiRouter.POST("", utils.GroupMiddleware(0, 1), controller.CreateUser)
		apiRouter.PUT("", utils.GroupMiddleware(0, 1, 2), controller.UpdateUser)
		apiRouter.PUT("/activate/:id", utils.GroupMiddleware(0, 1), controller.ChangeActivateStatus)
		apiRouter.PUT("/reset-password/:id", utils.GroupMiddleware(0, 1), controller.ResetPassword)
		apiRouter.PUT("/password", utils.GroupMiddleware(0, 1, 2), controller.ChangePassword)
		apiRouter.DELETE("/:id", utils.GroupMiddleware(0, 1), controller.DeleteUser)
	}

}

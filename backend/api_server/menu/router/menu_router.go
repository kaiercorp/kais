package router

import (
	"github.com/gin-gonic/gin"

	auth_svc "api_server/auth/service"
	"api_server/menu/controllers"
	"api_server/menu/database"
	"api_server/menu/service"
	"api_server/utils"
)

func InitRouter(r *gin.Engine) {
	menuController := controllers.NewMenuController(service.NewMenuService(database.NewMenuDAO()), auth_svc.NewUserService())

	apiRouter := r.Group(utils.API_BASE_URL_V1 + "/menu")
	{
		apiRouter.GET("", menuController.GetMenus)
		apiRouter.POST("", menuController.SetMenus)
	}
}

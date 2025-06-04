package router

import (
	"github.com/gin-gonic/gin"

	"api_server/utils"
)

func InitWebSocketRouter(r *gin.Engine) {
	wsRouter := r.Group(utils.WS_BASE_URL_v1)
	{
		wsRouter.GET("", HandleWebSocketMessage)
	}
}

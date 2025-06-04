package router

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	repo "api_server/auth/repository"
	"api_server/auth/service"
	"api_server/logger"
)

type LoginController struct {
	svc service.ILoginService
}

var once sync.Once
var instance *LoginController

func New(svc service.ILoginService) *LoginController {
	once.Do(func() {
		logger.Debug("Login Controller instance")
		instance = &LoginController{
			svc: svc,
		}
	})

	return instance
}

func (ctlr *LoginController) Login(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.LoginRequest{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		data, report := ctlr.svc.Login(reqDTO)
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", data.Token, 3600*6, "", "", false, true)
		logger.ApiResponse(c, report, data)
	}
}

func (ctlr *LoginController) Logout(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.LogoutRequest{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
	} else {
		report := ctlr.svc.Logout(reqDTO)
		logger.ApiResponse(c, report, nil)
	}
}

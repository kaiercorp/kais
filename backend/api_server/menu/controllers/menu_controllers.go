package controllers

import (
	"github.com/gin-gonic/gin"

	auth_svc "api_server/auth/service"
	"api_server/logger"
	"api_server/menu/dto"
	"api_server/menu/service"
)

type MenuController struct {
	menuService service.MenuServiceInterface
	userService auth_svc.IUserService
}

var menuControllerInstance *MenuController

func NewMenuController(menuService service.MenuServiceInterface, userService auth_svc.IUserService) *MenuController {
	if menuControllerInstance == nil {
		menuControllerInstance = &MenuController{
			menuService: menuService,
			userService: userService,
		}
	}

	return menuControllerInstance
}

// @Description Get menus
// @Tags menu
// @Success 200
// @Router /menu [get]
func (this *MenuController) GetMenus(c *gin.Context) {
	logger.ApiRequest(c)

	token := c.GetHeader("Authorization")
	group := this.userService.GetUserGroup(token)

	menus, r := this.menuService.ViewMenus(group)
	if r != nil {
		logger.ApiResponse(c, r, nil)
	}

	logger.ApiResponse(c, nil, menus)
}

// @Description Update usable status of menus
// @Tags menu
// @Success 200
// @Router /menu [post]
// @Param dto body dto.MenuDTO true "MenuDTO"
func (this *MenuController) SetMenus(c *gin.Context) {
	logger.ApiRequest(c)

	menus := []dto.MenuDTO{}
	if err := c.ShouldBindJSON(&menus); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	r := this.menuService.EditMenus(menus)
	if r != nil {
		logger.ApiResponse(c, r, nil)
		return
	}
}

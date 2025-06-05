package service

import (
	"context"

	"api_server/ent"
	"api_server/logger"
	"api_server/menu/database"
	"api_server/menu/dto"
)

type MenuServiceInterface interface {
	ViewMenus(int) ([]dto.MenuDTO, *logger.Report)
	EditMenus([]dto.MenuDTO) *logger.Report
}

type MenuService struct {
	ctx     context.Context
	menuDAO database.MenuDAOInterface
}

var menuServiceInstance *MenuService

func NewMenuService(menuDAO database.MenuDAOInterface) *MenuService {
	if menuServiceInstance == nil {
		menuServiceInstance = &MenuService{
			ctx:     context.Background(),
			menuDAO: menuDAO,
		}
	}

	return menuServiceInstance
}

func (this *MenuService) ViewMenus(group int) ([]dto.MenuDTO, *logger.Report) {
	selectedMenus, r := this.menuDAO.SelectMenusByGroup(this.ctx, group)

	if r != nil {
		return nil, r
	}

	menuDTO := this.convertMenuEntityToMenuDTO(selectedMenus)

	return menuDTO, nil
}

func (this *MenuService) EditMenus(menus []dto.MenuDTO) *logger.Report {
	return this.menuDAO.UpdateMenus(this.ctx, menus)
}

func (this *MenuService) convertMenuEntityToMenuDTO(menuEntity []*ent.Menu) []dto.MenuDTO {
	if menuEntity == nil {
		return nil
	}

	menuDTO := []dto.MenuDTO{}

	for _, menu := range menuEntity {
		menuDTO = append(menuDTO, dto.MenuDTO{
			Key:       menu.ID,
			Label:     menu.Label,
			Icon:      menu.Icon,
			Url:       menu.URL,
			IsUse:     menu.IsUse,
			IsTitle:   menu.IsTitle,
			ParentKey: menu.ParentKey,
			MenuOrder: menu.MenuOrder,
			Group:     menu.Group,
			Children:  this.convertMenuEntityToMenuDTO(menu.Edges.Children),
		})
	}

	return menuDTO
}

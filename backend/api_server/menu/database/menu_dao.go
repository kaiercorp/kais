package database

import (
	"context"

	"api_server/ent"
	"api_server/ent/menu"
	"api_server/logger"
	"api_server/menu/dto"
	"api_server/utils"

	"entgo.io/ent/dialect/sql"
)

type MenuDAOInterface interface {
	SelectMenus(context.Context) ([]*ent.Menu, *logger.Report)
	SelectMenusByGroup(context.Context, int) ([]*ent.Menu, *logger.Report)
	UpdateMenu(context.Context, dto.MenuDTO) *logger.Report
	UpdateMenus(context.Context, []dto.MenuDTO) *logger.Report
}

type MenuDAO struct {
	entClient *ent.Client
}

var menuDAOInstance *MenuDAO

func NewMenuDAO() *MenuDAO {
	if menuDAOInstance == nil {
		menuDAOInstance = &MenuDAO{
			entClient: utils.GetEntClient(),
		}
	}

	return menuDAOInstance
}

func (dao *MenuDAO) SelectMenus(ctx context.Context) ([]*ent.Menu, *logger.Report) {
	menus, err := dao.entClient.Menu.
		Query().
		Where(menu.Not(menu.HasParent())).
		WithChildren().
		Order(menu.ByMenuOrder(sql.OrderAsc())).
		All(ctx)

	if err != nil {
		return menus, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	return menus, nil
}

func (dao *MenuDAO) SelectMenusByGroup(ctx context.Context, group int) ([]*ent.Menu, *logger.Report) {
	if group == 0 {
		return dao.selectMenusForMaster(ctx)
	}

	return dao.selectMenusByGroupId(ctx, group)
}

func (dao *MenuDAO) selectMenusByGroupId(ctx context.Context, group int) ([]*ent.Menu, *logger.Report) {
	menus, err := dao.entClient.Menu.
		Query().
		Where(menu.And(menu.GroupGTE(group), menu.IsUse(true), menu.Not(menu.HasParent()))).
		WithChildren(func(m *ent.MenuQuery) {
			m.Where(menu.And(menu.GroupGTE(group), menu.IsUse(true)))
		}).
		Order(menu.ByMenuOrder(sql.OrderAsc())).
		All(ctx)

	if err != nil {
		return menus, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	return menus, nil
}

func (dao *MenuDAO) selectMenusForMaster(ctx context.Context) ([]*ent.Menu, *logger.Report) {
	menus, err := dao.entClient.Menu.
		Query().
		Where(menu.Not(menu.HasParent())).
		WithChildren().
		Order(menu.ByMenuOrder(sql.OrderAsc())).
		All(ctx)

	if err != nil {
		return menus, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	return menus, nil
}

func (dao *MenuDAO) UpdateMenu(ctx context.Context, menuDTO dto.MenuDTO) *logger.Report {
	err := dao.entClient.Menu.
		Update().
		Where(
			menu.ID(menuDTO.Key),
		).
		SetIsUse(menuDTO.IsUse).
		Exec(ctx)

	if err != nil {
		return logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	return nil
}

func (dao *MenuDAO) UpdateMenus(ctx context.Context, menus []dto.MenuDTO) *logger.Report {
	for _, menu := range menus {
		dao.UpdateMenu(ctx, menu)

		if menu.Children != nil {
			for _, child := range menu.Children {
				dao.UpdateMenu(ctx, child)
			}
		}
	}

	return nil
}

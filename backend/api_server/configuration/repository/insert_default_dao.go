package repository

import (
	"api_server/ent"
	"api_server/ent/configuration"
	"api_server/ent/menu"
	"api_server/ent/user"
	"api_server/ent/usergroup"
	"api_server/utils"
	"context"
	"sync"

	"entgo.io/ent/dialect/sql"
)

type IConfigInitDAO interface {
	InsertDefaultValues(ctx context.Context, version string)
}

type ConfigInitDAO struct {
	dbms *ent.Client
}

var onceInit sync.Once
var instanceInit *ConfigInitDAO

func NewConfigInit() *ConfigInitDAO {
	onceInit.Do(func() {
		println("Config Init DAO instance")
		instanceInit = &ConfigInitDAO{
			dbms: utils.GetEntClient(),
		}
	})

	return instanceInit
}

func (dao *ConfigInitDAO) InsertDefaultValues(ctx context.Context, version string) {
	dao.insertSystemConfigs(ctx)
	dao.insertUserConfigs(ctx, version)
	dao.insertUserGroup(ctx)
	dao.insertMasterUser(ctx)
	dao.insertMenu(ctx)
}

func (dao *ConfigInitDAO) insertUserGroup(ctx context.Context) {
	if err := dao.dbms.UserGroup.Create().
		SetIsUse(true).
		SetLevel(0).
		SetName("Master").
		OnConflict(
			sql.ConflictColumns(usergroup.FieldLevel),
		).
		Update(func(ugc *ent.UserGroupUpsert) {
			ugc.SetLevel(0)
		}).
		Exec(ctx); err != nil {
		println("INIT Auth config: ", err.Error())
		return
	}

	if err := dao.dbms.UserGroup.Create().
		SetIsUse(true).
		SetLevel(1).
		SetName("Kaier").
		OnConflict(
			sql.ConflictColumns(usergroup.FieldLevel),
		).
		Update(func(ugc *ent.UserGroupUpsert) {
			ugc.SetLevel(1)
		}).
		Exec(ctx); err != nil {
		println("INIT Auth config: ", err.Error())
	}
}

func (dao *ConfigInitDAO) insertMasterUser(ctx context.Context) {
	if err := dao.dbms.User.Create().
		SetName("MASTER").
		SetGroup(0).
		SetUsername("kaieradmin").
		SetPassword("uAtTnoT+KbJ/iD5CbV03Vw==").
		SetToken(("1234")).
		OnConflict(
			sql.ConflictColumns(user.FieldUsername),
		).
		Update(func(u *ent.UserUpsert) {
			u.SetUsername("kaieradmin")
		}).
		Exec(ctx); err != nil {
		println("INIT Auth config: ", err.Error())
	}
}

func (dao *ConfigInitDAO) insertSystemConfigs(ctx context.Context) {
	err := dao.dbms.Configuration.CreateBulk(
		dao.dbms.Configuration.Create().SetConfigType(utils.CONFIG_TYPE_SYSTEM).
			SetConfigKey("LOG_LEVEL").SetConfigVal("INFO"),
		dao.dbms.Configuration.Create().SetConfigType(utils.CONFIG_TYPE_SYSTEM).
			SetConfigKey("ROOT_PATH").SetConfigVal("/kaier/workspace"),
		dao.dbms.Configuration.Create().SetConfigType(utils.CONFIG_TYPE_SYSTEM).
			SetConfigKey("PATH_STATIC_TEST").SetConfigVal("/kaier/workspace/static"),
		dao.dbms.Configuration.Create().SetConfigType(utils.CONFIG_TYPE_SYSTEM).
			SetConfigKey("KAIS_PATH").SetConfigVal("/kaier"),
	).
		OnConflict(
			sql.ConflictColumns(configuration.FieldConfigType, configuration.FieldConfigKey),
		).
		DoNothing().
		Exec(ctx)

	if err != nil {
		println("INIT System config1: ", err.Error())
		return
	}
}

func (dao *ConfigInitDAO) insertUserConfigs(ctx context.Context, version string) {
	err := dao.dbms.Configuration.CreateBulk(
		dao.dbms.Configuration.Create().SetConfigType(utils.CONFIG_TYPE_USER).
			SetConfigKey("LANGUAGE").SetConfigVal("ko"),
		dao.dbms.Configuration.Create().SetConfigType(utils.CONFIG_TYPE_USER).
			SetConfigKey("VERSION").SetConfigVal(version),
		dao.dbms.Configuration.Create().SetConfigType(utils.CONFIG_TYPE_USER).
			SetConfigKey("PATH_LOG_DIR").SetConfigVal("/kaier/log"),
	).
		OnConflict(
			sql.ConflictColumns(configuration.FieldConfigType, configuration.FieldConfigKey),
		).
		DoNothing().
		Exec(ctx)

	if err != nil {
		println("INIT auth config: ", err.Error())
		return
	}
}

func (dao *ConfigInitDAO) insertMenu(ctx context.Context) {
	err := dao.dbms.Menu.CreateBulk(
		dao.dbms.Menu.Create().SetID("home").
			SetLabel("Home").SetIcon("mdi mdi-home").
			SetURL("/dashboard").SetIsUse(true).SetIsTitle(false).SetMenuOrder(1).SetGroup(2),
		dao.dbms.Menu.Create().SetID("vision").
			SetLabel("Vision data").SetIcon("mdi mdi-television-guide").
			SetURL("").SetIsUse(true).SetIsTitle(false).SetMenuOrder(2).SetGroup(2),
		dao.dbms.Menu.Create().SetID(utils.JOB_TYPE_VISION_CLS_SL).
			SetLabel("Single Label Classification").SetIcon("mdi mdi-book-settings").
			SetURL("/vision/vcls-sl").SetIsUse(true).SetIsTitle(false).SetMenuOrder(1).
			SetParentKey("vision").SetGroup(2),
		dao.dbms.Menu.Create().SetID(utils.JOB_TYPE_VISION_CLS_ML).
			SetLabel("Multi Label Classification").SetIcon("mdi mdi-book-settings").
			SetURL("/vision/vcls-ml").SetIsUse(true).SetIsTitle(false).SetMenuOrder(2).
			SetParentKey("vision").SetGroup(2),
		dao.dbms.Menu.Create().SetID(utils.JOB_TYPE_VISION_AD).
			SetLabel("Anomaly Detection").SetIcon("mdi mdi-book-settings").
			SetURL("/vision/vad").SetIsUse(false).SetIsTitle(false).SetMenuOrder(3).
			SetParentKey("vision").SetGroup(2),
		dao.dbms.Menu.Create().SetID("table").
			SetLabel("Table data").SetIcon("mdi mdi-table").
			SetURL("").SetIsUse(false).SetIsTitle(false).SetMenuOrder(3).SetGroup(2),
		dao.dbms.Menu.Create().SetID(utils.JOB_TYPE_TABLE_CLS).
			SetLabel("Classification").SetIcon("mdi mdi-book-settings").
			SetURL("/table/tcls").SetIsUse(true).SetIsTitle(false).SetMenuOrder(1).
			SetParentKey("table").SetGroup(2),
		dao.dbms.Menu.Create().SetID(utils.JOB_TYPE_TABLE_REG).
			SetLabel("Regression").SetIcon("mdi mdi-book-settings").
			SetURL("/table/treg").SetIsUse(true).SetIsTitle(false).SetMenuOrder(2).
			SetParentKey("table").SetGroup(2),
		dao.dbms.Menu.Create().SetID("ts").
			SetLabel("Time-series data").SetIcon("mdi mdi-clock-time-eight-outline").
			SetURL("").SetIsUse(false).SetIsTitle(false).SetMenuOrder(4).SetGroup(2),
		dao.dbms.Menu.Create().SetID(utils.JOB_TYPE_TS_AD).
			SetLabel("Anomaly Detection").SetIcon("mdi mdi-book-settings").
			SetURL("/ts/ad").SetIsUse(true).SetIsTitle(false).SetMenuOrder(1).
			SetParentKey("ts").SetGroup(2),
		dao.dbms.Menu.Create().SetID("dataset").
			SetLabel("Dataset Management").SetIcon("mdi mdi-folder-check-outline").
			SetURL("/dataset").SetIsUse(true).SetIsTitle(false).SetMenuOrder(5).SetGroup(2),
		dao.dbms.Menu.Create().SetID("configuration").
			SetLabel("Configuration").SetIcon("mdi mdi-cog").
			SetURL("").SetIsUse(true).SetIsTitle(false).SetMenuOrder(9).SetGroup(2),
		dao.dbms.Menu.Create().SetID("user").
			SetLabel("User").SetIcon("mdi mdi-cog").
			SetURL("/config/user").SetIsUse(true).SetIsTitle(false).SetMenuOrder(1).
			SetParentKey("configuration").SetGroup(0),
		dao.dbms.Menu.Create().SetID("setting").
			SetLabel("Setting").SetIcon("mdi mdi-cog").
			SetURL("/config/setting").SetIsUse(true).SetIsTitle(false).SetMenuOrder(2).
			SetParentKey("configuration").SetGroup(2),
		dao.dbms.Menu.Create().SetID("system").
			SetLabel("System").SetIcon("mdi mdi-cog").
			SetURL("/config/system").SetIsUse(true).SetIsTitle(false).SetMenuOrder(3).
			SetParentKey("configuration").SetGroup(0),
		dao.dbms.Menu.Create().SetID("menus").
			SetLabel("Menu").SetIcon("mdi mdi-cog").
			SetURL("/config/menu").SetIsUse(true).SetIsTitle(false).SetMenuOrder(4).
			SetParentKey("configuration").SetGroup(1),
		dao.dbms.Menu.Create().SetID("device").
			SetLabel("Device Management").SetIcon("mdi mdi-vector-circle").
			SetURL("/device").SetIsUse(true).SetIsTitle(false).SetMenuOrder(8).SetGroup(2),
	).
		OnConflict(
			sql.ConflictColumns(menu.FieldID),
		).
		UpdateURL().
		Exec(ctx)

	if err != nil {
		println("INIT Route config: ", err.Error())
		return
	}
}

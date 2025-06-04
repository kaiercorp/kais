package repository

import (
	"api_server/ent"
	"api_server/ent/configuration"
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

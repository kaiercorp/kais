package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"api_server/ent"
	"api_server/ent/configuration"
	"api_server/utils"

	"entgo.io/ent/dialect/sql"
)

type IConfigDAO interface {
	InsertOne(ctx context.Context, req ConfigDTO) (*ent.Configuration, error)
	SelectAllByTypeAndKey(ctx context.Context, cfgType string, cfgKey string) ([]*ent.Configuration, error)
	SelectOneByKey(ctx context.Context, cfgKey string) (*ent.Configuration, error)
	SelectAllByType(ctx context.Context, cfgType string) ([]*ent.Configuration, error)
	UpdateOne(ctx context.Context, req ConfigDTO) (*ent.Configuration, error)
	UpdateMany(ctx context.Context, configs []ConfigDTO) error
}

type ConfigDAO struct {
	dbms *ent.Client
}

var once sync.Once
var instance *ConfigDAO

func New() *ConfigDAO {
	once.Do(func() {
		println("Config DAO instance")
		instance = &ConfigDAO{
			dbms: utils.GetEntClient(),
		}
	})

	return instance
}

func (dao *ConfigDAO) InsertOne(ctx context.Context, req ConfigDTO) (*ent.Configuration, error) {
	fmt.Printf("%+v\n", req)
	return dao.dbms.Configuration.Create().
		SetConfigType(req.ConfigType).
		SetConfigKey(req.ConfigKey).
		SetConfigVal(req.ConfigVal).
		Save(ctx)
}

func (dao *ConfigDAO) SelectAllByTypeAndKey(ctx context.Context, cfgType string, cfgKey string) ([]*ent.Configuration, error) {
	// fmt.Printf("{'cfgType': %s, 'cfgKey': %s}\n", cfgType, cfgKey)
	return dao.dbms.Configuration.Query().
		Where(
			configuration.ConfigType(cfgType),
			configuration.ConfigKey(cfgKey),
		).
		Order(configuration.ByID(sql.OrderDesc())).
		All(ctx)
}

func (dao *ConfigDAO) SelectOneByKey(ctx context.Context, cfgKey string) (*ent.Configuration, error) {
	// fmt.Printf("{'cfgKey': %s}\n", cfgKey)
	return dao.dbms.Configuration.Query().
		Where(
			configuration.ConfigKey(cfgKey),
		).
		Order(configuration.ByID(sql.OrderDesc())).
		Limit(1).
		Only(ctx)
}

func (dao *ConfigDAO) SelectAllByType(ctx context.Context, cfgType string) ([]*ent.Configuration, error) {
	// fmt.Printf("{'cfgType': %s}\n", cfgType)
	return dao.dbms.Configuration.Query().
		Where(
			configuration.ConfigType(cfgType),
		).
		Order(configuration.ByID()).
		All(ctx)
}

func (dao *ConfigDAO) UpdateOne(ctx context.Context, req ConfigDTO) (*ent.Configuration, error) {
	// fmt.Printf("%+v\n", req)
	return dao.dbms.Configuration.
		UpdateOneID(req.ID).
		SetConfigVal(req.ConfigVal).
		SetUpdatedAt(time.Now()).
		Save(ctx)
}

func (dao *ConfigDAO) UpdateMany(ctx context.Context, configs []ConfigDTO) error {
	for _, item := range configs {
		if item.ID == 0 {
			if _, err := dao.InsertOne(ctx, item); err != nil {
				return err
			}
		} else {
			if _, err := dao.UpdateOne(ctx, item); err != nil {
				return err
			}
		}
	}

	return nil
}

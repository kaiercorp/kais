package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"api_server/ent"
	"api_server/ent/user"
	"api_server/logger"
	"api_server/utils"
)

type IAuthDAO interface {
	SelectByUsername(ctx context.Context, req LoginRequest) (*ent.User, error)
	SelectByToken(ctx context.Context, token string) (*ent.User, error)
	UpdateLogin(ctx context.Context, req *ent.User) (*ent.User, error)
	UpdateLogout(ctx context.Context, req LogoutRequest) error
	SelectUserByToken(ctx context.Context, token string) (*ent.User, *logger.Report)
}

type AuthDAO struct {
	dbms *ent.Client
}

var once sync.Once
var instance *AuthDAO

func New() *AuthDAO {
	once.Do(func() {
		logger.Debug("Auth DAO intance")
		instance = &AuthDAO{
			dbms: utils.GetEntClient(),
		}
	})

	return instance
}

func (dao *AuthDAO) SelectByUsername(ctx context.Context, req LoginRequest) (*ent.User, error) {
	logger.Debug(fmt.Sprintf("%+v", req))
	return dao.dbms.User.Query().
		Where(
			user.Username(req.Username),
			user.Password(req.Password),
		).
		Only(ctx)
}

func (dao *AuthDAO) SelectByToken(ctx context.Context, token string) (*ent.User, error) {
	logger.Debug(fmt.Sprintf(`{"token": %s}`, token))
	return dao.dbms.User.Query().
		Where(
			user.Token(token),
		).
		Only(ctx)
}

func (dao *AuthDAO) UpdateLogin(ctx context.Context, req *ent.User) (*ent.User, error) {
	logger.Debug(fmt.Sprintf("%+v", req))
	req.UpdatedAt = time.Now()
	if _, err := dao.dbms.User.Update().
		SetToken(req.Token).
		SetLoginAt(req.UpdatedAt).
		Where(
			user.Username(req.Username),
		).
		Save(ctx); err != nil {
		return req, err
	} else {
		return req, nil
	}
}

func (dao *AuthDAO) UpdateLogout(ctx context.Context, req LogoutRequest) error {
	logger.Debug(fmt.Sprintf("%+v", req))
	return dao.dbms.User.Update().
		Where(
			user.Username(req.Username),
			user.Token(req.Token),
		).
		SetToken("").
		Exec(ctx)
}

func (dao *AuthDAO) SelectUserByToken(ctx context.Context, token string) (*ent.User, *logger.Report) {
	user, err := dao.dbms.User.
		Query().
		Where(
			user.Token(token),
		).
		Only(ctx)

	if err != nil {
		return user, logger.CreateReport(&logger.CODE_TOKEN_EXPIRED, err)
	}

	return user, nil
}

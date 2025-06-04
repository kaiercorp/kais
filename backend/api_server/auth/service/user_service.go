package service

import (
	"context"
	"sync"

	repo "api_server/auth/repository"
	"api_server/logger"
)

type IUserService interface {
	ReadUserGroup(token string) int
}

type UserService struct {
	ctx context.Context
	dao repo.IAuthDAO
}

var onceUser sync.Once
var instanceUser *UserService

func NewUser(dao repo.IAuthDAO) *UserService {
	onceUser.Do(func() {
		logger.Debug("User Service instance")
		instanceUser = &UserService{
			ctx: context.Background(),
			dao: dao,
		}
	})

	return instanceUser
}

func (svc *UserService) ReadUserGroup(token string) int {
	user, r := svc.dao.SelectByToken(svc.ctx, token)
	if r != nil || user.Token == "" {
		return 2
	}

	return user.Group
}

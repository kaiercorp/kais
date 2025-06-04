package service

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	//todo: update to jwt-go/v5
	"github.com/dgrijalva/jwt-go/v4"

	repo "api_server/auth/repository"
	"api_server/logger"
	"api_server/utils"
)

type ILoginService interface {
	// Login은 사용자 이름과 비밀번호로 로그인 처리 후, 로그인 응답을 반환합니다.
	//
	// 매개변수:
	//   - req: 로그인 요청 DTO (LoginRequest)
	//
	// 반환값:
	//   - *repo.LoginResponse: 로그인 응답 객체 (사용자 정보, 토큰 등)
	//   - *logger.Report: 에러 발생 시 에러 보고 객체
	Login(req repo.LoginRequest) (*repo.LoginResponse, *logger.Report)

	// Logout은 주어진 요청에 따라 로그아웃 처리합니다.
	//
	// 매개변수:
	//   - req: 로그아웃 요청 DTO (LogoutRequest)
	//
	// 반환값:
	//   - *logger.Report: 에러 발생 시 에러 보고 객체
	Logout(req repo.LogoutRequest) *logger.Report

	// IsLogin은 주어진 토큰으로 로그인 상태를 확인합니다.
	//
	// 매개변수:
	//   - token: 확인할 JWT 토큰
	//
	// 반환값:
	//   - bool: 로그인 상태 (true: 로그인됨, false: 로그인 안됨)
	//   - *logger.Report: 에러 발생 시 에러 보고 객체
	IsLogin(token string) (bool, *logger.Report)
}

// LoginService는 ILoginService 인터페이스를 구현하는 구조체로,
// 실제 로그인, 로그아웃, 로그인 상태 체크 작업을 수행합니다.
type LoginService struct {
	ctx context.Context
	dao repo.IAuthDAO
}

var once sync.Once
var instance *LoginService

func NewAuth(dao repo.IAuthDAO) *LoginService {
	once.Do(func() {
		logger.Debug("Login Service instance")
		instance = &LoginService{
			ctx: context.Background(),
			dao: dao,
		}
	})

	return instance
}

// Login은 주어진 로그인 요청을 처리하여 로그인 응답을 반환합니다.
// 매개변수:
//   - req: 로그인 요청 DTO (LoginRequest)
//
// 반환값:
//   - *repo.LoginResponse: 로그인 응답 객체 (사용자 정보, 토큰 등)
//   - *logger.Report: 에러 발생 시 에러 보고 객체
func (svc *LoginService) Login(req repo.LoginRequest) (*repo.LoginResponse, *logger.Report) {
	logger.Debug(fmt.Sprintf("%+v", req))
	s := utils.CreateSecurity()

	// 비밀번호 암호화
	encPwd, _ := s.Encrypt(req.Password)
	req.Password = encPwd

	// 사용자 조회
	selectedUser, err := svc.dao.SelectByUsername(svc.ctx, req)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_LOGIN_PARAMS, err)
	}
	if !selectedUser.IsUse {
		return nil, logger.CreateReport(&logger.CODE_LOGIN_FAILED, fmt.Errorf("user is not available"))
	}

	// 토큰 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": selectedUser.Username,
		"group":    selectedUser.Group,
		"exp":      time.Now().Add(time.Hour * 10).Unix(),
	})

	// 토큰 서명
	if tokenString, err := token.SignedString([]byte(os.Getenv("SECRETKKAIER"))); err == nil {
		selectedUser.Token = tokenString
	} else {
		return nil, logger.CreateReport(&logger.CODE_LOGIN_PARAMS, err)
	}

	// 로그인 정보 업데이트
	updatedUser, r2 := svc.dao.UpdateLogin(svc.ctx, selectedUser)
	if r2 != nil {
		return nil, logger.CreateReport(&logger.CODE_LOGIN_FAILED, r2)
	}

	// 로그인 응답 반환
	return &repo.LoginResponse{
		Username: updatedUser.Username,
		Token:    updatedUser.Token,
		LoginAt:  updatedUser.LoginAt,
		Name:     updatedUser.Name,
		Group:    &updatedUser.Group,
	}, nil
}

// Logout은 주어진 로그아웃 요청을 처리하고 로그아웃 상태로 업데이트합니다.
// 매개변수:
//   - req: 로그아웃 요청 DTO (LogoutRequest)
//
// 반환값:
//   - *logger.Report: 에러 발생 시 에러 보고 객체
func (svc *LoginService) Logout(req repo.LogoutRequest) *logger.Report {
	logger.Debug(fmt.Sprintf("%+v", req))
	if err := svc.dao.UpdateLogout(svc.ctx, req); err != nil {
		return logger.CreateReport(&logger.CODE_LOGOUT_FAILED, err)
	}

	return nil
}

// IsLogin은 주어진 토큰으로 로그인 상태를 확인합니다.
// 매개변수:
//   - token: 확인할 JWT 토큰
//
// 반환값:
//   - bool: 로그인 상태 (true: 로그인됨, false: 로그인 안됨)
//   - *logger.Report: 에러 발생 시 에러 보고 객체
func (svc *LoginService) IsLogin(token string) (bool, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"token": %s}`, token))
	if token == "" {
		return false, logger.CreateReport(&logger.CODE_TOKEN_EXPIRED, nil)
	}

	// 토큰을 통한 사용자 조회
	if user, err := svc.dao.SelectByToken(svc.ctx, token); err != nil || user.Token == "" {
		return false, logger.CreateReport(&logger.CODE_TOKEN_EXPIRED, nil)
	}

	return true, nil
}

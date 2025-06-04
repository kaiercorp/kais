package router

import (
	"api_server/logger"
	repo "api_server/user/repository"
	"api_server/user/service"
	"api_server/utils"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	svc service.IUserService
}

var once sync.Once
var instance *UserController

func New(svc service.IUserService) *UserController {
	once.Do(func() { // atomic, does not allow repeating
		logger.Debug("User Controller instance")
		instance = &UserController{
			svc: svc,
		}
	})

	return instance
}

// CreateUser은 새로운 사용자 생성을 처리하는 HTTP POST 핸들러입니다.
//
// 요청 바디(JSON):
//   - name (string, required)
//   - username (string, required)
//   - password (string, required)
//   - group (int, required)
//   - is_use (bool, optional)
//
// 반환:
//   - 200 OK: 생성된 user_id 포함
//   - 400 Bad Request: 바인딩 실패 또는 서비스 오류
func (ctlr *UserController) CreateUser(c *gin.Context) {
	//todo: 중복 체크
	logger.ApiRequest(c)

	req := repo.UserDTO{}
	if err := c.ShouldBindJSON(&req); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	data, r := ctlr.svc.Create(req)
	logger.ApiResponse(c, r, gin.H{
		"user_id": data.ID,
	})
}

// UpdateUser는 기존 사용자의 정보를 수정하는 HTTP PUT 핸들러입니다.
//
// 요청 바디(JSON):
//   - id (int, required)
//   - name, username, password, group, is_use (optional)
//
// 반환:
//   - 200 OK: 수정된 user_id 포함
//   - 400 Bad Request: 바인딩 실패 또는 유효하지 않은 사용자
func (ctlr *UserController) UpdateUser(c *gin.Context) {
	logger.ApiRequest(c)

	req := repo.UserDTO{}
	if err := c.ShouldBindJSON(&req); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	data, r := ctlr.svc.Edit(req)
	logger.ApiResponse(c, r, gin.H{
		"user_id": data.ID,
	})
}

// GetUser는 사용자명을 기준으로 단일 사용자 정보를 조회하는 HTTP GET 핸들러입니다.
//
// URL Path 파라미터:
//   - username (string)
//
// 반환:
//   - 200 OK: 사용자 정보 JSON
//   - 404 Not Found: 사용자를 찾을 수 없는 경우
func (ctlr *UserController) GetUser(c *gin.Context) {
	logger.ApiRequest(c)

	username := c.Param("username")
	data, r := ctlr.svc.ReadByUsername(username)
	logger.ApiResponse(c, r, data)
}

// ChangeActivateStatus는 사용자의 활성화 상태(is_use)를 변경하는 HTTP PATCH 핸들러입니다.
//
// URL Path 파라미터:
//   - id (int)
//
// 요청 바디(JSON):
//   - is_use (bool, required)
//
// 반환:
//   - 200 OK: 상태 변경 성공
//   - 400 Bad Request: 파라미터 오류 또는 유효하지 않은 사용자
func (ctlr *UserController) ChangeActivateStatus(c *gin.Context) {
	logger.ApiRequest(c)

	req := repo.UserDTO{}
	if err := c.ShouldBindJSON(&req); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	_, r := ctlr.svc.Activate(id, *req.IsUse)
	logger.ApiResponse(c, r, nil)
}

// ResetPassword는 특정 사용자의 비밀번호를 초기화하는 HTTP PATCH 핸들러입니다.
//
// URL Path 파라미터:
//   - id (int)
//
// 반환:
//   - 200 OK: 초기화 성공
//   - 400/404: 사용자 없음 또는 서비스 오류
func (ctlr *UserController) ResetPassword(c *gin.Context) {
	logger.ApiRequest(c)

	id, _ := strconv.Atoi(c.Param("id"))
	_, r := ctlr.svc.ResetPassword(id)

	logger.ApiResponse(c, r, nil)
}

// DeleteUser는 특정 사용자를 삭제하는 HTTP DELETE 핸들러입니다.
//
// URL Path 파라미터:
//   - id (int)
//
// 반환:
//   - 200 OK: 삭제 성공
//   - 404 Not Found: 존재하지 않는 사용자
func (ctlr *UserController) DeleteUser(c *gin.Context) {
	logger.ApiRequest(c)

	id, _ := strconv.Atoi(c.Param("id"))
	r := ctlr.svc.Delete(id)
	logger.ApiResponse(c, r, nil)
}

// ChangePassword는 로그인한 사용자가 자신의 비밀번호를 변경하는 HTTP PATCH 핸들러입니다.
//
// 요청 바디(JSON):
//   - current_password (string, required)
//   - new_password (string, required)
//
// 인증:
//   - JWT 토큰에서 사용자 정보 추출
//
// 반환:
//   - 200 OK: 변경 성공
//   - 400 Bad Request: 인증 실패, 요청 형식 오류, 현재 비밀번호 불일치 등
func (ctlr *UserController) ChangePassword(c *gin.Context) {
	logger.ApiRequest(c)

	ctxData, err := utils.GetDataFromToken(c)
	if err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}
	type changePasswordReq struct {
		CurrentPassword string `json:"current_password,omitempty"`
		NewPassword     string `json:"new_password,omitempty"`
	}

	req := changePasswordReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	_, r := ctlr.svc.ChangePassword(ctxData.Username, req.CurrentPassword, req.NewPassword)
	logger.ApiResponse(c, r, nil)
}

// GetUsersByGroupGT는 주어진 그룹 값보다 큰 사용자들을 조회하는 HTTP GET 핸들러입니다.
//
// URL Path 파라미터:
//   - group (int)
//
// 반환:
//   - 200 OK: 조건에 해당하는 사용자 목록
//   - 400/500: 파싱 실패 또는 서버 오류
func (ctlr *UserController) GetUsersByGroupGT(c *gin.Context) {
	logger.ApiRequest(c)

	group, _ := strconv.Atoi(c.Param("group"))
	data, r := ctlr.svc.ReadManyByGroupGT(group)
	logger.ApiResponse(c, r, data)
}

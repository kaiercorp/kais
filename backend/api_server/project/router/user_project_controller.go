package router

import (
	"sync"

	"github.com/gin-gonic/gin"

	"api_server/logger"
	repo "api_server/project/repository"
	"api_server/project/service"
)

// UserProjectController는 사용자-프로젝트 관련 API 요청을 처리하는 컨트롤러입니다.
// 이 컨트롤러는 사용자-프로젝트 생성, 조회 및 삭제 등을 담당합니다.
type UserProjectController struct {
	svc service.IUserProjectService
}

var onceUsrProj sync.Once
var instanceUsrProj *UserProjectController

func NewUserProjectController(svc service.IUserProjectService) *UserProjectController {
	onceUsrProj.Do(func() { // atomic, does not allow repeating
		logger.Debug("User Project Controller instance")
		instanceUsrProj = &UserProjectController{
			svc: svc,
		}
	})

	return instanceUsrProj
}

// CreateUserProject는 새로운 사용자-프로젝트 매핑을 생성하는 API 핸들러입니다.
//
// 응답으로 생성된 데이터와 상태 보고서를 반환합니다.
func (ctlr *UserProjectController) CreateUserProject(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.UserProjectDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	data, report := ctlr.svc.Create(reqDTO)
	logger.ApiResponse(c, report, data)
}

// GetAllUsersByProject는 주어진 프로젝트 ID에 속한 모든 사용자를 조회하는 API 핸들러입니다.
//
// req 바디에서 프로젝트 ID를 받아 해당 프로젝트에 속한 사용자의 목록을 반환합니다.
func (ctlr *UserProjectController) GetAllUsersByProject(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.UserProjectDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	data, report := ctlr.svc.GetUsersByProjectId(reqDTO.ProjectId)
	logger.ApiResponse(c, report, data)
}

// GetAllProjectsByUser는 주어진 사용자명에 속한 모든 프로젝트 ID를 조회하는 API 핸들러입니다.
//
// req 바디에서 사용자명을 받아 해당 사용자에 속한 프로젝트 ID 목록을 반환합니다.
func (ctlr *UserProjectController) GetAllProjectsByUser(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.UserProjectDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	data, report := ctlr.svc.GetProjectIdsByUsername(reqDTO.Username)
	logger.ApiResponse(c, report, data)
}

// DeleteAllUsersByProject는 주어진 프로젝트에 속한 모든 사용자-프로젝트 매핑을 삭제하는 API 핸들러입니다.
//
// req 바디에서 프로젝트 ID를 받아 해당 프로젝트에 속한 모든 사용자-프로젝트 매핑을 삭제합니다.
func (ctlr *UserProjectController) DeleteAllUsersByProject(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.UserProjectDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	report := ctlr.svc.DeleteByProjectId(reqDTO.ProjectId)
	logger.ApiResponse(c, report, nil)
}

// DeleteAllProjectsByUser는 주어진 사용자에 속한 모든 사용자-프로젝트 매핑을 삭제하는 API 핸들러입니다.
//
// req 바디에서 사용자명을 받아 해당 사용자에 속한 모든 사용자-프로젝트 매핑을 삭제합니다.
func (ctlr *UserProjectController) DeleteAllProjectsByUser(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.UserProjectDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	report := ctlr.svc.DeleteByUsername(reqDTO.Username)
	logger.ApiResponse(c, report, nil)
}

// DeleteUserFromProject는 주어진 프로젝트 ID와 사용자명을 사용하여 특정 사용자-프로젝트 매핑을 삭제하는 API 핸들러입니다.
//
// req 바디에서 프로젝트 ID와 사용자명을 받아 특정 사용자-프로젝트 매핑을 삭제합니다.
func (ctlr *UserProjectController) DeleteUserFromProject(c *gin.Context) {
	logger.ApiRequest(c)

	reqDTO := repo.UserProjectDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	report := ctlr.svc.Delete(reqDTO.ProjectId, reqDTO.Username)
	logger.ApiResponse(c, report, nil)
}

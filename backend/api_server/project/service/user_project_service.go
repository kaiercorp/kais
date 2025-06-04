package service

import (
	"context"
	"fmt"
	"sync"

	"api_server/logger"
	repo "api_server/project/repository"
)

type IUserProjectService interface {
	// Create는 사용자-프로젝트 매핑을 생성합니다.
	//
	// 매개변수:
	//   - req: 생성할 사용자-프로젝트 매핑의 데이터 (UserProjectDTO)
	//
	// 반환값:
	//   - *UserProjectDTO: 생성된 사용자-프로젝트 매핑의 데이터 DTO
	//   - *logger.Report: 로깅 및 오류 보고서 객체
	Create(req repo.UserProjectDTO) (*repo.UserProjectDTO, *logger.Report)
	// GetUsersByProjectId는 주어진 프로젝트 ID에 속한 모든 사용자의 목록을 조회합니다.
	//
	// 매개변수:
	//   - project_id: 조회할 프로젝트의 ID
	//
	// 반환값:
	//   - []string: 해당 프로젝트에 속한 사용자의 사용자명 리스트
	//   - *logger.Report: 로깅 및 오류 보고서 객체
	GetUsersByProjectId(project_id int) ([]string, *logger.Report)
	// GetProjectIdsByUsername는 주어진 사용자명에 속한 모든 프로젝트 ID 목록을 조회합니다.
	//
	// 매개변수:
	//   - username: 조회할 사용자의 사용자명
	//
	// 반환값:
	//   - []int: 해당 사용자에 속한 프로젝트 ID 리스트
	//   - *logger.Report: 로깅 및 오류 보고서 객체
	GetProjectIdsByUsername(username string) ([]int, *logger.Report)
	// Delete는 특정 프로젝트와 사용자에 대한 매핑을 삭제합니다.
	//
	// 매개변수:
	//   - project_id: 삭제할 프로젝트의 ID
	//   - username: 삭제할 사용자의 사용자명
	//
	// 반환값:
	//   - *logger.Report: 로깅 및 오류 보고서 객체
	Delete(project_id int, username string) *logger.Report
	// DeleteByProjectId는 특정 프로젝트에 속한 모든 사용자-프로젝트 매핑을 삭제합니다.
	//
	// 매개변수:
	//   - project_id: 삭제할 프로젝트의 ID
	//
	// 반환값:
	//   - *logger.Report: 로깅 및 오류 보고서 객체
	DeleteByProjectId(project_id int) *logger.Report
	// DeleteByUsername은 특정 사용자에 대한 모든 사용자-프로젝트 매핑을 삭제합니다.
	//
	// 매개변수:
	//   - username: 삭제할 사용자의 사용자명
	//
	// 반환값:
	//   - *logger.Report: 로깅 및 오류 보고서 객체
	DeleteByUsername(username string) *logger.Report
}

// UserProjectService는 사용자와 프로젝트 간의 관계를 관리하는 서비스입니다.
// 사용자-프로젝트의 생성, 삭제 및 조회 등의 작업을 처리합니다.
type UserProjectService struct {
	ctx context.Context
	dao repo.IUserProjectDAO
}

var onceUserProject sync.Once
var instanceUserProject *UserProjectService

func NewUserProject(dao repo.IUserProjectDAO) *UserProjectService {
	onceUserProject.Do(func() { // atomic, does not allow repeating
		logger.Debug("User Project Service instance")
		instanceUserProject = &UserProjectService{
			ctx: context.Background(),
			dao: dao,
		}
	})

	return instanceUserProject
}

// Create는 사용자-프로젝트 매핑을 생성합니다.
//
// 매개변수:
//   - req: 생성할 사용자-프로젝트 매핑의 데이터 (UserProjectDTO)
//
// 반환값:
//   - *UserProjectDTO: 생성된 사용자-프로젝트 매핑의 데이터 DTO
//   - *logger.Report: 로깅 및 오류 보고서 객체
func (svc *UserProjectService) Create(req repo.UserProjectDTO) (*repo.UserProjectDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("%+v", req))
	// 사용자-프로젝트 매핑 생성
	project, err := svc.dao.InsertOne(svc.ctx, req)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_INSERT, err)
	}
	return repo.ConvertUserProjectEntToDTO(project), nil
}

// GetUsersByProjectId는 주어진 프로젝트 ID에 속한 모든 사용자의 목록을 조회합니다.
//
// 매개변수:
//   - project_id: 조회할 프로젝트의 ID
//
// 반환값:
//   - []string: 해당 프로젝트에 속한 사용자의 사용자명 리스트
//   - *logger.Report: 로깅 및 오류 보고서 객체
func (svc *UserProjectService) GetUsersByProjectId(project_id int) ([]string, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"project_id": %d}`, project_id))
	// 프로젝트에 속한 사용자 목록 조회
	ents, err := svc.dao.SelectByProjectID(svc.ctx, project_id)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}
	usernameList := make([]string, 0, len(ents))
	for _, ent := range ents {
		usernameList = append(usernameList, ent.Username)
	}
	return usernameList, nil
}

// GetProjectIdsByUsername는 주어진 사용자명에 속한 모든 프로젝트 ID 목록을 조회합니다.
//
// 매개변수:
//   - username: 조회할 사용자의 사용자명
//
// 반환값:
//   - []int: 해당 사용자에 속한 프로젝트 ID 리스트
//   - *logger.Report: 로깅 및 오류 보고서 객체
func (svc *UserProjectService) GetProjectIdsByUsername(username string) ([]int, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"username": %s}`, username))
	// 사용자에 속한 프로젝트 목록 조회
	ents, err := svc.dao.SelectByUsername(svc.ctx, username)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}
	projectIdList := make([]int, 0, len(ents))
	for _, ent := range ents {
		projectIdList = append(projectIdList, ent.ProjectID)
	}
	return projectIdList, nil
}

// Delete는 특정 프로젝트와 사용자에 대한 매핑을 삭제합니다.
//
// 매개변수:
//   - project_id: 삭제할 프로젝트의 ID
//   - username: 삭제할 사용자의 사용자명
//
// 반환값:
//   - *logger.Report: 로깅 및 오류 보고서 객체
func (svc *UserProjectService) Delete(project_id int, username string) *logger.Report {
	logger.Debug(fmt.Sprintf(`{"project_id": %d, "username": "%s"}`, project_id, username))

	// 프로젝트와 사용자 매핑 삭제
	if _, err := svc.dao.DeleteOne(svc.ctx, project_id, username); err != nil {
		return logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}
	return nil
}

// DeleteByProjectId는 특정 프로젝트에 속한 모든 사용자-프로젝트 매핑을 삭제합니다.
//
// 매개변수:
//   - project_id: 삭제할 프로젝트의 ID
//
// 반환값:
//   - *logger.Report: 로깅 및 오류 보고서 객체
func (svc *UserProjectService) DeleteByProjectId(project_id int) *logger.Report {
	logger.Debug(fmt.Sprintf(`{"project_id": %d}`, project_id))

	// 프로젝트 ID에 해당하는 모든 사용자-프로젝트 매핑 삭제
	if _, err := svc.dao.DeleteByProjectId(svc.ctx, project_id); err != nil {
		return logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}
	return nil
}

// DeleteByUsername은 특정 사용자에 대한 모든 사용자-프로젝트 매핑을 삭제합니다.
//
// 매개변수:
//   - username: 삭제할 사용자의 사용자명
//
// 반환값:
//   - *logger.Report: 로깅 및 오류 보고서 객체
func (svc *UserProjectService) DeleteByUsername(username string) *logger.Report {
	logger.Debug(fmt.Sprintf(`{"username": %s}`, username))

	// 사용자명에 해당하는 모든 사용자-프로젝트 매핑 삭제
	if _, err := svc.dao.DeleteByUsername(svc.ctx, username); err != nil {
		return logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}
	return nil
}

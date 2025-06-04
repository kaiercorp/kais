package service

import (
	"api_server/logger"
	"api_server/project/csvformat"
	repo "api_server/project/repository"
	svc_task "api_server/task/service"
	"api_server/utils"
	"context"
	"fmt"
	"sync"
)

type IProjectService interface {

	// Create는 새 프로젝트를 생성하고, 사용자-프로젝트 매핑을 추가합니다.
	//
	// 매개변수:
	//   - req: 생성할 프로젝트의 데이터 (ProjectDTO)
	//   - username: 사용자-프로젝트 매핑 생성에 필요한 사용자이름
	Create(req repo.ProjectDTO, username string) (*repo.ProjectPages, *logger.Report)
	// Read는 페이지 번호를 기준으로 프로젝트 목록을 조회합니다.
	Read(page int) (*repo.ProjectPages, *logger.Report)
	// Edit는 기존 프로젝트의 정보를 수정합니다.
	Edit(req repo.ProjectDTO) (*repo.ProjectDTO, *logger.Report)
	// Delete는 특정 프로젝트를 삭제합니다. 이때 관련된 작업과 사용자-프로젝트 매핑도 함께 삭제됩니다.
	Delete(project_id int) (*repo.ProjectPages, *logger.Report)
	// ExportProjectTableToCSV는 주어진 projectIDs에 해당하는 프로젝트 데이터를 CSV로 내보내는 서비스 함수입니다.
	//
	// 매개변수:
	//   - projectIDs: 내보낼 프로젝트의 ID 목록
	//
	// 반환 값:
	//   - string: 생성된 CSV 파일의 경로
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ExportProjectTableToCSV(projectIDs []int) (string, *logger.Report)
	// ImportProjectTableFromCSV는 CSV 파일에서 프로젝트 데이터를 가져와 데이터베이스에 삽입하는 서비스 함수입니다.
	//
	// 매개변수:
	//   - filename: 가져올 CSV 파일의 경로
	//   - username: 프로젝트에 연관된 사용자 이름
	//
	// 반환 값:
	//   - map[int]int: CSV에서 읽은 프로젝트 ID와 프로젝트 데이터를 매핑하는 맵
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ImportProjectTableFromCSV(filename string, username string) (map[int]int, *logger.Report)
	ReadByUsername(page int, username string) (*repo.ProjectPages, *logger.Report)
}

type ProjectService struct {
	ctx          context.Context
	dao          repo.IProjectDAO
	task_svc     svc_task.ITaskService
	usr_proj_svc IUserProjectService
}

var once sync.Once
var instance *ProjectService

func New(dao repo.IProjectDAO, task_svc svc_task.ITaskService, user_project_svc IUserProjectService) *ProjectService {
	once.Do(func() { // atomic, does not allow repeating
		logger.Debug("Project Service instance")
		instance = &ProjectService{
			ctx:          context.Background(),
			dao:          dao,
			usr_proj_svc: user_project_svc,
			task_svc:     task_svc,
		}
	})

	return instance
}

// Create는 새 프로젝트를 생성하고, 사용자-프로젝트 매핑을 추가합니다.
func (svc *ProjectService) Create(req repo.ProjectDTO, username string) (*repo.ProjectPages, *logger.Report) {
	logger.Debug(fmt.Sprintf("%+v", req))
	// 프로젝트 생성
	project, err := svc.dao.InsertOne(svc.ctx, req)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_INSERT, err)
	}
	// 사용자-프로젝트 매핑 생성
	req_usr_proj := repo.UserProjectDTO{
		ProjectId: project.ID,
		Username:  username,
	}
	svc.usr_proj_svc.Create(req_usr_proj)

	return svc.Read(1)
}

// Read는 페이지 번호를 기준으로 프로젝트 목록을 조회합니다.
func (svc *ProjectService) Read(page int) (*repo.ProjectPages, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"page": %d}`, page))
	if projectList, pageCount, hasMore, nextPage, err := svc.dao.SelectByPage(svc.ctx, page); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertEntsToProjectPages(projectList, pageCount, hasMore, nextPage), nil
	}
}

// Edit는 기존 프로젝트의 정보를 수정합니다.
func (svc *ProjectService) Edit(req repo.ProjectDTO) (*repo.ProjectDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("%+v", req))
	// 프로젝트 수정
	if project, err := svc.dao.UpdateOne(svc.ctx, req); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	} else {
		return repo.ConvertEntToDTO(project), nil
	}
}

// Delete는 특정 프로젝트를 삭제합니다. 이때 관련된 작업과 사용자-프로젝트 매핑도 함께 삭제됩니다.
func (svc *ProjectService) Delete(project_id int) (*repo.ProjectPages, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"project_id": %d}`, project_id))

	// 프로젝트에 관련된 작업 삭제
	if r := svc.task_svc.DeleteByProject(project_id); r != nil {
		return nil, r
	}

	// 사용자-프로젝트 매핑 삭제
	if r := svc.usr_proj_svc.DeleteByProjectId(project_id); r != nil {
		return nil, r
	}

	// 프로젝트 삭제
	if err := svc.dao.DeleteOne(svc.ctx, project_id); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}

	// 삭제 후 프로젝트 목록 조회
	return svc.Read(1)
}

// ExportProjectTableToCSV는 주어진 projectIDs에 해당하는 프로젝트 데이터를 CSV로 내보내는 서비스 함수입니다.
func (svc *ProjectService) ExportProjectTableToCSV(projectIDs []int) (string, *logger.Report) {
	// 디버깅을 위한 로그 출력
	logger.Debug(fmt.Sprintf(`{"project_ids": %v}`, projectIDs))

	// 프로젝트 데이터를 SelectMany 메서드를 호출하여 조회
	projects, err := svc.dao.SelectMany(svc.ctx, projectIDs)
	if err != nil {
		// 오류가 발생한 경우 에러 리포트 반환
		return "", logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	// ProjectCSVFormat 인스턴스를 생성하여 CSV 포맷을 정의
	projectCSVFormat := csvformat.ProjectCSVFormat{}
	// CSV 헤더 정의
	header := projectCSVFormat.GetHeader()

	// 프로젝트 데이터를 CSV 레코드로 변환
	var records [][]string
	for _, project := range projects {
		// 각 프로젝트를 CSV 레코드로 변환하여 추가
		records = append(records, projectCSVFormat.ConvertToRecord(project))
	}

	// utils 패키지의 ExportToCSV 함수 호출하여 CSV 파일로 저장
	return utils.ExportToCSV("project.csv", header, records)
}

// ImportProjectTableFromCSV는 CSV 파일에서 프로젝트 데이터를 가져와 데이터베이스에 삽입하는 서비스 함수입니다.
func (svc *ProjectService) ImportProjectTableFromCSV(filename string, username string) (map[int]int, *logger.Report) {
	// ProjectCSVFormat 인스턴스를 생성하여 CSV 레코드 변환 및 파싱 작업 수행
	projectCSVFormat := csvformat.ProjectCSVFormat{}
	// CSV 레코드에서 ProjectDTO 객체로 변환하는 함수 정의
	importFunc := func(record []string) (*repo.ProjectDTO, int, error) {
		// CSV 레코드를 ProjectDTO로 변환
		recordProject, _ := projectCSVFormat.ParseRecord(record)

		// 변환된 ProjectDTO 객체와 해당 ID 반환
		return recordProject, recordProject.ID, nil
	}

	// 새 프로젝트 데이터를 데이터베이스에 삽입하는 함수 정의
	insertFunc := func(project *repo.ProjectDTO) (int, error) {
		// 프로젝트 삽입
		newProject, err := svc.dao.InsertOne(svc.ctx, *project)
		if err != nil {
			// 오류가 발생한 경우 에러 반환
			return 0, err
		}

		// 프로젝트에 사용자 연관 정보 삽입
		req_usr_proj := repo.UserProjectDTO{
			ProjectId: newProject.ID,
			Username:  username,
		}
		// 사용자-프로젝트 연결 생성
		_, r := svc.usr_proj_svc.Create(req_usr_proj)
		if r != nil {
			// 사용자-프로젝트 연결에 실패한 경우 에러 반환
			return 0, fmt.Errorf("failed to insert user-project")
		}

		// 삽입된 프로젝트의 ID 반환
		return newProject.ID, nil
	}

	// utils 패키지의 ImportFromCSV 함수 호출하여 CSV 파일로부터 데이터를 가져와 삽입
	return utils.ImportFromCSV(filename, importFunc, insertFunc)
}

func (svc *ProjectService) ReadByUsername(page int, username string) (*repo.ProjectPages, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"page": %d}`, page))
	ids, r := svc.usr_proj_svc.GetProjectIdsByUsername(username)
	if r != nil {
		return nil, r
	}
	projectList, pageCount, hasMore, nextPage, err := svc.dao.SelectByPageWithIDs(svc.ctx, page, ids)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}
	return repo.ConvertEntsToProjectPages(projectList, pageCount, hasMore, nextPage), nil
}

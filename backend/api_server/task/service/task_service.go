package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	config_service "api_server/configuration/service"
	repo_dataset "api_server/dataset/repository"
	repo_device "api_server/device/repository"
	"api_server/ent"
	"api_server/logger"
	"api_server/task/csvformat"
	repo "api_server/task/repository"
	"api_server/utils"
)

type ITaskService interface {
	Create(req repo.TaskDTO) (*repo.TaskDTO, *repo.ModelingDTO, *logger.Report)
	ReadByProject(project_id int) (*repo.TaskPages, *logger.Report)
	ReadOne(id int) (*repo.TaskDTO, *logger.Report)
	Edit(req repo.TaskDTO) (*repo.TaskDTO, *logger.Report)
	DeleteOne(id int) *logger.Report
	DeleteByProject(project_id int) *logger.Report
	makeModelingDTO(task *ent.Task) (*repo.ModelingDTO, *logger.Report)
	makeModelingParams(task *ent.Task) ([]string, *logger.Report)
	makeCommonParams(task *ent.Task, params map[string]interface{}) *logger.Report
	makeParamsVCLSSL(arams map[string]interface{}) *logger.Report
	makeParamsVCLSML(arams map[string]interface{}) *logger.Report
	makeParamsVAD(arams map[string]interface{}) *logger.Report
	makeParamsTCLS(arams map[string]interface{}) *logger.Report
	makeParamsTREG(arams map[string]interface{}) *logger.Report

	// ReadIDsByProjectIDs는 주어진 projectIDs에 해당하는 Task의 ID 목록을 조회하는 서비스 함수입니다.
	//
	// 매개변수:
	//   - projectIDs: 조회할 프로젝트의 ID 목록
	//
	// 반환 값:
	//   - []int: 조회된 Task의 ID 목록
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ReadIDsByProjectIDs(projectIDs []int) ([]int, *logger.Report)
	// ExportTaskTableToCSV는 주어진 taskIDs에 해당하는 Task 데이터를 CSV 파일로 내보내는 서비스 함수입니다.
	//
	// 매개변수:
	//   - taskIDs: 내보낼 Task의 ID 목록
	//
	// 반환 값:
	//   - string: 생성된 CSV 파일의 경로
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ExportTaskTableToCSV(taskIDs []int) (string, *logger.Report)
	// ImportTaskTableFromCSV는 CSV 파일에서 Task 데이터를 가져와 데이터베이스에 삽입하는 서비스 함수입니다.
	//
	// 매개변수:
	//   - filename: 가져올 CSV 파일의 경로
	//   - projectIdMap: CSV에서 읽은 프로젝트 ID를 실제 프로젝트 ID로 변환하기 위한 맵
	//
	// 반환 값:
	//   - map[int]int: CSV에서 읽은 Task의 ID와 프로젝트 ID가 매핑된 맵
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ImportTaskTableFromCSV(filename string, projectIdMap map[int]int) (map[int]int, *logger.Report)

	// ReadAll 함수는 전체 태스크(Task) 목록을 조회합니다.
	//
	// 반환값:
	//   - []*repo.TaskDTO: 태스크 DTO 포인터 리스트
	//   - *logger.Report: 오류 발생 시 반환되는 리포트 객체
	ReadAll() ([]*repo.TaskDTO, *logger.Report)

	// ReadByEngineType 함수는 주어진 엔진 타입에 해당하는 태스크 목록을 조회합니다.
	//
	// 매개변수:
	//   - engine_type: 조회할 엔진 타입 문자열
	//
	// 반환값:
	//   - []*repo.TaskDTO: 해당 엔진 타입에 해당하는 태스크 DTO 포인터 리스트
	//   - *logger.Report: 오류 발생 시 반환되는 리포트 객체
	ReadByEngineType(engine_type string) ([]*repo.TaskDTO, *logger.Report)
}

type TaskService struct {
	ctx          context.Context
	dao          repo.ITaskDAO
	dao_modeling repo.IModelingDAO
	dao_dataset  repo_dataset.DatasetDAOInterface
	dao_device   repo_device.IDeviceDAO
}

var onceTask sync.Once
var instanceTask *TaskService

func NewTaskService(dao repo.ITaskDAO,
	dao_modeling repo.IModelingDAO,
	dao_dataset repo_dataset.DatasetDAOInterface,
	dao_device repo_device.IDeviceDAO,
) *TaskService {
	onceTask.Do(func() {
		logger.Debug("Task Service instance")
		instanceTask = &TaskService{
			ctx:          context.Background(),
			dao:          dao,
			dao_modeling: dao_modeling,
			dao_dataset:  dao_dataset,
			dao_device:   dao_device,
		}
	})

	return instanceTask
}
func (svc *TaskService) Create(req repo.TaskDTO) (*repo.TaskDTO, *repo.ModelingDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("%+v", req))

	req.UserParams["engine_type"] = req.EngineType

	if jsonstr, err := json.Marshal(req.UserParams); err != nil {
		return nil, nil, logger.CreateReport(&logger.CODE_JSON_MARSHAL, err)
	} else {
		req.Params = []string{string(jsonstr)}
	}

	task, err := svc.dao.InsertOne(svc.ctx, req)
	if err != nil {
		return nil, nil, logger.CreateReport(&logger.CODE_DB_INSERT, err)
	}
	modeling, r := svc.makeModelingDTO(task)
	if r != nil {
		return nil, nil, r
	}
	m, err := svc.dao_modeling.InsertOne(svc.ctx, *modeling)
	if err != nil {
		return nil, nil, logger.CreateReport(&logger.CODE_DB_INSERT, err)
	}
	return repo.ConvertTaskEntToDTO(task), repo.ConvertModelingEntToDTO(m), nil
}

func (svc *TaskService) ReadByProject(project_id int) (
	*repo.TaskPages, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"project_id": %d}`, project_id))
	if taskList, pageCount, err := svc.dao.SelectByProject(svc.ctx, project_id); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertTaskEntsToTaskPages(taskList, pageCount), nil
	}
}

func (svc *TaskService) ReadOne(id int) (*repo.TaskDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	if task, err := svc.dao.SelectOne(svc.ctx, id); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertTaskEntToDTO(task), nil
	}
}

func (svc *TaskService) Edit(req repo.TaskDTO) (*repo.TaskDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("%+v", req))
	if task, err := svc.dao.UpdateOne(svc.ctx, req); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	} else {
		return repo.ConvertTaskEntToDTO(task), nil
	}
}

func (svc *TaskService) DeleteOne(id int) *logger.Report {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))

	modelings, err := svc.dao_modeling.SelectByTask(svc.ctx, id)
	if err != nil {
		return logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	for _, modeling := range modelings {
		if modeling.ModelingStep != utils.MODELING_STEP_COMPLETE && modeling.ModelingStep != utils.MODELING_STEP_FAIL && modeling.ModelingStep != utils.MODELING_STEP_CANCEL {
			return logger.CreateReport(&logger.CODE_MODELING_IN_PROGRESS, nil)
		}
	}

	if task, err := svc.dao.SelectOne(svc.ctx, id); err != nil {
		return logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		cf := config_service.NewStatic()
		targetPath := cf.Get("ROOT_PATH") + "/task/" + task.EngineType + "/" + strconv.Itoa(task.ID) + "/"

		// modeling 단위로 여러 머신에 흩어져 있을 수 있다.
		// 모든 engine들을 호출해서 remove 해야 한다.
		if devices, err := svc.dao_device.SelectActive(svc.ctx); err != nil {
			return logger.CreateReport(&logger.CODE_DB_SELECT, err)
		} else if len(devices) < 1 {
			return logger.CreateReport(&logger.CODE_DB_SELECT, errors.New("there are no active devices"))
		} else {
			for _, device := range devices {
				addr := "http://" + device.IP + ":" + strconv.Itoa(device.Port) + "/api/remove"
				params := bytes.NewReader([]byte(fmt.Sprintf(`{"target_path":"%s"}`, targetPath)))
				logger.Debug(addr, string(fmt.Sprintf(`{"target_path":%s}`, targetPath)))
				if resp, err := http.Post(addr, "application/json", params); err != nil {
					logger.Debug(&logger.CODE_REMOTE_SELECT, err)
				} else {
					logger.Debug(resp)
				}
			}
		}
	}

	if err := svc.dao.DeleteOne(svc.ctx, id); err != nil {
		return logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}

	return nil
}

func (svc *TaskService) DeleteByProject(project_id int) *logger.Report {
	logger.Debug(fmt.Sprintf(`{"project_id": %d}`, project_id))

	if tasks, _, err := svc.dao.SelectByProject(svc.ctx, project_id); len(tasks) == 0 {
		return nil
	} else if err != nil {
		logger.Debug(&logger.CODE_REMOTE_SELECT, err)
	} else {
		for _, task := range tasks {
			svc.DeleteOne(task.ID)
		}
	}

	if err := svc.dao.DeleteByProject(svc.ctx, project_id); err != nil {
		return logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}

	return nil
}

func (svc *TaskService) makeModelingDTO(task *ent.Task) (*repo.ModelingDTO, *logger.Report) {
	if params, err := svc.makeModelingParams(task); err != nil {
		return nil, err
	} else {
		modeling := &repo.ModelingDTO{
			Params:       params,
			ModelingType: utils.MODELING_TYPE_INITIAL,
			ModelingStep: utils.MODELING_STEP_IDLE,
			TaskID:       task.ID,
		}

		return modeling, nil
	}
}

func (svc *TaskService) makeModelingParams(task *ent.Task) ([]string, *logger.Report) {
	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(task.Params[0]), &params); err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	}

	if err := svc.makeCommonParams(task, params); err != nil {
		return nil, err
	}

	var r *logger.Report
	switch task.EngineType {
	case utils.JOB_TYPE_VISION_CLS_SL:
		r = svc.makeParamsVCLSSL(params)
	case utils.JOB_TYPE_VISION_CLS_ML:
		r = svc.makeParamsVCLSML(params)
	case utils.JOB_TYPE_VISION_AD:
		r = svc.makeParamsVAD(params)
	case utils.JOB_TYPE_TABLE_CLS:
		r = svc.makeParamsTCLS(params)
		// fallthrough
	case utils.JOB_TYPE_TABLE_REG:
		r = svc.makeParamsTREG(params)
	}

	if r != nil {
		return nil, r
	}

	if paramstr, err := json.Marshal(params); err != nil {
		return []string{}, logger.CreateReport(&logger.CODE_JSON_MARSHAL, err)
	} else {
		return []string{string(paramstr)}, nil
	}
}

func (svc *TaskService) makeCommonParams(task *ent.Task, params map[string]interface{}) *logger.Report {
	cf := config_service.NewStatic()

	params["data_path"], _ = svc.dao_dataset.SelectDataPathByDataSetId(svc.ctx, task.DatasetID)
	params["save_path"] = cf.Get("ROOT_PATH") + "/task/" + task.EngineType + "/" + strconv.Itoa(task.ID) + "/"
	params["target_metric"] = task.TargetMetric
	params["device_ids"] = params["gpus"]
	delete(params, "gpus")
	params["engine_type"] = task.EngineType

	return nil
}

func (svc *TaskService) makeParamsVCLSSL(params map[string]interface{}) *logger.Report {
	params["multi_label"] = false
	params["img_height"] = params["image_resolution"]
	params["img_width"] = params["image_resolution"]
	delete(params, "image_resolution")

	return nil
}

func (svc *TaskService) makeParamsVCLSML(params map[string]interface{}) *logger.Report {
	params["multi_label"] = true
	params["img_height"] = params["image_resolution"]
	params["img_width"] = params["image_resolution"]
	delete(params, "image_resolution")

	return nil
}

func (svc *TaskService) makeParamsVAD(params map[string]interface{}) *logger.Report {
	params["img_height"] = params["image_resolution"]
	params["img_width"] = params["image_resolution"]
	delete(params, "image_resolution")

	return nil
}

func (svc *TaskService) makeParamsTCLS(params map[string]interface{}) *logger.Report {
	params["task_mode"] = "classification"

	return nil
}

func (svc *TaskService) makeParamsTREG(params map[string]interface{}) *logger.Report {
	params["task_mode"] = "regression"

	return nil
}

// ReadIDsByProjectIDs는 주어진 projectIDs에 해당하는 Task의 ID 목록을 조회하는 서비스 함수입니다.
func (svc *TaskService) ReadIDsByProjectIDs(projectIDs []int) ([]int, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"project_ids": %v}`, projectIDs))

	// TaskDAO의 SelectIDsByProjectIDs 메서드를 호출하여 Task의 ID 목록을 조회
	taskIDs, err := svc.dao.SelectIDsByProjectIDs(svc.ctx, projectIDs)
	if err != nil {
		// 오류가 발생한 경우 에러 리포트 반환
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	// 성공적으로 조회된 ID 목록 반환
	return taskIDs, nil
}

// ExportTaskTableToCSV는 주어진 taskIDs에 해당하는 Task 데이터를 CSV 파일로 내보내는 서비스 함수입니다.
func (svc *TaskService) ExportTaskTableToCSV(taskIDs []int) (string, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"task_ids": %v}`, taskIDs))

	// TaskDAO의 SelectMany 메서드를 호출하여 Task 데이터 조회
	tasks, err := svc.dao.SelectMany(svc.ctx, taskIDs)
	if err != nil {
		// 오류가 발생한 경우 에러 리포트 반환
		return "", logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	// CSV 포맷 정의
	taskCSVFormat := csvformat.TaskCSVFormat{}
	// CSV 헤더 정의
	header := taskCSVFormat.GetHeader()
	// CSV 데이터로 변환
	var records [][]string
	for _, task := range tasks {
		// 각 Task를 CSV 레코드로 변환하여 추가
		records = append(records, taskCSVFormat.ConvertToRecord(task))
	}

	// utils 패키지의 ExportToCSV 함수 호출하여 CSV 파일로 저장
	return utils.ExportToCSV("task.csv", header, records)
}

// ImportTaskTableFromCSV는 CSV 파일에서 Task 데이터를 가져와 데이터베이스에 삽입하는 서비스 함수입니다.
func (svc *TaskService) ImportTaskTableFromCSV(filename string, projectIdMap map[int]int) (map[int]int, *logger.Report) {
	// TaskCSVFormat 인스턴스를 생성하여 CSV 레코드 변환 및 파싱 작업 수행
	taskCSVFormat := csvformat.TaskCSVFormat{}

	// CSV 레코드를 TaskDTO로 변환하는 함수
	importFunc := func(record []string) (*repo.TaskDTO, int, error) {
		// CSV 레코드에서 TaskDTO 객체로 변환
		recordTask, _ := taskCSVFormat.ParseRecord(record)
		// 프로젝트 ID 변환
		*recordTask.ProjectID = projectIdMap[*recordTask.ProjectID]

		// 변환된 TaskDTO 반환
		return recordTask, recordTask.ID, nil
	}

	// 새 Task 데이터를 데이터베이스에 삽입하는 함수
	insertFunc := func(task *repo.TaskDTO) (int, error) {
		// TaskDAO의 InsertOne 메서드를 호출하여 Task 삽입
		newTask, err := svc.dao.InsertOne(svc.ctx, *task)
		if err != nil {
			// 오류가 발생한 경우 에러 반환
			return 0, err
		}
		// 삽입된 Task의 ID 반환
		return newTask.ID, nil
	}

	// utils 패키지의 ImportFromCSV 함수 호출하여 CSV 파일로부터 데이터를 가져와 삽입
	return utils.ImportFromCSV(filename, importFunc, insertFunc)
}

// ReadAll 함수는 전체 태스크(Task) 목록을 조회합니다.
func (svc *TaskService) ReadAll() ([]*repo.TaskDTO, *logger.Report) {
	logger.Debug("ReadAll()")
	if tasks, err := svc.dao.SelectAll(svc.ctx); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertTaskEntsToDTOs(tasks), nil
	}
}

// ReadByEngineType 함수는 주어진 엔진 타입에 해당하는 태스크 목록을 조회합니다.
func (svc *TaskService) ReadByEngineType(engine_type string) ([]*repo.TaskDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"project_id": %s}`, engine_type))
	if tasks, err := svc.dao.SelectByEngineType(svc.ctx, engine_type); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertTaskEntsToDTOs(tasks), nil
	}
}

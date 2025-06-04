package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	config_service "api_server/configuration/service"
	repo_dataset "api_server/dataset/repository"
	repo_device "api_server/device/repository"
	"api_server/logger"
	"api_server/task/csvformat"
	repo "api_server/task/repository"
	"api_server/utils"
)

type IModelingService interface {
	CreateEvaluation(req repo.EvaluationDTO) (*repo.ModelingDTO, *logger.Report)
	ReadByTask(task_id int) ([]*repo.ModelingDB, *logger.Report)
	ReadModelingType(task_id int) ([]*repo.ModelingDB, *logger.Report)
	ReadOne(id int) (*repo.ModelingDTO, *logger.Report)
	ReadFull(id int) (*repo.ModelingDB, *logger.Report)
	StopModelingTask(id int) *logger.Report
	makeModelingDTO(req repo.EvaluationDTO) (*repo.ModelingDTO, *logger.Report)
	makeModelingParams(req repo.EvaluationDTO, parent *repo.ModelingDTO) ([]string, *logger.Report)
	makeCommonParams(req repo.EvaluationDTO, params map[string]interface{}, parentParams map[string]interface{}) *logger.Report
	makeParamsVCLSSL(arams map[string]interface{}) *logger.Report
	makeParamsVCLSML(arams map[string]interface{}) *logger.Report
	makeParamsVAD(arams map[string]interface{}) *logger.Report
	makeParamsTCLS(arams map[string]interface{}) *logger.Report
	makeParamsTREG(arams map[string]interface{}) *logger.Report
	DeleteOne(id int) *logger.Report

	// ReadIDsByTaskIDs는 주어진 taskIDs에 해당하는 Task의 ID 목록을 조회하는 서비스 함수입니다.
	//
	// 매개변수:
	//   - taskIDs: 조회할 Task의 ID 목록
	//
	// 반환 값:
	//   - []int: 조회된 Task의 ID 목록
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ReadIDsByTaskIDs(taskIDs []int) ([]int, *logger.Report)
	// ReadModelingModelsIDsByModelingIDs는 주어진 modelingIDs에 해당하는 ModelingModel의 ID 목록을 조회하는 서비스 함수입니다.
	//
	// 매개변수:
	//   - modelingIDs: 조회할 Modeling의 ID 목록
	//
	// 반환 값:
	//   - []int: 조회된 ModelingModel의 ID 목록
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ReadModelingModelsIDsByModelingIDs(modelingIDs []int) ([]int, *logger.Report)
	// ReadModelingDetailsIDsByModelingIDs는 주어진 modelingIDs에 해당하는 ModelingDetail의 ID 목록을 조회하는 서비스 함수입니다.
	//
	// 매개변수:
	//   - modelingIDs: 조회할 Modeling의 ID 목록
	//
	// 반환 값:
	//   - []int: 조회된 ModelingDetail의 ID 목록
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ReadModelingDetailsIDsByModelingIDs(modelingIDs []int) ([]int, *logger.Report)
	// ExportModelingTableToCSV는 주어진 modelingIDs에 해당하는 Modeling 데이터를 CSV로 내보내는 서비스 함수입니다.
	//
	// 매개변수:
	//   - modelingIDs: 내보낼 Modeling의 ID 목록
	//
	// 반환 값:
	//   - string: 생성된 CSV 파일의 경로
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ExportModelingTableToCSV(modelingIDs []int) (string, *logger.Report)
	// ImportModelingTableFromCSV는 CSV 파일에서 Modeling 데이터를 가져와 데이터베이스에 삽입하는 서비스 함수입니다.
	//
	// 매개변수:
	//   - filename: 가져올 CSV 파일의 경로
	//   - taskIdMap: CSV에서 읽은 Task ID를 실제 새로 변환된 Task ID로 변환하기 위한 맵.
	//
	// 반환 값:
	//   - map[int]int: CSV에서 읽은 Modeling의 ID와 Task ID가 매핑된 맵
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ImportModelingTableFromCSV(filename string, taskIdMap map[int]int) (map[int]int, *logger.Report)
	// ExportModelingDetailsTableToCSV는 주어진 modelingDetailIDs에 해당하는 ModelingDetail 데이터를 CSV로 내보내는 서비스 함수입니다.
	//
	// 매개변수:
	//   - modelingDetailIDs: 내보낼 ModelingDetail의 ID 목록
	//
	// 반환 값:
	//   - string: 생성된 CSV 파일의 경로
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ExportModelingDetailsTableToCSV(modelingDetailIDs []int) (string, *logger.Report)
	// ImportModelingDetailsTableFromCSV는 CSV 파일에서 ModelingDetail 데이터를 가져와 데이터베이스에 삽입하는 서비스 함수입니다.
	//
	// 매개변수:
	//   - filename: 가져올 CSV 파일의 경로
	//   - modelingIdMap: CSV에서 읽은 Modeling ID를 실제 Modeling ID로 변환하기 위한 맵
	//
	// 반환 값:
	//   - map[int]int: CSV에서 읽은 ModelingDetail의 ID와 Modeling ID가 매핑된 맵
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ImportModelingDetailsTableFromCSV(filename string, modelingIdMap map[int]int) (map[int]int, *logger.Report)
	// ExportModelingModelsTableToCSV는 주어진 modelingModelIDs에 해당하는 ModelingModel 데이터를 CSV로 내보내는 서비스 함수입니다.
	//
	// 매개변수:
	//   - modelingModelIDs: 내보낼 ModelingModel의 ID 목록
	//
	// 반환 값:
	//   - string: 생성된 CSV 파일의 경로
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ExportModelingModelsTableToCSV(modelingModelIDs []int) (string, *logger.Report)
	// ImportModelingModelsTableFromCSV는 CSV 파일에서 ModelingModel 데이터를 가져와 데이터베이스에 삽입하는 서비스 함수입니다.
	//
	// 매개변수:
	//   - filename: 가져올 CSV 파일의 경로
	//   - modelingIdMap: CSV에서 읽은 Modeling ID를 실제 Modeling ID로 변환하기 위한 맵
	//
	// 반환 값:
	//   - map[int]int: CSV에서 읽은 ModelingModel의 ID와 Modeling ID가 매핑된 맵
	//   - *logger.Report: 오류가 발생한 경우 에러 정보를 담고 있는 리포트
	ImportModelingModelsTableFromCSV(filename string, modelingIdMap map[int]int) (map[int]int, *logger.Report)

	// GetBestModelNameByModelingID 함수는 특정 모델링 ID에 대한 최적의 모델 이름을 조회합니다.
	//
	// 매개변수:
	//   - modelingID: 조회할 모델링 ID
	//
	// 반환값:
	//   - string: 최적 모델의 이름
	//   - *logger.Report: 오류 발생 시 반환되는 리포트 객체
	GetBestModelNameByModelingID(modelingID int) (string, *logger.Report)

	// GetBestModelNameListByModelingID 함수는 특정 모델링 ID에 대한 모든 최적 모델 이름 리스트를 반환합니다.
	//
	// 매개변수:
	//   - modelingID: 조회할 모델링 ID
	//
	// 반환값:
	//   - []string: 모델 이름 리스트
	//   - *logger.Report: 오류 발생 시 반환되는 리포트 객체
	GetBestModelNameListByModelingID(modelingID int) ([]string, *logger.Report)

	// ParseModelingDetailDataByType 함수는 특정 모델링 ID와 모델 이름에 해당하는 상세 데이터를 타입별로 파싱합니다.
	//
	// 매개변수:
	//   - modelingID: 모델링 ID
	//   - modelName: 모델 이름
	//
	// 반환값:
	//   - map[string]interface{}: 타입별 파싱된 상세 데이터 맵
	//   - *logger.Report: 오류 발생 시 반환되는 리포트 객체
	ParseModelingDetailDataByType(modelingID int, modelName string) (map[string]interface{}, *logger.Report)

	// ReadModelingModelsByTypeAndModelingID 함수는 특정 모델링 ID와 데이터 타입에 해당하는 모델 리스트를 조회합니다.
	//
	// 매개변수:
	//   - modelingID: 모델링 ID
	//   - data_type: 조회할 데이터 타입
	//
	// 반환값:
	//   - *repo.ModelingModels: 조회된 모델 리스트 DTO
	//   - *logger.Report: 오류 발생 시 반환되는 리포트 객체
	ReadModelingModelsByTypeAndModelingID(modelingID int, data_type string) (*repo.ModelingModels, *logger.Report)
}

type ModelingService struct {
	ctx         context.Context
	dao         repo.IModelingDAO
	dao_details repo.IModelinDetailDAO
	dao_dataset repo_dataset.DatasetDAOInterface
	dao_device  repo_device.IDeviceDAO
}

var onceModeling sync.Once
var instanceModeling *ModelingService

func NewModelingService(dao repo.IModelingDAO,
	dao_details repo.IModelinDetailDAO,
	dao_dataset repo_dataset.DatasetDAOInterface,
	dao_device repo_device.IDeviceDAO,
) *ModelingService {
	onceModeling.Do(func() {
		logger.Debug("Modeling service instance")
		instanceModeling = &ModelingService{
			ctx:         context.Background(),
			dao:         dao,
			dao_details: dao_details,
			dao_dataset: dao_dataset,
			dao_device:  dao_device,
		}
	})

	return instanceModeling
}

func (svc *ModelingService) CreateEvaluation(req repo.EvaluationDTO) (*repo.ModelingDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("%+v", req))
	if modeling, err := svc.makeModelingDTO(req); err != nil {
		return nil, err
	} else {
		if inserted, err := svc.dao.InsertOne(svc.ctx, *modeling); err != nil {
			return nil, logger.CreateReport(&logger.CODE_DB_INSERT, err)
		} else {
			return repo.ConvertModelingEntToDTO(inserted), nil
		}
	}
}

func (svc *ModelingService) ReadByTask(task_id int) ([]*repo.ModelingDB, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"task_id": %d}`, task_id))
	if modelingList, err := svc.dao.SelectByTask(svc.ctx, task_id); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		for _, dto := range modelingList {
			sourceID := dto.ID
			if dto.ParentID > 0 {
				sourceID = dto.ParentID
			}
			if model, err := svc.dao.SelectBestModelByModeling(svc.ctx, sourceID); err == nil && model != "" {
				dto.Scores = svc.dao.SelectTestScoreByModelAndModeling(svc.ctx, model, dto.ID)
				dto.InfTime = svc.dao.SelectTestInfTimeByModelAndModeling(svc.ctx, model, dto.ID)
			}
		}
		return modelingList, nil
	}
}

func (svc *ModelingService) ReadModelingType(task_id int) ([]*repo.ModelingDB, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"task_id": %d}`, task_id))
	if modelingList, err := svc.dao.SelectModelingType(svc.ctx, task_id); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return modelingList, nil
	}
}

func (svc *ModelingService) ReadOne(id int) (*repo.ModelingDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	if modeling, err := svc.dao.SelectOne(svc.ctx, id); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertModelingEntToDTO(modeling), nil
	}
}

func (svc *ModelingService) ReadFull(id int) (*repo.ModelingDB, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	if modeling, err := svc.dao.SelectFull(svc.ctx, id); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return modeling, nil
	}
}

func (svc *ModelingService) StopModelingTask(id int) *logger.Report {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	scheduler := NewTaskScheduler()
	if modeling, err := svc.dao.SelectOne(svc.ctx, id); err != nil {
		return logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		scheduler.CancelTask(modeling)
	}

	return nil
}

func (svc *ModelingService) makeModelingDTO(req repo.EvaluationDTO) (*repo.ModelingDTO, *logger.Report) {
	if parent, err := svc.ReadOne(req.ParentID); err != nil {
		return nil, err
	} else {
		if params, err := svc.makeModelingParams(req, parent); err != nil {
			return nil, err
		} else {
			modeling := &repo.ModelingDTO{
				Params:       params,
				ModelingStep: utils.MODELING_STEP_IDLE,
				TaskID:       parent.TaskID,
				ParentID:     parent.ID,
				DatasetID:    req.DatasetID,
			}

			if req.IsPath {
				modeling.ModelingType = utils.MODELING_TYPE_BLIND
			} else {
				modeling.ModelingType = utils.MODELING_TYPE_EVALUATION
			}

			return modeling, nil
		}
	}
}

func (svc *ModelingService) makeModelingParams(req repo.EvaluationDTO, parent *repo.ModelingDTO) ([]string, *logger.Report) {
	params := req.UserParams

	parentParams := make(map[string]interface{})
	if err := json.Unmarshal([]byte(parent.Params[0]), &parentParams); err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	}

	req.TaskID = parent.TaskID
	if err := svc.makeCommonParams(req, params, parentParams); err != nil {
		return nil, err
	}

	var r *logger.Report
	switch params["engine_type"] {
	case utils.JOB_TYPE_VISION_CLS_SL:
		r = svc.makeParamsVCLSSL(params)
	case utils.JOB_TYPE_VISION_CLS_ML:
		r = svc.makeParamsVCLSML(params)
	case utils.JOB_TYPE_VISION_AD:
		r = svc.makeParamsVAD(params)
	case utils.JOB_TYPE_TABLE_CLS:
		r = svc.makeParamsTCLS(params)
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

func (svc *ModelingService) makeCommonParams(req repo.EvaluationDTO, params map[string]interface{}, parentParams map[string]interface{}) *logger.Report {
	cf := config_service.NewStatic()

	params["data_path"], _ = svc.dao_dataset.SelectDataPathByDataSetId(svc.ctx, req.DatasetID)
	params["save_path"] = cf.Get("ROOT_PATH") + "/task/" + req.EngineType + "/" + strconv.Itoa(req.TaskID) + "/"
	params["origin_path"] = parentParams["save_path"].(string) + "/" + strconv.Itoa(req.ParentID) + "/"
	params["origin_id"] = req.ParentID
	params["device_ids"] = params["gpus"]
	delete(params, "gpus")
	params["target_metric"] = parentParams["target_metric"]
	params["engine_type"] = parentParams["engine_type"]
	if req.IsPath {
		params["mode"] = "test"
	} else if !req.IsPath {
		params["mode"] = "train"
	}

	return nil
}

func (svc *ModelingService) makeParamsVCLSSL(params map[string]interface{}) *logger.Report {
	params["multi_label"] = false

	return nil
}

func (svc *ModelingService) makeParamsVCLSML(params map[string]interface{}) *logger.Report {
	params["multi_label"] = true
	delete(params, "image_resolution")

	return nil
}

func (svc *ModelingService) makeParamsVAD(params map[string]interface{}) *logger.Report {
	return nil
}

func (svc *ModelingService) makeParamsTCLS(params map[string]interface{}) *logger.Report {
	return nil
}

func (svc *ModelingService) makeParamsTREG(params map[string]interface{}) *logger.Report {
	return nil
}

func (svc *ModelingService) DeleteOne(id int) *logger.Report {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))

	modeling, err := svc.dao.SelectOne(svc.ctx, id)
	if err != nil {
		return logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	if modeling.ModelingStep == utils.MODELING_STEP_IDLE {
		if err := svc.dao.DeleteModelingAndTask(svc.ctx, id, modeling.TaskID); err != nil {
			return logger.CreateReport(&logger.CODE_DB_DELETE, err)
		} else {
			return nil
		}
	}

	if modeling.ModelingStep != utils.MODELING_STEP_COMPLETE && modeling.ModelingStep != utils.MODELING_STEP_FAIL && modeling.ModelingStep != utils.MODELING_STEP_CANCEL {
		return logger.CreateReport(&logger.CODE_MODELING_IN_PROGRESS, nil)
	}

	engineParams := repo.EngineParams{}
	if err := json.Unmarshal([]byte(modeling.Params[0]), &engineParams); err != nil {
		return logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	} else {
		if engineParams.SavePath != "" {
			devices, err := repo_device.New().SelectByGPU(svc.ctx, engineParams.DeviceIDs)
			if err != nil || len(devices) < 1 {
				return logger.CreateReport(&logger.CODE_MODELING_DEVICE_NOT_EXIST, err)
			}

			targetPath := engineParams.SavePath + "/" + strconv.Itoa(id) + "/"
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

	if err := svc.dao.DeleteModelingAndTask(svc.ctx, id, modeling.TaskID); err != nil {
		return logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}

	return nil
}

// ReadIDsByTaskIDs는 주어진 taskIDs에 해당하는 Task의 ID 목록을 조회하는 서비스 함수입니다.
func (svc *ModelingService) ReadIDsByTaskIDs(taskIDs []int) ([]int, *logger.Report) {
	// 디버깅을 위한 로그 출력
	logger.Debug(fmt.Sprintf(`{"task_ids": %v}`, taskIDs))

	// TaskDAO의 SelectIDsByTaskIDs 메서드를 호출하여 해당 Task ID 목록을 조회
	taskIDs, err := svc.dao.SelectIDsByTaskIDs(svc.ctx, taskIDs)
	if err != nil {
		// 오류가 발생한 경우 에러 리포트 반환
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	// 성공적으로 조회된 Task ID 목록 반환
	return taskIDs, nil
}

// ReadModelingModelsIDsByModelingIDs는 주어진 modelingIDs에 해당하는 ModelingModel의 ID 목록을 조회하는 서비스 함수입니다.
func (svc *ModelingService) ReadModelingModelsIDsByModelingIDs(modelingIDs []int) ([]int, *logger.Report) {
	// 디버깅을 위한 로그 출력
	logger.Debug(fmt.Sprintf(`{"modeling_ids": %v}`, modelingIDs))

	// ModelingDAO의 SelectModelingModelsIDsByModelingIDs 메서드를 호출하여 해당 ModelingModel ID 목록을 조회
	modelingModelIDs, err := svc.dao.SelectModelingModelsIDsByModelingIDs(svc.ctx, modelingIDs)
	if err != nil {
		// 오류가 발생한 경우 에러 리포트 반환
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	// 성공적으로 조회된 ModelingModel ID 목록 반환
	return modelingModelIDs, nil
}

// ReadModelingDetailsIDsByModelingIDs는 주어진 modelingIDs에 해당하는 ModelingDetail의 ID 목록을 조회하는 서비스 함수입니다.
func (svc *ModelingService) ReadModelingDetailsIDsByModelingIDs(modelingIDs []int) ([]int, *logger.Report) {
	// 디버깅을 위한 로그 출력
	logger.Debug(fmt.Sprintf(`{"modeling_ids": %v}`, modelingIDs))

	// ModelingDetailDAO의 SelectIDsByModelingIDs 메서드를 호출하여 해당 ModelingDetail ID 목록을 조회
	modelingDetailIDs, err := svc.dao_details.SelectIDsByModelingIDs(svc.ctx, modelingIDs)
	if err != nil {
		// 오류가 발생한 경우 에러 리포트 반환
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	// 성공적으로 조회된 ModelingDetail ID 목록 반환
	return modelingDetailIDs, nil
}

// ExportModelingTableToCSV는 주어진 modelingIDs에 해당하는 Modeling 데이터를 CSV로 내보내는 서비스 함수입니다.
func (svc *ModelingService) ExportModelingTableToCSV(modelingIDs []int) (string, *logger.Report) {
	// 디버깅을 위한 로그 출력
	logger.Debug(fmt.Sprintf(`{"modeling_ids": %v}`, modelingIDs))

	// ModelingDAO의 SelectMany 메서드를 호출하여 Modeling 데이터를 조회
	modelings, err := svc.dao.SelectMany(svc.ctx, modelingIDs)
	if err != nil {
		// 오류가 발생한 경우 에러 리포트 반환
		return "", logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	// ModelingCSVFormat 인스턴스를 생성하여 CSV 포맷을 정의
	modelingCSVFormat := csvformat.ModelingCSVFormat{}
	// CSV 헤더 정의
	header := modelingCSVFormat.GetHeader()
	// Modeling 데이터를 CSV 레코드로 변환
	var records [][]string
	for _, modeling := range modelings {
		// 각 Modeling을 CSV 레코드로 변환하여 추가
		records = append(records, modelingCSVFormat.ConvertToRecord(modeling))
	}

	// utils 패키지의 ExportToCSV 함수 호출하여 CSV 파일로 저장
	return utils.ExportToCSV("modeling.csv", header, records)
}

// ImportModelingTableFromCSV는 CSV 파일에서 Modeling 데이터를 가져와 데이터베이스에 삽입하는 서비스 함수입니다.
func (svc *ModelingService) ImportModelingTableFromCSV(filename string, taskIdMap map[int]int) (map[int]int, *logger.Report) {
	// ModelingCSVFormat 인스턴스를 생성하여 CSV 레코드 변환 및 파싱 작업 수행
	modelingCSVFormat := csvformat.ModelingCSVFormat{}
	// CSV 레코드에서 ModelingDTO로 변환하는 함수 정의
	importFunc := func(record []string) (*repo.ModelingDTO, int, error) {
		// CSV 레코드에서 ModelingDTO 객체로 변환
		recordModeling, _ := modelingCSVFormat.ParseRecord(record)
		// Task ID 변환
		recordModeling.TaskID = taskIdMap[recordModeling.TaskID]

		// 변환된 ModelingDTO 반환
		return recordModeling, recordModeling.ID, nil
	}

	// 새 Modeling 데이터를 데이터베이스에 삽입하는 함수 정의
	insertFunc := func(modeling *repo.ModelingDTO) (int, error) {
		// ModelingDAO의 InsertOne 메서드를 호출하여 Modeling 삽입
		newModeling, err := svc.dao.InsertOne(svc.ctx, *modeling)
		if err != nil {
			// 오류가 발생한 경우 에러 반환
			return -1, err
		}
		// 삽입된 Modeling의 ID 반환
		return newModeling.ID, nil
	}

	// utils 패키지의 ImportFromCSV 함수 호출하여 CSV 파일로부터 데이터를 가져와 삽입
	return utils.ImportFromCSV(filename, importFunc, insertFunc)
}

// ExportModelingDetailsTableToCSV는 주어진 modelingDetailIDs에 해당하는 ModelingDetail 데이터를 CSV로 내보내는 서비스 함수입니다.
func (svc *ModelingService) ExportModelingDetailsTableToCSV(modelingDetailIDs []int) (string, *logger.Report) {
	// 디버깅을 위한 로그 출력
	logger.Debug(fmt.Sprintf(`{"modeling_detail_ids": %v}`, modelingDetailIDs))

	// ModelingDetailDAO의 SelectMany 메서드를 호출하여 ModelingDetail 데이터를 조회
	modelingDetails, err := svc.dao_details.SelectMany(svc.ctx, modelingDetailIDs)
	if err != nil {
		// 오류가 발생한 경우 에러 리포트 반환
		return "", logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	// ModelingDetailCSVFormat 인스턴스를 생성하여 CSV 포맷을 정의
	detailsCSVFormat := csvformat.ModelingDetailCSVFormat{}
	// CSV 헤더 정의
	header := detailsCSVFormat.GetHeader()

	// ModelingDetail 데이터를 CSV 레코드로 변환
	var records [][]string
	for _, detail := range modelingDetails {
		// 각 ModelingDetail을 CSV 레코드로 변환하여 추가
		records = append(records, detailsCSVFormat.ConvertToRecord(detail))
	}

	// utils 패키지의 ExportToCSV 함수 호출하여 CSV 파일로 저장
	return utils.ExportToCSV("modeling_details.csv", header, records)
}

// ImportModelingDetailsTableFromCSV는 CSV 파일에서 ModelingDetail 데이터를 가져와 데이터베이스에 삽입하는 서비스 함수입니다.
func (svc *ModelingService) ImportModelingDetailsTableFromCSV(filename string, modelingIdMap map[int]int) (map[int]int, *logger.Report) {
	// ModelingDetailCSVFormat 인스턴스를 생성하여 CSV 레코드 변환 및 파싱 작업 수행
	detailsCSVFormat := csvformat.ModelingDetailCSVFormat{}
	// CSV 레코드에서 ModelingDetailDTO로 변환하는 함수 정의
	importFunc := func(record []string) (*repo.ModelingDetailDTO, int, error) {
		// CSV 레코드에서 ModelingDetailDTO 객체로 변환
		details, _ := detailsCSVFormat.ParseRecord(record)
		// Modeling ID 변환
		details.ModelingID = modelingIdMap[details.ModelingID]

		// 변환된 ModelingDetailDTO 반환
		return details, details.ID, nil
	}

	// 새 ModelingDetail 데이터를 데이터베이스에 삽입하는 함수 정의
	insertFunc := func(modelingDetail *repo.ModelingDetailDTO) (int, error) {
		// ModelingDetailDAO의 InsertOne 메서드를 호출하여 ModelingDetail 삽입
		newModelingDetail, err := svc.dao_details.InsertOne(svc.ctx, *modelingDetail)
		if err != nil {
			// 오류가 발생한 경우 에러 반환
			return -1, err
		}
		// 삽입된 ModelingDetail의 ID 반환
		return newModelingDetail.ID, nil
	}

	// utils 패키지의 ImportFromCSV 함수 호출하여 CSV 파일로부터 데이터를 가져와 삽입
	return utils.ImportFromCSV(filename, importFunc, insertFunc)
}

// ExportModelingModelsTableToCSV는 주어진 modelingModelIDs에 해당하는 ModelingModel 데이터를 CSV로 내보내는 서비스 함수입니다.
func (svc *ModelingService) ExportModelingModelsTableToCSV(modelingModelIDs []int) (string, *logger.Report) {
	// 디버깅을 위한 로그 출력
	logger.Debug(fmt.Sprintf(`{"modeling_model_ids": %v}`, modelingModelIDs))

	// ModelingModelDAO의 SelectManyModelingModels 메서드를 호출하여 ModelingModel 데이터를 조회
	modelingModels, err := svc.dao.SelectManyModelingModels(svc.ctx, modelingModelIDs)
	if err != nil {
		// 오류가 발생한 경우 에러 리포트 반환
		return "", logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	// ModelingModelCSVFormat 인스턴스를 생성하여 CSV 포맷을 정의
	modelsCSVFormat := csvformat.ModelingModelCSVFormat{}
	// CSV 헤더 정의
	header := modelsCSVFormat.GetHeader()

	// ModelingModel 데이터를 CSV 레코드로 변환
	var records [][]string
	for _, model := range modelingModels {
		// 각 ModelingModel을 CSV 레코드로 변환하여 추가
		records = append(records, modelsCSVFormat.ConvertToRecord(model))
	}

	// utils 패키지의 ExportToCSV 함수 호출하여 CSV 파일로 저장
	return utils.ExportToCSV("modeling_models.csv", header, records)
}

// ImportModelingModelsTableFromCSV는 CSV 파일에서 ModelingModel 데이터를 가져와 데이터베이스에 삽입하는 서비스 함수입니다.
func (svc *ModelingService) ImportModelingModelsTableFromCSV(filename string, modelingIdMap map[int]int) (map[int]int, *logger.Report) {
	// ModelingModelCSVFormat 인스턴스를 생성하여 CSV 레코드 변환 및 파싱 작업 수행
	modelsCSVFormat := csvformat.ModelingModelCSVFormat{}
	// CSV 레코드에서 ModelingModels로 변환하는 함수 정의
	importFunc := func(record []string) (*repo.ModelingModels, int, error) {
		// CSV 레코드에서 ModelingModels 객체로 변환
		models, _ := modelsCSVFormat.ParseRecord(record)
		// Modeling ID 변환
		models.ModelingID = modelingIdMap[models.ModelingID]

		// 변환된 ModelingModels 반환
		return models, models.ID, nil
	}

	// 새 ModelingModel 데이터를 데이터베이스에 삽입하는 함수 정의
	insertFunc := func(modelingModel *repo.ModelingModels) (int, error) {
		// ModelingModelDAO의 InsertModelingModels 메서드를 호출하여 ModelingModel 삽입
		newModelingModel, err := svc.dao.InsertModelingModels(svc.ctx, *modelingModel)
		if err != nil {
			// 오류가 발생한 경우 에러 반환
			return -1, err
		}
		// 삽입된 ModelingModel의 ID 반환
		return newModelingModel.ID, nil
	}

	// utils 패키지의 ImportFromCSV 함수 호출하여 CSV 파일로부터 데이터를 가져와 삽입
	return utils.ImportFromCSV(filename, importFunc, insertFunc)
}

// GetBestModelNameByModelingID 함수는 특정 모델링 ID에 대한 최적의 모델 이름을 조회합니다.
func (svc *ModelingService) GetBestModelNameByModelingID(modelingID int) (string, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d}`, modelingID))
	if model, err := svc.dao.SelectBestModelByModeling(svc.ctx, modelingID); err != nil {
		return "", logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return model, nil
	}
}

// GetBestModelNameListByModelingID 함수는 특정 모델링 ID에 대한 모든 최적 모델 이름 리스트를 반환합니다.
func (svc *ModelingService) GetBestModelNameListByModelingID(modelingID int) ([]string, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d}`, modelingID))
	if model_dict_str, err := svc.dao.SelectBestModelsByModelingId(svc.ctx, modelingID); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		var parsed map[string]map[string][]any
		err := json.Unmarshal([]byte(model_dict_str), &parsed)
		if err != nil {
			logger.Error(err)
		}

		var filenames []string
		for _, subMap := range parsed {
			for _, arr := range subMap {
				if len(arr) > 0 {
					if path, ok := arr[0].(string); ok {
						filename := filepath.Base(path)
						name := strings.TrimSuffix(filename, ".kaier") // 확장자 제거
						filenames = append(filenames, name)
					}
				}
			}
		}

		return filenames, nil
	}
}

// ParseModelingDetailDataByType 함수는 특정 모델링 ID와 모델 이름에 해당하는 상세 데이터를 타입별로 파싱합니다.
func (svc *ModelingService) ParseModelingDetailDataByType(modelingID int, modelName string) (map[string]any, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d, "model_name": "%s"}`, modelingID, modelName))

	records, err := svc.dao_details.SelectManyByModelingIDAndModelName(svc.ctx, modelingID, modelName)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	typedDataMap := make(map[string]any)
	for _, record := range records {
		if len(record.Data) == 0 {
			continue
		}
		var dataJson any
		if err := json.Unmarshal([]byte(record.Data[0]), &dataJson); err != nil {
			return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
		}

		typedDataMap[record.DataType] = dataJson
	}

	return typedDataMap, nil
}

// ReadModelingModelsByTypeAndModelingID 함수는 특정 모델링 ID와 데이터 타입에 해당하는 모델 리스트를 조회합니다.
func (svc *ModelingService) ReadModelingModelsByTypeAndModelingID(modelingID int, data_type string) (*repo.ModelingModels, *logger.Report) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d, "data_type": "%s"}`, modelingID, data_type))
	if modelingModel, err := svc.dao.SelectModelingModelsByTypeAndModelingID(svc.ctx, modelingID, data_type); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertModelingModelEntToDTO(modelingModel), nil
	}
}

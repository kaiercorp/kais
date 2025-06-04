package service

import (
	modules_dataset "api_server/dataset/modules"
	repo_dataset "api_server/dataset/repository"
	service_dataset "api_server/dataset/service"
	service_device "api_server/device/service"
	"api_server/logger"
	repo "api_server/tapi/repository"
	repo_task "api_server/task/repository"
	service_task "api_server/task/service"

	"api_server/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ITapiService interface {
	// StartModeling 함수는 새로운 모델링 작업을 시작합니다. req 를 가지고 task 를 생성
	//
	// 매개변수:
	//   - req: 모델링 시작 요청 정보를 담고 있는 StartModelingRequest 객체
	//   - engine_type: 사용할 엔진 타입 (예: Vision Classification 등)
	//
	// 반환값:
	//   - *repo.StartModelingResponse: 모델링 작업 시작 결과를 담은 응답 객체
	//   - *logger.Report: 오류 발생 시 해당 정보를 담은 Report 객체
	StartModeling(req repo.StartModelingRequest, engineType string) (*repo.StartModelingResponse, *logger.Report)

	// StopModeling 함수는 실행 중인 모델링 작업을 중단합니다.
	//
	// 매개변수:
	//   - modelingID: 중단할 모델링 작업의 ID
	//
	// 반환값:
	//   - *logger.Report: 오류 발생 시 해당 정보를 담은 Report 객체
	StopModeling(modelingID int) *logger.Report

	// ReadModelingList는 주어진 엔진 타입에 해당하는 모든 태스크에 대해 모델링 목록을 조회하고 성능 및 추론 시간 정보를 포함하여 반환합니다.
	//
	// 매개변수:
	//   - engineType: 조회할 엔진 타입 (예: vision 분류, 테이블 분류 등)
	//
	// 반환값:
	//   - []*repo.Modeling: 모델링 목록 (성능 및 추론 시간 포함)
	//   - *logger.Report: 오류 발생 시 리포트, 없으면 nil
	ReadModelingList(engineType string) ([]*repo.Modeling, *logger.Report)

	// 모델링 상세 정보 조회
	//
	// 매개변수:
	//   - modelingID: 조회할 모델링 ID
	//   - thresholdKey: 사용할 threshold 값 (예: "0.1")
	//
	// 반환값:
	//   - *repo.ModelingDetail: 모델링 상세 정보
	//   - *logger.Report: 실패 시 에러 리포트
	ReadModelingDetail(modelingID int, threshold string) (*repo.ModelingDetail, *logger.Report)
	LoadedModels() ([]repo.DeviceModelGroup, *logger.Report)
	LoadModel(reqDTO repo.TestDTO, engineType string) (*repo.TapiLoaddedModel, *logger.Report)
	UnloadModel(testId int, gpuId int) *logger.Report
	UnloadByGPUAndModelName(gpuID int, modelName string) *logger.Report
	GetTestIdByGPUAndModelName(gpuID int, modelName string) (int, error)
	InferenceVCLS(testId int, filename string, image multipart.File, heatmap string) (map[string]interface{}, *logger.Report)
	InferenceTabular(testId int, xFeatures []map[string]interface{}) (map[string]interface{}, *logger.Report)

	// GetSystemInfomation은 시스템에 등록된 활성 디바이스들과 해당 디바이스에 연결된 GPU 정보를 수집하여 반환합니다.
	//
	// 매개변수:
	//   - 없음
	//
	// 반환값:
	//   - *repo.SystemInformation: 시스템 디바이스 및 GPU 정보를 포함한 구조체
	//   - *logger.Report: 오류 발생 시 리포트, 없으면 nil
	GetSystemInfomation() (*repo.SystemInformation, *logger.Report)
	GenerateNewDatasetPath(filename string) ([]string, *logger.Report)
	GetNewDataset(unique_name string) (*repo_dataset.DatasetDTO, *logger.Report)
}

type TapiService struct {
	ctx             context.Context
	device_svc      service_device.IDeviceService
	gpu_svc         service_device.IGPUService
	task_svc        service_task.ITaskService
	modeling_svc    service_task.IModelingService
	dataset_svc     service_dataset.DatasetServiceInterface
	datasetroot_svc service_dataset.IDatasetRootService
	dataset_watcher modules_dataset.DatasetWatcherInterface

	//dao               repo_task.ITaskDAO
	dao_modeling      repo_task.IModelingDAO
	tapiLoaddedModels []repo.TapiLoaddedModel
	tapi_devices      []repo.Device
}

var once sync.Once
var instance *TapiService

func New(dv_svc service_device.IDeviceService, g_svc service_device.IGPUService,
	t_svc service_task.ITaskService, m_svc service_task.IModelingService,
	ds_svc service_dataset.DatasetServiceInterface, dsr_svc service_dataset.IDatasetRootService,
	ds_watcher modules_dataset.DatasetWatcherInterface, dao_modeling repo_task.IModelingDAO) *TapiService {
	once.Do(func() { // atomic, does not allow repeating
		logger.Debug("Tapi Service instance")
		instance = &TapiService{
			ctx:             context.Background(),
			device_svc:      dv_svc,
			gpu_svc:         g_svc,
			task_svc:        t_svc,
			modeling_svc:    m_svc,
			dataset_svc:     ds_svc,
			datasetroot_svc: dsr_svc,
			dataset_watcher: ds_watcher,
			dao_modeling:    dao_modeling,
		}
	})

	return instance
}

func (svc *TapiService) getActiveDevices() ([]*repo.Device, *logger.Report) {
	devices, r := svc.device_svc.ReadActive()
	if r != nil {
		return nil, r
	}
	tapi_devices := make([]*repo.Device, len(devices))

	for d_idx, device := range devices {
		tapi_devices[d_idx] = &repo.Device{
			ID:    device.ID,
			Name:  device.Name,
			IP:    device.IP,
			Port:  device.Port,
			IsUse: device.IsUse,
			Type:  device.Type,
		}

		gpus, r := svc.gpu_svc.ReadManyByDeviceID(device.ID)
		if r != nil {
			return nil, r
		}
		tapi_devices[d_idx].GPU = make([]*repo.DeviceGPU, len(gpus))
		for g_idx, gpu := range gpus {
			tapi_devices[d_idx].GPU[g_idx] = &repo.DeviceGPU{
				ID:       gpu.ID,
				Index:    gpu.Index,
				Name:     gpu.Name,
				UUID:     gpu.UUID,
				IsUse:    gpu.IsUse,
				State:    gpu.State,
				DeviceID: gpu.DeviceID,
			}
		}
	}

	return tapi_devices, nil
}

// GetSystemInfomation은 시스템에 등록된 활성 디바이스들과 해당 디바이스에 연결된 GPU 정보를 수집하여 반환합니다.
func (svc *TapiService) GetSystemInfomation() (*repo.SystemInformation, *logger.Report) {
	tapi_devices, err := svc.getActiveDevices()
	if err != nil {
		return nil, err
	}

	info := repo.SystemInformation{
		Devices: tapi_devices,
		VERSION: utils.SW_VERSION,
	}

	return &info, nil
}

// ReadModelingList는 주어진 엔진 타입에 해당하는 모든 태스크에 대해 모델링 목록을 조회하고 성능 및 추론 시간 정보를 포함하여 반환합니다.
func (svc *TapiService) ReadModelingList(engineType string) ([]*repo.Modeling, *logger.Report) {
	tasks, r := svc.task_svc.ReadByEngineType(engineType)
	if r != nil {
		return nil, r
	}

	var (
		wg            sync.WaitGroup
		mutex         sync.Mutex
		tapiModelings []*repo.Modeling
		errOnce       sync.Once
		firstErr      *logger.Report
	)

	for _, task := range tasks {
		task := task // 고루틴 캡처 주의

		wg.Add(1)
		go func() {
			defer wg.Done()

			modelings, r := svc.modeling_svc.ReadByTask(task.ID)
			if r != nil {
				errOnce.Do(func() { firstErr = r })
				return
			}

			if err := json.Unmarshal([]byte(task.Params[0]), &task.UserParams); err != nil {
				errOnce.Do(func() {
					firstErr = logger.CreateReport(&logger.CODE_JSON_MARSHAL, err)
				})
				return
			}
			datasetName := svc.getString(task.UserParams["dataset_name"])

			var localModelings []*repo.Modeling

			for _, m := range modelings {
				tapiM := &repo.Modeling{
					ID:           m.ID,
					State:        m.ModelingStep,
					Name:         task.Title,
					TargetMetric: task.TargetMetric,
					CreatedAt:    m.CreatedAt,
					UpdatedAt:    m.UpdatedAt,
					Dataset:      datasetName,
				}

				bestModelName, r := svc.modeling_svc.GetBestModelNameByModelingID(m.ID)
				if r != nil || bestModelName == "" {
					localModelings = append(localModelings, tapiM)
					continue
				}

				dataMap, r := svc.modeling_svc.ParseModelingDetailDataByType(m.ID, bestModelName)
				if r != nil || len(dataMap) == 0 {
					localModelings = append(localModelings, tapiM)
					continue
				}

				perf, infTime := svc.extractPerformanceAndInferenceTime(task.EngineType, dataMap)
				if perf != nil {
					tapiM.Performance = &perf
				}
				if infTime != nil {
					tapiM.InferenceTime = infTime
				}

				localModelings = append(localModelings, tapiM)
			}

			mutex.Lock()
			tapiModelings = append(tapiModelings, localModelings...)
			mutex.Unlock()
		}()
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return tapiModelings, nil
}

// extractPerformanceAndInferenceTime는 주어진 데이터 맵에서 엔진 타입에 따라 모델의 성능 지표와 추론 시간을 추출합니다.
//
// 매개변수:
//   - engineType: 엔진 타입 (예: Vision 분류, Table 분류 등)
//   - dataMap: 모델 상세 정보를 담은 맵
//
// 반환값:
//   - map[string]any: 성능 지표 (Performance)
//   - *float64: 평균 추론 시간 (Inference Time)
func (svc *TapiService) extractPerformanceAndInferenceTime(engineType string, dataMap map[string]any) (map[string]any, *float64) {
	var infTime *float64
	var perf map[string]any

	switch engineType {
	case utils.JOB_TYPE_VISION_CLS_ML:
		perf = make(map[string]any)
		for _, threshold := range []string{"0.1", "0.2", "0.3", "0.4", "0.5"} {
			perf[threshold] = make(map[string]float64)
		}

		metrics := []string{
			"image_accuracy", "image_precision", "image_recall", "image_f1_score",
			"label_accuracy", "label_precision", "label_recall", "label_f1_score",
		}
		for _, metric := range metrics {
			// 예: threshold_image_accuracy 형태의 키를 찾아 threshold 별 값 추출
			if metricData, ok := dataMap["threshold_"+metric].(map[string]any); ok {
				if thresholds, ok := metricData["threshold"].(map[string]any); ok {
					for k, v := range map[string]string{"0.1": "100", "0.2": "200",
						"0.3": "300", "0.4": "400", "0.5": "500"} {
						if val, ok := thresholds[v].(float64); ok {
							perf[k].(map[string]float64)[metric] = val
						}
					}
				}
			}
		}
		if infData, ok := dataMap["test_avg_inference_time"].(map[string]any); ok {
			if t, ok := infData["avg inference time"].(float64); ok {
				infTime = &t
			}
		}
	case utils.JOB_TYPE_VISION_CLS_SL:
		if testset_score, ok := dataMap["testset_score"].(map[string]any); ok {
			perf = testset_score
		}
		if infData, ok := dataMap["test_avg_inference_time"].(map[string]any); ok {
			if t, ok := infData["avg inference time"].(float64); ok {
				infTime = &t
			}
		}
	case utils.JOB_TYPE_TABLE_CLS, utils.JOB_TYPE_TABLE_REG:
		if testset_score, ok := dataMap["testset_score"].(map[string]any); ok {
			perf = testset_score
		}
		if t, ok := dataMap["test_inference_time"].(float64); ok {
			infTime = &t
		}

	}

	return perf, infTime
}

// 모델링 상세 정보 조회
func (svc *TapiService) ReadModelingDetail(modelingID int, thresholdKey string) (*repo.ModelingDetail, *logger.Report) {
	threshold, report := svc.parseThreshold(thresholdKey)
	if report != nil {
		return nil, report
	}

	modeling, report := svc.modeling_svc.ReadOne(modelingID)
	if report != nil {
		return nil, report
	}

	task, report := svc.task_svc.ReadOne(modeling.TaskID)
	if report != nil {
		return nil, report
	}

	if !svc.isSupportedEngine(task.EngineType) {
		return nil, logger.CreateReport(&logger.CODE_FAILE, nil)
	}

	if err := json.Unmarshal([]byte(task.Params[0]), &task.UserParams); err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	}

	imageSize, err := svc.parseImageSize(task.UserParams)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	}

	tapiModeling := &repo.ModelingDetail{
		Name:         task.Title,
		Dataset:      svc.getString(task.UserParams["dataset_name"]),
		TargetMetric: task.TargetMetric,
		ImageSize:    &imageSize,
		ID:           modeling.ID,
		State:        modeling.ModelingStep,
		CreatedAt:    modeling.CreatedAt,
		UpdatedAt:    modeling.UpdatedAt,
		Models:       []*repo.ModelingModel{},
	}

	modelNameList, report := svc.modeling_svc.GetBestModelNameListByModelingID(modeling.ID)
	if report != nil || modelNameList == nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, nil)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, modelName := range modelNameList {
		wg.Add(1)
		go func(modelName string) {
			defer wg.Done()
			model := svc.buildModelDetail(modeling.ID, modelName, threshold, task, modeling.UpdatedAt)
			if model != nil {
				mu.Lock()
				tapiModeling.Models = append(tapiModeling.Models, model)
				mu.Unlock()
			}
		}(modelName)
	}
	wg.Wait()

	return tapiModeling, nil
}

// thresholdKey 문자열을 내부 사용 포맷으로 변환 (예: "0.1" → "100")
func (svc *TapiService) parseThreshold(thresholdKey string) (string, *logger.Report) {
	val, err := strconv.ParseFloat(thresholdKey, 64)
	if err != nil {
		return "", logger.CreateReport(&logger.CODE_FAILE, err)
	}
	return strconv.Itoa(int(val * 1000)), nil
}

// 지원하는 엔진 타입인지 여부 확인
func (svc *TapiService) isSupportedEngine(engineType string) bool {
	return engineType == utils.JOB_TYPE_VISION_CLS_ML || engineType == utils.JOB_TYPE_VISION_AD
}

// image_resolution 파라미터를 정수로 파싱
func (svc *TapiService) parseImageSize(params map[string]any) (int, error) {
	if str, ok := params["image_resolution"].(string); ok {
		return strconv.Atoi(str)
	}
	return 0, fmt.Errorf("invalid image_resolution format")
}

// any 타입에서 문자열로 변환 (실패 시 빈 문자열 반환)
func (svc *TapiService) getString(v any) string {
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

// 모델 상세 정보 구성
//
// 매개변수:
//   - modelingID: 모델링 ID
//   - modelName: 모델 이름
//   - threshold: threshold 값 (예: "100")
//   - task: 관련 태스크 정보
//   - updatedAt: 업데이트 시간
//
// 반환값:
//   - *repo.ModelingModel: 모델 성능 결과 정보
func (svc *TapiService) buildModelDetail(modelingID int, modelName, threshold string, task *repo_task.TaskDTO, updatedAt time.Time) *repo.ModelingModel {
	dataMap, r := svc.modeling_svc.ParseModelingDetailDataByType(modelingID, modelName)
	if r != nil || dataMap == nil {
		return nil
	}

	// 모델 이름에서 실제 이름만 추출 (형식: [모델타입]모델이름)
	re := regexp.MustCompile(`\[(.*?)\](.+)`)
	match := re.FindStringSubmatch(modelName)
	if len(match) < 3 {
		return nil
	}

	model := &repo.ModelingModel{
		ModelName: match[2],
		UpdatedAt: updatedAt,
		Result:    make(map[string]any),
	}

	// 분류 모델인 경우 각 성능 지표를 추출하여 결과에 추가
	if task.EngineType == utils.JOB_TYPE_VISION_CLS_ML {
		for _, metric := range []string{
			"image_accuracy", "image_precision", "image_recall", "image_f1_score",
			"label_accuracy", "label_precision", "label_recall", "label_f1_score",
		} {
			val := svc.extractMetric(dataMap, "threshold_"+metric, threshold)
			if val != nil {
				model.Result[metric] = val
			}
		}

		if infData, ok := dataMap["test_avg_inference_time"].(map[string]any); ok {
			if t, ok := infData["avg inference time"].(float64); ok {
				model.InferenceTime = t
			}
		}
	}

	if score, ok := model.Result[task.TargetMetric].(float64); ok {
		model.Score = score
	}

	return model
}

// metric 추출 (threshold 기반 값 파싱)
//
// 매개변수:
//   - dataMap: 전체 결과 데이터
//   - key: 찾을 metric key (예: "threshold_image_accuracy")
//   - threshold: threshold 값 (예: "100")
//
// 반환값:
//   - *float64: 해당 성능 지표 값 (없으면 nil)
func (svc *TapiService) extractMetric(dataMap map[string]any, key, threshold string) *float64 {
	if metricData, ok := dataMap[key].(map[string]any); ok {
		if thresholdMap, ok := metricData["threshold"].(map[string]any); ok {
			thresholdVal := thresholdMap[threshold].(float64)
			return &thresholdVal
		}
	}
	return nil
}

// StartModeling 함수는 새로운 모델링 작업을 시작합니다. req 를 가지고 task 를 생성
func (svc *TapiService) StartModeling(req repo.StartModelingRequest, engine_type string) (*repo.StartModelingResponse, *logger.Report) {
	dataset, r := svc.dataset_svc.ReadDataset(req.DatasetID)
	if r != nil {
		return nil, r
	}
	userParams := make(map[string]any)
	userParams["image_resolution"] = req.ImageSize
	userParams["gpu_auto"] = true
	userParams["gpus"] = req.GpuID
	userParams["index_column"] = req.IndexColumn
	userParams["output_columns"] = req.OutputColumns
	userParams["input_columns"] = req.InputColumns
	userParams["task_mode"] = req.TaskMode
	userParams["dataset_name"] = dataset.Name

	// 임의의 projectID 사용
	default_projectID := 1
	taskReq := repo_task.TaskDTO{
		ProjectID:    &default_projectID,
		DatasetID:    &req.DatasetID,
		EngineType:   engine_type,
		TargetMetric: req.TargetMetric,
		UserParams:   userParams,
		Title:        "",
		Description:  "",
	}
	// taskid 로 모델링 id 검색
	t, m, r := svc.task_svc.Create(taskReq)
	if r != nil {
		return nil, r
	}

	return &repo.StartModelingResponse{
		TrialID:      m.ID,
		State:        m.ModelingStep,
		TargetMetric: t.TargetMetric,
		GpuAuto:      true,
		DatasetName:  dataset.Name,
	}, nil
}

// StopModeling 함수는 실행 중인 모델링 작업을 중단합니다.
func (svc *TapiService) StopModeling(modelingID int) *logger.Report {

	return svc.modeling_svc.StopModelingTask(modelingID)
}

func (svc *TapiService) LoadedModelsGroupedByDevice() ([]repo.DeviceModelGroup, *logger.Report) {
	// UUID → DeviceID, DeviceName 매핑
	uuidToDevice := make(map[string]struct {
		ID   int
		Name string
	})
	for _, dev := range svc.tapi_devices {
		for _, gpu := range dev.GPU {
			uuidToDevice[strings.TrimSpace(gpu.UUID)] = struct {
				ID   int
				Name string
			}{
				ID:   dev.ID,
				Name: dev.Name,
			}
		}
	}

	// DeviceID → Models 그룹핑
	deviceToModels := make(map[int][]repo.TapiLoaddedModel)
	deviceToName := make(map[int]string)

	for _, model := range svc.tapiLoaddedModels {
		cleanUUID := strings.TrimSpace(model.GPUUUID)
		if devInfo, ok := uuidToDevice[cleanUUID]; ok {
			deviceToModels[devInfo.ID] = append(deviceToModels[devInfo.ID], model)
			deviceToName[devInfo.ID] = devInfo.Name
		}
	}

	// 결과 조립
	var result []repo.DeviceModelGroup
	for devID, models := range deviceToModels {
		result = append(result, repo.DeviceModelGroup{
			DeviceID:   devID,
			DeviceName: deviceToName[devID],
			Models:     models,
		})
	}

	return result, nil
}

func (svc *TapiService) LoadedModels() ([]repo.DeviceModelGroup, *logger.Report) {
	logger.Debug("Get loadedModels ")

	devices, err1 := instance.getActiveDevices()
	if err1 != nil {
		return nil, err1
	}

	svc.tapi_devices = make([]repo.Device, 0, len(devices))
	for _, d := range devices {
		if d != nil {
			svc.tapi_devices = append(svc.tapi_devices, *d)
		}
	}

	result, err2 := instance.LoadedModelsGroupedByDevice()
	if err2 != nil {
		return nil, err2
	}
	return result, nil
}

func (svc *TapiService) LoadModel(reqDTO repo.TestDTO, engineType string) (*repo.TapiLoaddedModel, *logger.Report) {
	logger.Debug("Loading model: "+reqDTO.ModelName, reqDTO.ModelingID, engineType, reqDTO.GpuId)

	bestModelStr, err := svc.dao_modeling.SelectBestModelsByModelingId(svc.ctx, reqDTO.ModelingID)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	var modelMap map[string]map[string][]interface{}

	// JSON 문자열을 맵으로 파싱
	errModelMap := json.Unmarshal([]byte(bestModelStr), &modelMap)
	if errModelMap != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, errModelMap)
	}

	modelPath := ""
	for _, metricData := range modelMap {
		for _, modelInfo := range metricData {
			if strings.Contains(modelInfo[0].(string), reqDTO.ModelName) {
				modelPath = modelInfo[0].(string)
				break
			}
		}
	}

	// API 서버 요청 준비
	url := "http://localhost:5000/api/load"

	stringGpuId := strconv.Itoa(reqDTO.GpuId)

	// 요청 파라미터 구성
	reqData := map[string]interface{}{
		"device_id":  stringGpuId,
		"model_name": reqDTO.ModelName,
		"model_path": modelPath,
		"model_type": engineType,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_MARSHAL, err)
	}

	// HTTP 요청 생성
	req, err := http.NewRequestWithContext(svc.ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_CREATE_REQ, err)
	}

	req.Header.Set("Content-Type", "application/json")

	// HTTP 요청 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_REQUEST, err)
	}
	defer resp.Body.Close()

	// 응답 처리
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_RESPONSE, err)
	}

	// 응답 데이터 파싱
	var apiResponse struct {
		Name   string `json:"name"`
		Models []struct {
			ModelNum  int    `json:"model_num"`
			ModelFile string `json:"model_file"`
			Engine    string `json:"engine"`
		} `json:"models"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	}

	loadded := repo.TapiLoaddedModel{
		TestId:    len(svc.tapiLoaddedModels),
		ModelName: reqDTO.ModelName,
		ModelNum:  apiResponse.Models[0].ModelNum,
		GPUIndex:  reqDTO.GpuId,
	}
	// gpuService := svc.gpu_svc.NewStatic()
	gpuService := svc.gpu_svc
	gpu, r := gpuService.ViewGpuByIndex(stringGpuId)
	if r != nil {
		logger.Error(fmt.Errorf("failed to get GPU info [gpu_index: %d]: %w", reqDTO.GpuId, err))
	} else {
		loadded.GPUID = gpu.ID
		loadded.GPUUUID = gpu.UUID
		svc.tapiLoaddedModels = append(svc.tapiLoaddedModels, loadded)
	}

	// 성공 응답 반환
	return &loadded, nil
}

func (s *TapiService) UnloadModel(testId int, gpuId int) *logger.Report {
	logger.Debug("TAPI Unload model: ", testId)

	url := "http://localhost:5000/api/model"

	for _, loadded := range s.tapiLoaddedModels {
		if loadded.TestId == testId {

			reqData := map[string]interface{}{
				"device_id":  strconv.Itoa(gpuId),
				"model_name": loadded.ModelName,
				"model_num":  loadded.ModelNum,
			}

			jsonData, err := json.Marshal(reqData)
			if err != nil {
				return logger.CreateReport(&logger.CODE_JSON_MARSHAL, err)
			}

			// HTTP 요청 생성
			req, err := http.NewRequestWithContext(s.ctx, "DELETE", url, bytes.NewBuffer(jsonData))
			if err != nil {
				return logger.CreateReport(&logger.CODE_REMOTE_CREATE_REQ, err)
			}

			req.Header.Set("Content-Type", "application/json")

			// HTTP 요청 실행
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return logger.CreateReport(&logger.CODE_REMOTE_REQUEST, err)
			}
			defer resp.Body.Close()

			// 응답 처리
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return logger.CreateReport(&logger.CODE_REMOTE_RESPONSE, err)
			}

			bodyStr := string(body)
			logger.Debug(bodyStr)

			// "Successfully" 포함 시 성공 처리
			if strings.Contains(bodyStr, "Successfully") {
				targetModelName := loadded.ModelName
				filtered := make([]repo.TapiLoaddedModel, 0)
				for _, m := range s.tapiLoaddedModels {
					if m.ModelName != targetModelName {
						filtered = append(filtered, m)
					}
				}
				s.tapiLoaddedModels = filtered

				return nil
			}

			return logger.CreateReport(&logger.CODE_REMOTE_RESPONSE, fmt.Errorf("unload failed: %s", bodyStr))
		}
	}

	// 모델을 찾지 못한 경우
	return logger.CreateReport(&logger.CODE_REMOTE_NOT_FOUND_MODEL, fmt.Errorf("model with test ID %d not found", testId))
}

func (s *TapiService) UnloadByGPUAndModelName(gpuID int, modelName string) *logger.Report {
	var report *logger.Report
	found := false

	for _, loadded := range s.tapiLoaddedModels {
		if loadded.GPUID == gpuID && loadded.ModelName == modelName {
			found = true
			report = s.UnloadModel(loadded.TestId, loadded.GPUID)
		}
	}

	if !found {
		return logger.CreateReport(&logger.CODE_TESTID_NOT_EXIST, errors.New("no test_id found for given GPU ID and Model Name"))
	}
	return report
}

func (s *TapiService) GetTestIdByGPUAndModelName(gpuID int, modelName string) (int, error) {
	logger.Debug(fmt.Errorf("get test id by gpu_id[%d] and model_name[%v]", gpuID, modelName))
	testId := -1
	found := false

	for _, loadded := range s.tapiLoaddedModels {
		if loadded.GPUID == gpuID && loadded.ModelName == modelName {
			found = true
			testId = loadded.TestId
			return testId, nil
		}
	}

	if !found {
		return -1, errors.New("test id not found for given GPU and model name")
	}

	return testId, nil
}

func (s *TapiService) InferenceVCLS(testId int, filename string, image multipart.File, heatmap string) (map[string]interface{}, *logger.Report) {
	logger.Debug("Inference VCLS for test ID: ", testId)
	// 로드된 모델 찾기
	var model *repo.TapiLoaddedModel
	for _, loadded := range s.tapiLoaddedModels {
		if loadded.TestId == testId {
			model = &loadded
			break
		}
	}

	if model == nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_NOT_FOUND_MODEL, fmt.Errorf("model with test ID %d not found", testId))
	}

	// multipart/form-data 준비
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// 이미지 파일 추가
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REQUEST, err)
	}

	// 파일 위치를 처음으로 되돌리기 (이미 읽혔을 수 있음)
	if seeker, ok := image.(io.Seeker); ok {
		_, err = seeker.Seek(0, io.SeekStart)
		if err != nil {
			return nil, logger.CreateReport(&logger.CODE_REQUEST, err)
		}
	}

	if _, err = io.Copy(fw, image); err != nil {
		return nil, logger.CreateReport(&logger.CODE_REQUEST, err)
	}

	// 추가 파라미터 설정
	params := map[string]string{
		"device_id":  "0",
		"model_name": model.ModelName,
		"model_num":  fmt.Sprintf("%d", model.ModelNum),
		"heatmap":    heatmap,
	}

	for key, value := range params {
		if err = w.WriteField(key, value); err != nil {
			return nil, logger.CreateReport(&logger.CODE_REQUEST, err)
		}
	}

	// multipart writer 닫기
	if err = w.Close(); err != nil {
		return nil, logger.CreateReport(&logger.CODE_REQUEST, err)
	}

	// HTTP 요청 준비
	url := "http://localhost:5000/api/vcls"
	req, err := http.NewRequestWithContext(s.ctx, "POST", url, &b)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_CREATE_REQ, err)
	}

	// Content-Type 헤더 설정
	req.Header.Set("Content-Type", w.FormDataContentType())

	// HTTP 요청 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_REQUEST, err)
	}
	defer resp.Body.Close()

	// 응답 본문 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_RESPONSE, err)
	}

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_RESPONSE,
			fmt.Errorf("server returned non-OK status: %d, body: %s", resp.StatusCode, string(body)))
	}

	// 응답 파싱
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	}

	var data map[string]interface{}
	data = utils.ConvertVCLSResult(result)

	// 성공 응답 반환
	return data, nil
}

func (s *TapiService) InferenceTabular(testId int, xFeatures []map[string]interface{}) (map[string]interface{}, *logger.Report) {
	logger.Debug("Inference TABULAR for test ID: ", testId)

	// 로드된 모델 찾기
	var model *repo.TapiLoaddedModel
	for _, loaded := range s.tapiLoaddedModels {
		if loaded.TestId == testId {
			model = &loaded
			break
		}
	}

	if model == nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_NOT_FOUND_MODEL, fmt.Errorf("model with test ID %d not found", testId))
	}

	// xFeatures를 JSON 문자열로 변환
	xInputJSON, err := json.Marshal(xFeatures)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_MARSHAL, err)
	}

	// multipart/form-data 준비
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// 파라미터 설정
	params := map[string]string{
		"device_id":  "0",
		"model_name": model.ModelName,
		"model_num":  fmt.Sprintf("%d", model.ModelNum),
		"x_input":    string(xInputJSON),
	}

	for key, value := range params {
		if err := w.WriteField(key, value); err != nil {
			return nil, logger.CreateReport(&logger.CODE_REQUEST, err)
		}
	}

	// multipart writer 닫기
	if err := w.Close(); err != nil {
		return nil, logger.CreateReport(&logger.CODE_REQUEST, err)
	}

	// HTTP 요청 준비
	url := "http://localhost:5000/api/tabular"
	req, err := http.NewRequestWithContext(s.ctx, "POST", url, &b)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_CREATE_REQ, err)
	}

	// Content-Type 헤더 설정
	req.Header.Set("Content-Type", w.FormDataContentType())

	// HTTP 요청 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_REQUEST, err)
	}
	defer resp.Body.Close()

	// 응답 본문 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_RESPONSE, err)
	}

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_RESPONSE,
			fmt.Errorf("server returned non-OK status: %d, body: %s", resp.StatusCode, string(body)))
	}

	// JSON 형식이 아닐 경우 대비
	if !json.Valid(body) {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_RESPONSE, fmt.Errorf("response is not valid JSON: %s", string(body)))
	}

	// 응답 파싱
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	}

	// 성공 응답 반환
	return result, nil
}

func (s *TapiService) GenerateNewDatasetPath(filename string) ([]string, *logger.Report) {
	datasetroot, r := s.datasetroot_svc.ViewDatasetrootActive()
	if r != nil {
		return nil, r
	}
	// UUID 생성
	uniqueID := uuid.New().String()
	ext := filepath.Ext(filename)                              // 확장자
	name := filename[:len(filename)-len(ext)] + "_" + uniqueID // 파일 이름 (확장자 제외)
	dir_path := datasetroot[0].Path + "/" + name + "/train"
	// 디렉토리가 없다면 생성
	if err := os.MkdirAll(dir_path, os.ModePerm); err != nil {
		return nil, logger.CreateReport(&logger.CODE_REQUEST,
			fmt.Errorf("unable to create directory"))
	}

	return []string{dir_path, name, ext}, nil
}

func (s *TapiService) GetNewDataset(unique_name string) (*repo_dataset.DatasetDTO, *logger.Report) {
	s.dataset_watcher.DetectDatasetModification()
	return s.dataset_svc.GetDatasetByName(unique_name)
}

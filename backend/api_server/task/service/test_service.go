package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"api_server/logger"
	repo "api_server/task/repository"
)

type ITestService interface {
	LoadModel(reqDTO repo.TestDTO) (*loaddedModel, *logger.Report)
	UnloadModel(testId int) *logger.Report
	InferenceVCLS(testId int, filename string, image multipart.File) (map[string]interface{}, *logger.Report)
	InferenceTabular(testId int, xFeatures []map[string]interface{}) (map[string]interface{}, *logger.Report)
	PredictTabular(reqDTO repo.TabularTestRequest) (map[string]interface{}, *logger.Report)
	FeatureImportanceLIME(reqDTO repo.TabularTestRequest) (map[string]interface{}, *logger.Report)
	GetDatasetColumns(modelingId int) (*repo.ColumnsResponse, *logger.Report)
	LoadedModels() ([]loaddedModel, *logger.Report)
}

type loaddedModel struct {
	TestId    int
	ModelName string
	ModelType string
	ModelNum  int
}
type tapiLoaddedModel struct {
	TestId    int    `json:"test_id"`
	ModelName string `json:"model_name"`
	ModelNum  int    `json:"model_num"`
	GPUIndex  string `json:"gpu_index"`
	GPUID     int    `json:"gpu_id"`
	GPUUUID   string `json:"gpu_uuid"`
}

// type tapiLoaddedModels struct {
// 	DeviceID   int                `json:"device_id"`
// 	DeviceName string             `json:"device_name"`
// 	Models     []tapiLoaddedModel `json:"models"`
// }

type TestService struct {
	ctx               context.Context
	dao               repo.ITaskDAO
	dao_modeling      repo.IModelingDAO
	loaddedModels     []loaddedModel
	tapiLoaddedModels []tapiLoaddedModel
}

// type LoaddedModels struct {
// 	DeviceID   int            `json:"device_id"`
// 	DeviceName string         `json:"device_name"`
// 	Models     []loaddedModel `json:"models"`
// }

var onceTest sync.Once
var instanceTest *TestService

func NewTestService(
	dao repo.ITaskDAO,
	dao_modeling repo.IModelingDAO,
) *TestService {
	onceTest.Do(func() {
		logger.Debug("Test Service instance")
		instanceTest = &TestService{
			ctx:          context.Background(),
			dao:          dao,
			dao_modeling: dao_modeling,
		}
	})

	return instanceTest
}

func (s *TestService) LoadModel(reqDTO repo.TestDTO) (*loaddedModel, *logger.Report) {
	logger.Debug("Loading model: " + reqDTO.ModelName)

	bestModelStr, err := s.dao_modeling.SelectBestModelsByModelingId(s.ctx, reqDTO.ModelingID)
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

	// 요청 파라미터 구성
	reqData := map[string]string{
		"device_id":  "0",
		"model_name": reqDTO.ModelName,
		"model_path": modelPath,
		"model_type": reqDTO.ModelType,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_MARSHAL, err)
	}

	// HTTP 요청 생성
	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewBuffer(jsonData))
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

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		logger.Debug("Non-JSON response body: ", string(body))
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, errors.New(string(body)))
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

	loadded := loaddedModel{
		TestId:    len(s.loaddedModels),
		ModelName: reqDTO.ModelName,
		ModelType: reqDTO.ModelType,
		ModelNum:  apiResponse.Models[0].ModelNum,
	}

	s.loaddedModels = append(s.loaddedModels, loadded)

	// 성공 응답 반환
	return &loadded, nil
}

func (s *TestService) UnloadModel(testId int) *logger.Report {
	logger.Debug("TAPI Unload model: ", testId, s.tapiLoaddedModels)

	url := "http://localhost:5000/api/model"

	for _, loadded := range s.loaddedModels {
		if loadded.TestId == testId {
			reqData := map[string]interface{}{
				"device_id":  "0",
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
				filtered := make([]loaddedModel, 0)
				for _, m := range s.loaddedModels {
					if m.ModelName != targetModelName {
						filtered = append(filtered, m)
					}
				}
				s.loaddedModels = filtered

				return nil
			}
			return logger.CreateReport(&logger.CODE_REMOTE_RESPONSE, fmt.Errorf("unload failed: %s", bodyStr))
		}
	}

	// 모델을 찾지 못한 경우
	return logger.CreateReport(&logger.CODE_REMOTE_NOT_FOUND_MODEL, fmt.Errorf("model with test ID %d not found", testId))
}

func (s *TestService) InferenceVCLS(testId int, filename string, image multipart.File) (map[string]interface{}, *logger.Report) {
	logger.Debug("Inference VCLS for test ID: ", testId)

	// type loaddedModel struct {
	// 	TestId    int
	// 	DeviceId  string
	// 	ModelName string
	// 	ModelType string
	// 	ModelNum  int
	// }

	// 로드된 모델 찾기
	var model *loaddedModel
	for _, loaded := range s.loaddedModels {
		if loaded.TestId == testId {
			model = &loaded
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
		"heatmap":    "true",
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

	if heatmapObj, exists := result["heatmap"].(map[string]interface{}); exists {
		// Iterate through all keys in the map
		for key, value := range heatmapObj {
			// Try to process each string value we find
			if base64Data, ok := value.(string); ok {
				// Check and process the base64 data
				if !strings.HasPrefix(base64Data, "data:image") && !strings.HasPrefix(base64Data, "http") {
					// If image data is raw base64 without a prefix, convert it to the proper format
					heatmapObj[key] = "data:image/png;base64," + base64Data
					logger.Debug("Added prefix to base64 image data for key: " + key)
				}
			}
		}
	}

	// 성공 응답 반환
	return result, nil
}

func (s *TestService) InferenceTabular(testId int, xFeatures []map[string]interface{}) (map[string]interface{}, *logger.Report) {
	logger.Debug("Inference TABULAR for test ID: ", testId)

	// 로드된 모델 찾기
	var model *loaddedModel
	for _, loaded := range s.loaddedModels {
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

func (s *TestService) PredictTabular(reqDTO repo.TabularTestRequest) (map[string]interface{}, *logger.Report) {
	logger.Debug("Predict TABULAR for modeling ID: ", reqDTO.ModelingID)

	// 입력 데이터 유효성 검사
	if len(reqDTO.YInputs) == 0 {
		return nil, logger.CreateReport(&logger.CODE_REQUEST, fmt.Errorf("empty features data"))
	}

	modeling, err := s.dao_modeling.SelectOne(s.ctx, reqDTO.ModelingID)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	engineParams := repo.EngineParams{}
	if err := json.Unmarshal([]byte(modeling.Params[0]), &engineParams); err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	}

	trial := strings.Split(reqDTO.ModelName, "_")[1]
	modelingPath := engineParams.SavePath + "/" + strconv.Itoa(reqDTO.ModelingID) + "/"

	// 요청 파라미터 구성
	params := map[string]interface{}{
		"origin_id":        reqDTO.ModelingID,
		"selected_trial":   trial,
		"y_input":          reqDTO.YInputs,
		"data_path":        engineParams.DataPath,
		"engine":           engineParams.EngineType,
		"pp_path":          modelingPath + "/preprocessing",
		"data_inform_path": modelingPath + "/data_inform.json",
	}

	// JSON 요청 본문 생성
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_MARSHAL, err)
	}

	// HTTP 요청 준비
	url := "http://localhost:5000/api/tabular/predict_x"
	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_CREATE_REQ, err)
	}

	// Content-Type 헤더 설정
	req.Header.Set("Content-Type", "application/json")

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

	// 이전 함수와 달리 이 엔드포인트는 응답 본문이 없을 수 있음
	// 응답 본문이 비어있지 않다면 파싱 시도
	var result map[string]interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
		}
		return result, nil
	}

	// 응답 본문이 비어있을 경우 (HTTP 200만 반환)
	return map[string]interface{}{"status": "success"}, nil
}

func (s *TestService) GetDatasetColumns(modelingId int) (*repo.ColumnsResponse, *logger.Report) {
	dataInformPath := ""

	modeling, err := s.dao_modeling.SelectOne(s.ctx, modelingId)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_API_PARAM_ENGINE, err)
	}

	engineParams := repo.EngineParams{}
	if err := json.Unmarshal([]byte(modeling.Params[0]), &engineParams); err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	} else {
		if engineParams.SavePath != "" {
			dataInformPath = engineParams.SavePath + "/" + strconv.Itoa(modelingId) + "/"
		}
	}

	// multipart/form-data 생성
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// form 필드 추가
	err = writer.WriteField("data_inform_path", dataInformPath)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_API_PARAM_ENGINE, err)
	}

	// form 닫기
	err = writer.Close()
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_API_PARAM_ENGINE, err)
	}

	// Flask 서버로 요청 보내기
	req, err := http.NewRequest("POST", "http://localhost:5000/api/tabular/columns", body)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_CREATE_REQ, err)
	}

	// 헤더 설정
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 요청 보내기
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_REQUEST, err)
	}
	defer resp.Body.Close()

	// 응답 읽기
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_RESPONSE, err)
	}

	// 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_RESPONSE, err)
	}

	// 성공 응답 파싱
	var columnsResp repo.ColumnsResponse
	err = json.Unmarshal(respBody, &columnsResp)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	}

	return &columnsResp, nil
}

func (s *TestService) FeatureImportanceLIME(reqDTO repo.TabularTestRequest) (map[string]interface{}, *logger.Report) {
	logger.Debug("Feature importance LIME for modeling ID: ", reqDTO.ModelingID)

	modeling, err := s.dao_modeling.SelectOne(s.ctx, reqDTO.ModelingID)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	engineParams := repo.EngineParams{}
	if err := json.Unmarshal([]byte(modeling.Params[0]), &engineParams); err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	}

	trial := strings.Split(reqDTO.ModelName, "_")[1]
	modelingPath := engineParams.SavePath + "/" + strconv.Itoa(reqDTO.ModelingID) + "/"

	// 요청 파라미터 구성
	params := map[string]interface{}{
		"origin_id":        reqDTO.ModelingID,
		"x_input":          reqDTO.XInputs,
		"selected_trial":   trial,
		"data_path":        engineParams.DataPath,
		"engine":           engineParams.EngineType,
		"pp_path":          modelingPath + "/preprocessing",
		"data_inform_path": modelingPath + "/data_inform.json",
	}

	// JSON 요청 본문 생성
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_MARSHAL, err)
	}

	// HTTP 요청 준비
	url := "http://localhost:5000/api/tabular/lime"
	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_REMOTE_CREATE_REQ, err)
	}

	// Content-Type 헤더 설정
	req.Header.Set("Content-Type", "application/json")

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

	// 이전 함수와 달리 이 엔드포인트는 응답 본문이 없을 수 있음
	// 응답 본문이 비어있지 않다면 파싱 시도
	var result map[string]interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
		}
		return result, nil
	}

	// 응답 본문이 비어있을 경우 (HTTP 200만 반환)
	return map[string]interface{}{"status": "success"}, nil
}

func (s *TestService) LoadedModels() ([]loaddedModel, *logger.Report) {
	logger.Debug("Get loadedModels ", s.loaddedModels)

	return s.loaddedModels, nil
}

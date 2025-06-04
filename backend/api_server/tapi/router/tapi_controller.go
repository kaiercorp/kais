package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"api_server/logger"
	repo "api_server/tapi/repository"
	"api_server/tapi/service"
	"api_server/utils"
)

type TapiController struct {
	svc service.ITapiService
}

var once sync.Once
var instance *TapiController

func New(svc service.ITapiService) *TapiController {
	once.Do(func() { // atomic, does not allow repeating
		logger.Debug("Tapi Controller instance")
		instance = &TapiController{
			svc: svc,
		}
	})

	return instance
}

// @Summary Start modeling
// @Description Create and start new modling task
// @Tags Modeling
// @Success 200
// @Router /{engine}/modeling [post]
// @Accept json
// @Param engine path string true "vcls-ml || vcls-sl"
// @Param Body body dto.ApiModelRequest true "modeling parameters"
func (ctlr *TapiController) StartModeling(c *gin.Context) {
	logger.ApiRequest(c)

	engineType := strings.ReplaceAll(c.Param("engine"), "_", "-")

	trialRequest := repo.StartModelingRequest{}
	if errParam := c.ShouldBindJSON(&trialRequest); errParam != nil {
		r := logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, errParam)
		logger.ApiResponse(c, r, nil)
		return
	}

	data, r := ctlr.svc.StartModeling(trialRequest, engineType)
	logger.ApiResponse(c, r, data)
}

// @Summary Stop modeling
// @Description Stop modeling task by modeling id
// @Tags Modeling
// @Success 200
// @Router /{engine}/modeling/{modeling_id} [delete]
// @Param modeling_id path int true "modeling_id"
func (ctlr *TapiController) StopModeling(c *gin.Context) {
	logger.ApiRequest(c)

	modeling_id, errParam := strconv.Atoi(c.Param("modeling_id"))
	if errParam != nil {
		r := logger.CreateReport(&logger.CODE_API_PARAM_ENGINE, errParam)
		logger.ApiResponse(c, r, nil)
		return
	}

	r := ctlr.svc.StopModeling(modeling_id)

	logger.ApiResponse(c, r, gin.H{"trial_id": modeling_id})
}

// @Summary Get modeling detail
// @Description Get modeling detail information by modeling id
// @Tags Modeling
// @Success 200
// @Router /{engine}/modeling/{modeling_id}/{threshold} [get]
// @Param modeling_id path int true "modeling_id"
// @Param threshold path string true "0.1 or 0.2 or 0.3 or 0.4 or 0.5"
func (ctlr *TapiController) ListModelingModel(c *gin.Context) {
	logger.ApiRequest(c)

	modeling_id, err := strconv.Atoi(c.Param("modeling_id"))
	if err != nil {
		r := logger.CreateReport(&logger.CODE_API_PARAM_ENGINE, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	threshold := c.Param("threshold")
	modeling, r := ctlr.svc.ReadModelingDetail(modeling_id, threshold)
	logger.ApiResponse(c, r, modeling)
}

// @Summary Get Modeling list
// @Description Get modeling list
// @Tags Modeling
// @Success 200
// @Router /{engine}/modeling/list [get]
// @Param engine path string true "vcls-ml || vcls-sl"
func (ctlr *TapiController) ListModeling(c *gin.Context) {
	logger.ApiRequest(c)

	engineType := strings.ReplaceAll(c.Param("engine"), "_", "-")
	data, r := ctlr.svc.ReadModelingList(engineType)
	logger.ApiResponse(c, r, data)
}

// @Summary Get recent modeling logs
// @Description Get recent modeling logs by modeling id
// @Tags Modeling
// @Success 200
// @Router /{engine}/logs/{modeling_id} [get]
// @Param engine path string true "vcls-ml || vcls-sl"
// @Param modeling_id path int true "modeling_id"
/*
func (ctlr *TapiController) ListModelingLogs(c *gin.Context) {
	logger.ApiRequest(c)

	modeling_id, errParam := strconv.Atoi(c.Param("modeling_id"))
	if errParam != nil {
		r := logger.CreateReport(&logger.ERROR_CODE_REQUEST_PARAM, errParam)
		logger.ApiResponse(c, r, nil)
		return
	}

	logs, err := engine_log_service.ViewLogsByModelingID(modeling_id)

	logger.ApiResponse(c, err, logs)
}
*/
// @Description Get System information
// @Tags system
// @Success 200
// @Router /sys [get]
func (ctlr *TapiController) GetSystemStatus(c *gin.Context) {
	logger.ApiRequest(c)

	data, r := ctlr.svc.GetSystemInfomation()

	logger.ApiResponse(c, r, data)
}

func (ctrl *TapiController) GetLoaddedModels(c *gin.Context) {
	logger.ApiRequest(c)

	result, err := ctrl.svc.LoadedModels()

	logger.TApiResponse(c, err, result)
}

func (ctrl *TapiController) PostLoadModel(c *gin.Context) {
	logger.ApiRequest(c)
	engineType := c.Param("engine")

	tapiReq := repo.TestDTO{}
	if err := c.ShouldBindJSON(&tapiReq); err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	} else {
		// gpu_id와 model_name, modelingid 넘겨주는 걸로 api 바꿀 예정
		reqDTO := repo.TestDTO{
			ModelingID: tapiReq.ModelingID,
			ModelName:  tapiReq.ModelName,
			GpuId:      tapiReq.GpuId,
		}

		_, err := ctrl.svc.LoadModel(reqDTO, engineType)
		logger.TApiResponse(c, err, nil)
	}
}

func (ctlr *TapiController) DeleteUnloadModel(c *gin.Context) {
	logger.Debug("Delete unload model")
	logger.ApiRequest(c)

	gpuID, errParam := strconv.Atoi(c.Param("gpu_id"))
	if errParam != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, errParam)
		logger.ApiResponse(c, r, nil)
		return
	}

	modelName := c.Param("model_name")
	if modelName != "" {
		report := ctlr.svc.UnloadByGPUAndModelName(gpuID, modelName)
		logger.TApiResponse(c, report, nil)
		return
	}
}

func (ctlr *TapiController) PostStartFileTest(c *gin.Context) {
	logger.ApiRequest(c)

	engineType := c.Param("engine")

	switch engineType {
	// Vision 모델 처리
	case utils.JOB_TYPE_VISION_CLS_ML, utils.JOB_TYPE_VISION_CLS_SL, utils.TAPI_JOB_TYPE_VISION_CLS_ML, utils.TAPI_JOB_TYPE_VISION_CLS_SL:
		var testRequest repo.InferenceRequest
		if err := c.ShouldBindWith(&testRequest, binding.FormMultipart); err != nil {
			logger.ApiResponse(c, logger.CreateReport(&logger.CODE_REQUEST, err), nil)
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			logger.ApiResponse(c, logger.CreateReport(&logger.CODE_REQUEST, err), nil)
			return
		}
		defer file.Close()

		if testRequest.ModelName == "" {
			logger.Debug("model_name is empty")
			return
		}

		if testRequest.Heatmap == "" {
			testRequest.Heatmap = "false"
		}

		testId, err := ctlr.svc.GetTestIdByGPUAndModelName(testRequest.GPUID, testRequest.ModelName)
		if err != nil || testId < 0 {
			logger.Debug("cannot get test id by gpu_id and model_name")
			return
		}

		data, report := ctlr.svc.InferenceVCLS(testId, header.Filename, file, testRequest.Heatmap)
		logger.TApiResponse(c, report, data)

	// Tabular 모델 처리
	case utils.JOB_TYPE_TABLE_CLS, utils.JOB_TYPE_TABLE_REG:
		var testRequest repo.InferenceTabularRequest
		if err := c.ShouldBindWith(&testRequest, binding.FormMultipart); err != nil {
			logger.ApiResponse(c, logger.CreateReport(&logger.CODE_REQUEST, err), nil)
			return
		}

		var xInputs []map[string]interface{}
		if err := json.Unmarshal([]byte(testRequest.XInputs), &xInputs); err != nil {
			logger.ApiResponse(c, logger.CreateReport(&logger.CODE_REQUEST, err), nil)
			return
		}

		if testRequest.ModelName == "" {
			logger.Debug("model_name is empty")
			return
		}

		testId, err := ctlr.svc.GetTestIdByGPUAndModelName(testRequest.GPUID, testRequest.ModelName)
		if err != nil || testId < 0 {
			logger.Debug("cannot get test id by gpu_id and model_name")
			return
		}

		data, report := ctlr.svc.InferenceTabular(testId, xInputs)
		logger.TApiResponse(c, report, data)

	default:
		logger.ApiResponse(c, logger.CreateReport(&logger.CODE_REQUEST, fmt.Errorf("unsupported engine type: %s", engineType)), nil)
	}
}

func (ctlr *TapiController) UploadDataset(c *gin.Context) {
	logger.ApiRequest(c)
	// 'file'이라는 이름의 form 데이터에서 파일을 가져옴
	file, err := c.FormFile("file")
	if err != nil {
		logger.ApiResponse(c, logger.CreateReport(&logger.CODE_REQUEST,
			fmt.Errorf("unable to get file from form")), nil)
		return
	}
	dirNameExt, r := ctlr.svc.GenerateNewDatasetPath(file.Filename)
	if r != nil {
		logger.ApiResponse(c, r, nil)
		return
	}
	filePath := dirNameExt[0] + "/" + dirNameExt[1] + dirNameExt[2]

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		logger.ApiResponse(c, logger.CreateReport(&logger.CODE_REQUEST,
			fmt.Errorf("unable to save file")), nil)
		return
	}

	ds, r := ctlr.svc.GetNewDataset(dirNameExt[1])
	if r != nil {
		logger.ApiResponse(c, r, nil)
		return
	}

	logger.ApiResponse(c, nil, gin.H{
		"dataset_name": file.Filename,
		"dataset_path": dirNameExt[0],
		"dataset_id":   ds.ID,
	})
}

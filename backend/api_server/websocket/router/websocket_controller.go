package router

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"api_server/logger"
	project_repo "api_server/project/repository"
	task_repo "api_server/task/repository"
	"api_server/utils"
	"api_server/websocket/dto"
)

func HandleWebSocketMessage(c *gin.Context) {

	ws, err := utils.WSUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}
	defer ws.Close()

	ctx := context.Background()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			logger.Error(err)
			return
		}

		webSocketDTO := dto.WebSocketDTO{}
		_ = json.Unmarshal([]byte(string(message)), &webSocketDTO)
		messageType := webSocketDTO.MessageType
		data := webSocketDTO.Data

		dao := task_repo.NewModelingDetailDAO()
		enginelog_dao := task_repo.NewEngineLogDAO()
		switch messageType {
		case "GET_PROJECT":
			getProject(ctx, ws) //, data)
		case "GET_TASK":
			getTask(ctx, ws, data)
		case "LOSS_CHART":
			logger.Debug(webSocketDTO.Data)
			getLossChart(ws, data, dao)
		case "MODEL_PERF":
			logger.Debug(webSocketDTO.Data)
			getModelPerf(ws, data, dao)
		case "CONFUSION_MATRIX":
			logger.Debug(webSocketDTO.Data)
			getConfusionMatrix(ws, data, dao)
		case "PRED_RESULT":
			logger.Debug(webSocketDTO.Data)
			getPredResult(ws, data, dao)
		case "THRESHOLD_RESULT":
			logger.Debug(webSocketDTO.Data)
			getThresholdResult(ws, data, dao)
		case "FEATURE_IMPORTANCE":
			logger.Debug(webSocketDTO.Data)
			getFeatureImportanceResult(ws, data, dao)
		case "HEATMAP_IMAGE":
			logger.Debug(webSocketDTO.Data)
			getHeatmapResult(ws, data, dao)
		case "ENGINE_LOG":
			logger.Debug(webSocketDTO.Data)
			getEngineLog(ws, data, enginelog_dao)
		}
	}
}

// getProject는 특정 사용자의 프로젝트 목록을 가져와서 WebSocket을 통해 클라이언트로 전송하는 함수입니다.
//
// *프론트엔드 ProjectListFetcher.tsx 참조*
func getProject(ctx context.Context, ws *websocket.Conn) { //, data map[string]interface{}) {
	/*
		// "username" 값을 data에서 추출하고, 올바른 타입인지 확인
		if username, ok := data["username"].(string); ok {
			// 사용자-프로젝트 매핑 데이터를 조회하는 DAO 인스턴스 생성
			userProjectDAO := project_repo.NewUserProject()
			// 특정 사용자에 해당하는 프로젝트 ID 목록을 조회
			userProjectIds, _ := userProjectDAO.SelectProjectIdsByUsername(ctx, username)

			// 조회된 프로젝트 ID 목록을 프로젝트 ID 배열로 변환
			projectIds := make([]int, len(userProjectIds))
			for _, up := range userProjectIds {
				projectIds = append(projectIds, up.ProjectID)
			}
	*/

	projectDAO := project_repo.New()
	// 프로젝트 ID 목록에 해당하는 프로젝트 정보를 조회
	//projects, _ := projectDAO.SelectMany(ctx, projectIds)
	projects, _ := projectDAO.SelectAll(ctx)

	message, _ := json.Marshal(dto.WebSocketDTO{
		MessageType: "GET_PROJECT",                                // 메시지 타입 설정
		Data:        map[string]interface{}{"projects": projects}, // 프로젝트 데이터 포함
	})

	_ = ws.WriteMessage(websocket.TextMessage, message)
	//}
}

func getTask(ctx context.Context, ws *websocket.Conn, data map[string]interface{}) {
	id := int(data["id"].(float64))
	dao_task := task_repo.NewTaskDAO()
	entity, _ := dao_task.SelectOne(ctx, id)

	message, _ := json.Marshal(dto.WebSocketDTO{
		MessageType: strings.Join([]string{"GET_MODELING_", strconv.Itoa(id)}, ""),
		Data:        map[string]interface{}{"modeling": entity},
	})

	_ = ws.WriteMessage(websocket.TextMessage, message)
}

func getLossChart(ws *websocket.Conn, data map[string]interface{}, dao *task_repo.ModelingDetailDAO) {
	if modeling_id, ok := data["modeling_id"].(float64); ok {
		engine_type, _ := dao.SelectModelingType(int(modeling_id))
		if engine_type == utils.JOB_TYPE_TABLE_CLS || engine_type == utils.JOB_TYPE_TABLE_REG {
			getTableLossChart(ws, int(modeling_id), dao)
		} else {
			getVisionLossChart(ws, int(modeling_id), dao)
		}
	}
}

func getVisionLossChart(ws *websocket.Conn, modeling_id int, dao *task_repo.ModelingDetailDAO) {
	if result, err := dao.SelectLossChart(modeling_id); err != nil {
		logger.Error(err)
	} else {
		message, _ := json.Marshal(dto.WebSocketDTO{
			MessageType: strings.Join([]string{
				"LOSS_CHART_",
				strconv.Itoa(int(modeling_id)),
			},
				"",
			),
			Data: map[string]interface{}{"chart_item": result, "modeling_step": result.ModelingStep},
		})

		_ = ws.WriteMessage(websocket.TextMessage, message)
	}
}

func getTableLossChart(ws *websocket.Conn, modeling_id int, dao *task_repo.ModelingDetailDAO) {
	if result, err := dao.SelectTabularChart(modeling_id); err != nil {
		logger.Error(err)
	} else {
		message, _ := json.Marshal(dto.WebSocketDTO{
			MessageType: strings.Join([]string{
				"LOSS_CHART_",
				strconv.Itoa(int(modeling_id)),
			},
				"",
			),
			Data: map[string]interface{}{"chart_item": result, "modeling_step": result.ModelingStep},
		})

		_ = ws.WriteMessage(websocket.TextMessage, message)
	}
}

func getModelPerf(ws *websocket.Conn, data map[string]interface{}, dao *task_repo.ModelingDetailDAO) {
	modeling_id, ok := data["modeling_id"].(float64)
	dataset_type, ok2 := data["dataset_type"].(string)

	if ok && ok2 {
		engine_type, _ := dao.SelectModelingType(int(modeling_id))
		if engine_type == utils.JOB_TYPE_TABLE_CLS || engine_type == utils.JOB_TYPE_TABLE_REG {
			getTabularModelPerf(ws, int(modeling_id), dataset_type, dao)
		} else {
			getVisionModelPerf(ws, int(modeling_id), dataset_type, dao)
		}
	}
}

func getVisionModelPerf(ws *websocket.Conn, modeling_id int, dataset_type string, dao *task_repo.ModelingDetailDAO) {
	if result, err := dao.SelectModelPerformance(modeling_id, dataset_type); err != nil {
		logger.Error(err)
	} else {
		message, _ := json.Marshal(dto.WebSocketDTO{
			MessageType: strings.Join([]string{
				"MODEL_PERF_",
				strconv.Itoa(modeling_id),
				"_",
				dataset_type,
			},
				"",
			),
			Data: map[string]interface{}{"rows": result.Rows},
		})

		_ = ws.WriteMessage(websocket.TextMessage, message)
	}
}

func getTabularModelPerf(ws *websocket.Conn, modeling_id int, dataset_type string, dao *task_repo.ModelingDetailDAO) {
	if result, err := dao.SelectTabularModelPerformance(modeling_id, dataset_type); err != nil {
		logger.Error(err)
	} else {
		message, _ := json.Marshal(dto.WebSocketDTO{
			MessageType: strings.Join([]string{
				"MODEL_PERF_",
				strconv.Itoa(modeling_id),
				"_",
				dataset_type,
			},
				"",
			),
			Data: map[string]interface{}{"rows": result.Rows},
		})

		_ = ws.WriteMessage(websocket.TextMessage, message)
	}
}

func getConfusionMatrix(ws *websocket.Conn, data map[string]interface{}, dao *task_repo.ModelingDetailDAO) {
	modeling_id, ok := data["modeling_id"].(float64)
	dataset_type, ok2 := data["dataset_type"].(string)
	model_name, ok3 := data["model"].(string)

	if ok && ok2 && ok3 {
		engine_type, _ := dao.SelectModelingType(int(modeling_id))
		switch engine_type {
		case utils.JOB_TYPE_TABLE_CLS:
			fallthrough
		case utils.JOB_TYPE_TABLE_REG:
			getTabularModelConfusionMatrix(ws, int(modeling_id), dataset_type, model_name, dao)
		case utils.JOB_TYPE_VISION_CLS_SL:
			getVisionSLModelConfusionMatrix(ws, int(modeling_id), dataset_type, model_name, dao)
		case utils.JOB_TYPE_VISION_CLS_ML:
			getVisionMLModelConfusionMatrix(ws, int(modeling_id), dataset_type, model_name, dao)
		default:
			logger.Error("Unknown Engine Type")
		}
	}
}

func getVisionSLModelConfusionMatrix(ws *websocket.Conn, modeling_id int, dataset_type string, model_name string, dao *task_repo.ModelingDetailDAO) {
	if result, err := dao.SelectVisionSLModelConfusionMatrix(modeling_id, dataset_type, model_name); err != nil {
		logger.Error(err)
	} else {
		message, _ := json.Marshal(dto.WebSocketDTO{
			MessageType: strings.Join([]string{
				"CONFUSION_MATRIX_",
				strconv.Itoa(int(modeling_id)),
				"_",
				dataset_type,
				"_",
				model_name,
			},
				"",
			),
			Data: map[string]interface{}{"rows": result.Rows, "summaries": result.Summaries, "sum": result.Sum, "acc": result.Acc},
		})

		_ = ws.WriteMessage(websocket.TextMessage, message)
	}
}

func getVisionMLModelConfusionMatrix(ws *websocket.Conn, modeling_id int, dataset_type string, model_name string, dao *task_repo.ModelingDetailDAO) {
	if result, err := dao.SelectVisionMLModelConfusionMatrix(modeling_id, dataset_type, model_name); err != nil {
		logger.Error(err)
	} else {
		message, _ := json.Marshal(dto.WebSocketDTO{
			MessageType: strings.Join([]string{
				"CONFUSION_MATRIX_",
				strconv.Itoa(int(modeling_id)),
				"_",
				dataset_type,
				"_",
				model_name,
			},
				"",
			),
			Data: map[string]interface{}{"rows": result.Rows, "summaries": result.Summaries},
		})

		_ = ws.WriteMessage(websocket.TextMessage, message)
	}
}

func getTabularModelConfusionMatrix(ws *websocket.Conn, modeling_id int, dataset_type string, model_name string, dao *task_repo.ModelingDetailDAO) {
	if result, err := dao.SelectVisionSLModelConfusionMatrix(modeling_id, dataset_type, model_name); err != nil {
		logger.Error(err)
	} else {
		message, _ := json.Marshal(dto.WebSocketDTO{
			MessageType: strings.Join([]string{
				"CONFUSION_MATRIX_",
				strconv.Itoa(int(modeling_id)),
				"_",
				dataset_type,
				"_",
				model_name,
			},
				"",
			),
			Data: map[string]interface{}{"rows": result.Rows, "summaries": result.Summaries, "sum": result.Sum, "acc": result.Acc},
		})
		_ = ws.WriteMessage(websocket.TextMessage, message)
	}
}

func getPredResult(ws *websocket.Conn, data map[string]interface{}, dao *task_repo.ModelingDetailDAO) {
	modeling_id, ok := data["modeling_id"].(float64)
	dataset_type, ok2 := data["dataset_type"].(string)
	model_name, ok3 := data["model"].(string)

	if ok && ok2 && ok3 {
		if result, err := dao.SelectModelSampleTest(int(modeling_id), dataset_type, model_name); err != nil {
			logger.Error(err)
		} else {
			message, _ := json.Marshal(dto.WebSocketDTO{
				MessageType: strings.Join([]string{
					"PRED_RESULT_",
					strconv.Itoa(int(modeling_id)),
					"_",
					dataset_type,
					"_",
					model_name,
				},
					"",
				),
				Data: map[string]interface{}{"rows": result.Rows},
			})

			_ = ws.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func getHeatmapResult(ws *websocket.Conn, data map[string]interface{}, dao *task_repo.ModelingDetailDAO) {
	task_id, ok := data["task_id"].(float64)
	engine_type, ok1 := data["engine_type"].(string)
	modeling_id, ok2 := data["modeling_id"].(float64)
	dataset_type, ok3 := data["dataset_type"].(string)
	model_name, ok4 := data["model"].(string)

	if ok && ok1 && ok2 && ok3 && ok4 {
		if result, err := dao.SelectHeatmapImage(int(task_id), engine_type, int(modeling_id), dataset_type, model_name); err != nil {
			logger.Error(err)
		} else {
			message, _ := json.Marshal(dto.WebSocketDTO{
				MessageType: strings.Join([]string{
					"HEATMAP_IMAGE_",
					strconv.Itoa(int(modeling_id)),
					"_",
					dataset_type,
					"_",
					model_name,
				},
					"",
				),
				Data: map[string]interface{}{"rows": result.KeyRows},
			})

			_ = ws.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func getEngineLog(ws *websocket.Conn, data map[string]interface{}, dao *task_repo.EngineLogDAO) {
	modeling_id, ok1 := data["modeling_id"].(float64)
	page, ok2 := data["page"].(float64)

	if ok1 && ok2 {
		if result, err := dao.SelectEngineLog(int(modeling_id)); err != nil {
			logger.Error(err)
		} else {
			message, _ := json.Marshal(dto.WebSocketDTO{
				MessageType: strings.Join([]string{
					"ENGINE_LOG_",
					strconv.Itoa(int(modeling_id)),
					"_",
					strconv.Itoa(int(page)),
				},
					"",
				),
				Data: map[string]interface{}{"rows": result},
			})

			_ = ws.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func getThresholdResult(ws *websocket.Conn, data map[string]interface{}, dao *task_repo.ModelingDetailDAO) {
	modeling_id, ok := data["modeling_id"].(float64)
	dataset_type, ok2 := data["dataset_type"].(string)
	model_name, ok3 := data["model"].(string)

	if ok && ok2 && ok3 {
		if result, err := dao.SelectThreshold(int(modeling_id), dataset_type, model_name); err != nil {
			logger.Error(err)
		} else {
			message, _ := json.Marshal(dto.WebSocketDTO{
				MessageType: strings.Join([]string{
					"THRESHOLD_RESULT_",
					strconv.Itoa(int(modeling_id)),
					"_",
					dataset_type,
					"_",
					model_name,
				},
					"",
				),
				Data: map[string]interface{}{"rows": result},
			})

			_ = ws.WriteMessage(websocket.TextMessage, message)
		}
	}

}

func getFeatureImportanceResult(ws *websocket.Conn, data map[string]interface{}, dao *task_repo.ModelingDetailDAO) {
	modeling_id, ok := data["modeling_id"].(float64)
	dataset_type, ok2 := data["dataset_type"].(string)
	model_name, ok3 := data["model"].(string)

	if ok && ok2 && ok3 {
		if result, err := dao.SelectFeatureImportanceChart(int(modeling_id), dataset_type, model_name); err != nil {
			logger.Error(err)
		} else {
			message, _ := json.Marshal(dto.WebSocketDTO{
				MessageType: strings.Join([]string{
					"FEATURE_IMPORTANCE_",
					strconv.Itoa(int(modeling_id)),
					"_",
					dataset_type,
					"_",
					model_name,
				},
					"",
				),
				Data: map[string]interface{}{"result": result},
			})

			_ = ws.WriteMessage(websocket.TextMessage, message)
		}
	}
}

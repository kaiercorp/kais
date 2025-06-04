package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	device_repo "api_server/device/repository"
	"api_server/ent"
	"api_server/logger"
	repo "api_server/task/repository"
	"api_server/utils"
)

type ITaskScheduler interface {
	WatchTasks()
	matchGPU(task *ent.Modeling) (string, io.Reader, []int, error)
	updateGPUSRunning(gpu_ids []int, state string)
	updateGPUIdle(gpu_id int)
	runComplete(task *ent.Modeling)
	runModeling(task *ent.Modeling)
	runEvaluation(task *ent.Modeling)
	completeTasks()
	runTasks()
	CancelTask(task *ent.Modeling)
}

type TaskScheduler struct {
	ctx          context.Context
	dao_task     repo.ITaskDAO
	dao_modeling repo.IModelingDAO
	dao_device   device_repo.DeviceDAO
	dao_gpu      device_repo.GPUDAO
}

var onceTaskScheduler sync.Once
var instanceTaskScheduler *TaskScheduler

func NewTaskScheduler() *TaskScheduler {
	onceTaskScheduler.Do(func() {
		logger.Debug("TaskScheduler Service instance")
		instanceTaskScheduler = &TaskScheduler{
			ctx:          context.Background(),
			dao_task:     repo.NewTaskDAO(),
			dao_modeling: repo.NewModelingDAO(),
			dao_device:   *device_repo.New(),
			dao_gpu:      *device_repo.NewGPUDAO(),
		}
	})

	return instanceTaskScheduler
}

func (scheduler *TaskScheduler) WatchTasks() {
	// KAI.S 재시작 하는 것에 엔진이 영향을 받나? 안 받을 것 같은데?
	// scheduler.dao_modeling.UpdateCancel(scheduler.ctx)

	for {
		time.Sleep(10 * time.Second)
		scheduler.completeTasks()
		scheduler.runTasks()
	}
}

func (scheduler *TaskScheduler) matchGPU(task *ent.Modeling) (string, io.Reader, []int, error) {
	// engineParams := make(map[string]interface{})
	engineParams := repo.EngineParams{}
	addr := ""
	if err := json.Unmarshal([]byte(task.Params[0]), &engineParams); err != nil {
		return "", nil, nil, err
	}

	engineParams.ModelingID = task.ID

	if engineParams.GPUAuto {
		// 자동 배치
		gpus, err := scheduler.dao_gpu.SelectIdle(scheduler.ctx)
		if err != nil {
			return "", nil, nil, err
		} else if len(gpus) < 1 {
			return "", nil, nil, errors.New("no idle gpus")
		}

		engineParams.DeviceIDs = []int{gpus[0].ID}

		if jsonBytes, err := json.Marshal(engineParams); err == nil {
			scheduler.dao_modeling.UpdateParams(scheduler.ctx, task.ID, []string{string(jsonBytes)})
		}
	}

	if devices, err := scheduler.dao_device.SelectIdleByGPU(scheduler.ctx, engineParams.DeviceIDs); err != nil {
		return "", nil, nil, err
	} else {
		if len(devices) > 1 {
			// TODO : multi node
			return "", nil, nil, errors.New("not implemented yet")
		} else if len(devices) > 0 {
			// gpu가 idle 상태면 run
			addr = "http://" + devices[0].IP + ":" + strconv.Itoa(devices[0].Port)
		} else {
			return "", nil, nil, errors.New("not enough gpus")
		}
	}

	logger.Debug(fmt.Sprintf("%+v", engineParams))

	return addr, bytes.NewReader([]byte(fmt.Sprintf(`{"modeling_id":%d}`, task.ID))), engineParams.DeviceIDs, nil
}

func (scheduler *TaskScheduler) updateGPUIdle(gpu_ids []int) {
	scheduler.dao_gpu.UpdateManyState(scheduler.ctx, gpu_ids, utils.GPU_STATE_IDLE)
}

func (scheduler *TaskScheduler) runComplete(task *ent.Modeling) {
	if addr, params, _, err := scheduler.matchGPU(task); err != nil {
		logger.Error(err)
	} else if resp, err := http.Post(
		addr+"/api/train/finish", "application/json",
		params); err != nil {
		logger.Error(err)
	} else if respBody, err := io.ReadAll(resp.Body); err != nil {
		logger.Error(err)
	} else {
		logger.Debug(respBody)
	}
}

func (scheduler *TaskScheduler) runModeling(task *ent.Modeling) {
	if addr, params, _, err := scheduler.matchGPU(task); err != nil {
		logger.Error(err)
	} else if resp, err := http.Post(
		addr+"/api/train", "application/json",
		params); err != nil {
		logger.Error(err)
	} else if respBody, err := io.ReadAll(resp.Body); err != nil {
		logger.Error(err)
	} else {
		logger.Debug(string(respBody))
		scheduler.dao_modeling.UpdateState(scheduler.ctx, task.ID, utils.MODELING_STEP_REQUEST)
	}
}

func (scheduler *TaskScheduler) runEvaluation(task *ent.Modeling) {
	if addr, params, _, err := scheduler.matchGPU(task); err != nil {
		logger.Error(err)
	} else if resp, err := http.Post(
		addr+"/api/evaluation",
		"application/json", params); err != nil {
		logger.Error(err)
	} else if respBody, err := io.ReadAll(resp.Body); err != nil {
		logger.Error(err)
	} else {
		logger.Debug(string(respBody))
		scheduler.dao_modeling.UpdateState(scheduler.ctx, task.ID, utils.MODELING_STEP_REQUEST)
	}
}

func (scheduler *TaskScheduler) runBlind(task *ent.Modeling) {
	if addr, params, _, err := scheduler.matchGPU(task); err != nil {
		logger.Error(err)
	} else if resp, err := http.Post(addr+"/api/vision/blind", "application/json", params); err != nil {
		logger.Error(err)
	} else {
		logger.Debug(resp)
	}
}

func (scheduler *TaskScheduler) completeTasks() {
	if taskEnts, err := scheduler.dao_modeling.SelectManyFinish(scheduler.ctx); err != nil {
		logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		for _, task := range taskEnts {
			scheduler.runComplete(task)
		}
	}
}

func (scheduler *TaskScheduler) CancelTask(task *ent.Modeling) {
	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(task.Params[0]), &params); err != nil {
		logger.Error(err)
		return
	}

	gpus := []float64{}
	gpu_ids := []int{}
	deviceBytes, _ := json.Marshal(params["device_ids"])
	if err := json.Unmarshal(deviceBytes, &gpus); err == nil {
		for _, v := range gpus {
			gpu_ids = append(gpu_ids, int(v))
		}
	}

	if devices, err := scheduler.dao_device.SelectByGPU(scheduler.ctx, gpu_ids); err != nil {
		return
	} else if len(devices) > 0 {
		client := &http.Client{}
		addr := "http://" + devices[0].IP + ":" + strconv.Itoa(devices[0].Port) + "/api/train/" + strconv.Itoa(task.ID)
		if req, err := http.NewRequest("DELETE", addr, nil); err != nil {
			logger.Error(err)
		} else if resp, err := client.Do(req); err != nil {
			logger.Error(err)
		} else if respBody, err := io.ReadAll(resp.Body); err != nil {
			logger.Error(err)
		} else {
			logger.Debug(respBody)
			scheduler.updateGPUIdle(gpu_ids)
		}
	}

	scheduler.dao_modeling.UpdateState(scheduler.ctx, task.ID, utils.MODELING_STEP_CANCEL)
}

func (scheduler *TaskScheduler) runTasks() {
	if taskEnts, err := scheduler.dao_modeling.SelectManyIdle(scheduler.ctx); err != nil {
		logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		for _, task := range taskEnts {
			switch task.ModelingType {
			case utils.MODELING_TYPE_INITIAL:
				fallthrough
			case utils.MODELING_TYPE_UPDATE:
				scheduler.runModeling(task)
			case utils.MODELING_TYPE_EVALUATION:
				scheduler.runEvaluation(task)
			case utils.MODELING_TYPE_BLIND:
				scheduler.runBlind(task)
			default:
				logger.Error("Invalid modeling task")
			}
		}
	}
}

package repository

import (
	"mime/multipart"
	"time"
)

type StartModelingRequest struct {
	DatasetID     int      `json:"dataset_id"`               // modeling.task_id	// task.id
	ImageSize     *int     `json:"image_size,omitempty"`     // modeling.task_type	// task.task_type
	IndexColumn   *string  `json:"index_column,omitempty"`   // modeling.task_state	// task.state
	OutputColumns []string `json:"output_columns,omitempty"` // modeling.task_state	// task.state
	InputColumns  []string `json:"input_columns,omitempty"`  // modeling.task_state	// task.state
	TaskMode      string   `json:"task_mode,omitempty"`      // modeling.task_type	// task.task_type
	TargetMetric  string   `json:"target_metric"`            // modeling.target_metric	// task.target_metric
	GpuID         []int    `json:"gpu_ids,omitempty"`        // modeling.task_state	// task.state
}
type StartModelingResponse struct {
	TrialID      int    `json:"trial_id"`           // modeling.task_id	// task.id
	State        string `json:"state"`              // modeling.task_type	// task.task_type
	GpuAuto      bool   `json:"gpu_auto,omitempty"` // modeling.task_state	// task.state
	TargetMetric string `json:"target_metric"`      // modeling.target_metric	// task.target_metric
	DatasetName  string `json:"dataset_name"`       //task.datasetId 로 검색 dataset.name
}

// 'GET' 'http://localhost:8900/api/vcls_ml/list'
type Modeling struct {
	ID    int    `json:"id"`    //modeling.id
	State string `json:"state"` // -> modeling.modeling_step
	//GPU           string         `json:"gpu"`            // -> 빼죠
	Name          string          `json:"name"`                     // task.title
	Dataset       string          `json:"dataset"`                  //task.datasetId 로 검색 dataset.name -> modeling.params 파싱하는게 나아요. task.dataId는 제거할 예정이에요
	Performance   *map[string]any `json:"performance,omitempty"`    //modeling.performance -> 이 컬럼 제거 예정입니다.
	InferenceTime *float64        `json:"inference_time,omitempty"` //??
	TargetMetric  string          `json:"target_metric"`            //task.target_metric
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

/*
performance:
V.CLS.ML
- modeling_modes의 best_model_dict에서 modeling의 target metric값에 해당 하는 1순위 모델명 검색
  - modeling_details에서 modeling_id, model(모델명) 검색
  - data_type에 threshold_{target metric}을 return 값에 맞게 변환
    - DB 데이터에서 threshold 1 -> 0.001 임
    - 따라서 100, 200, 300, 400, 500 값을 찾아서 API 문서의 포맷에 맞게 넣어주면 됨
  - data_type = test_avg_inference_time에서 "avg inference time" 값을 InferenceTime 값으로 쓰면 됨
V.CLS.SL
- modeling_modes의 best_model_dict에서 modeling의 target metric값에 해당 하는 1순위 모델명 검색
  - data_type=testset_score 를 performance
  - data_type=test_avg_inference_time에서 "avg inference time" 값을 InferenceTime 값으로 쓰면 됨
Tabular
- performance: V.CLS.SL과 동일
- InferenceTime: test_inference_time 값을 그대로 쓰면 됨
*/

// 'GET' 'http://localhost:8900/api/vcls_ml/modeling/{modeling_id}/{threshold}'
type ModelingDetail struct {
	ID    int    `json:"id"`    //modeling.id
	Name  string `json:"name"`  //task.title
	State string `json:"state"` //??
	//GPU           string           `json:"gpu"`            //??
	Dataset      string           `json:"dataset"`       //task.datasetId 로 검색 dataset.name
	TargetMetric string           `json:"target_metric"` //task.target_metric
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	ImageSize    *int             `json:"img_size,omitempty"` //task.params.image_resoultion
	Models       []*ModelingModel `json:"models,omitempty"`
}

//-> 다른 값들은 list와 동일하게

type ModelingModel struct {
	//ModelID       int            `json:"model_id"`   //modeling_models.id
	ModelName     string         `json:"model_name"` //??
	UpdatedAt     time.Time      `json:"updated_at"`
	InferenceTime float64        `json:"inference_time,omitempty"` //??
	Result        map[string]any `json:"result,omitempty"`         // performance
	Score         float64        `json:"score,omitempty"`          //??
}

/*
-> result: list에서 performance를 찾는 과정과 유사하지만 threshold 파라미터에 해당하는 값만 담아서
-> score: result중에서 모델링에 사용한 target metric 값
-> model_id -> 없애죠. load model 시에 modeling_id와 model name으로 처리합니다
*/

type TapiLoaddedModel struct {
	TestId    int    `json:"test_id"`
	ModelName string `json:"model_name"`
	ModelNum  int    `json:"model_num"`
	GPUIndex  int    `json:"gpu_index"`
	GPUID     int    `json:"gpu_id"`
	GPUUUID   string `json:"gpu_uuid"`
}

type Device struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	IP    string `json:"ip,omitempty"`
	Port  *int   `json:"port,omitempty"`
	IsUse *bool  `json:"is_use,omitempty"`
	Type  string `json:"type,omitempty"`

	GPU []*DeviceGPU `json:"gpu,omitempty"`
	//Pytorch string       `json:"pytorch,omitempty"` //??
}

type DeviceGPU struct {
	ID       int    `json:"id"`
	Index    int    `json:"index"`
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	IsUse    bool   `json:"is_use,omitempty"` // usable
	State    string `json:"state"`            // idle | train | test | load
	DeviceID int    `json:"device_id"`        // GPU machine ID
}

type DeviceModelGroup struct {
	DeviceID   int                `json:"device_id"`
	DeviceName string             `json:"device_name"`
	Models     []TapiLoaddedModel `json:"models"`
}

type SystemInformation struct {
	VERSION string    `json:"version"`
	Devices []*Device `json:"devices"`
}

type TestDTO struct {
	ModelingID int    `json:"modeling_id"`
	ModelName  string `json:"model_name"`
	GpuId      int    `json:"gpu_id"`
}

type InferenceRequest struct {
	GPUID     int                   `form:"device_id"`
	File      *multipart.FileHeader `form:"file"`
	Heatmap   string                `form:"heatmap"`
	Threshold float64               `form:"threshold"`
	ModelName string                `form:"model_name"`
}

type InferenceTabularRequest struct {
	EngineType string `form:"engine_type"`
	GPUID      int    `form:"gpu_id"`
	ModelingID int    `form:"modeling_id"`
	ModelName  string `form:"model_name"`
	XInputs    string `form:"x_input"`
	// YInputs    []map[string]interface{} `form:"y_input"`
}

package repository

import "mime/multipart"

type TestDTO struct {
	ModelingID int    `json:"modeling_id"`
	ModelName  string `json:"model_name"`
	ModelType  string `json:"model_type"`
}

type ColumnsResponse struct {
	Label         map[string]interface{} `json:"label"`
	LabelColNames []string               `json:"label_col_names"`
	InputColNames []string               `json:"input_col_names"`
	NumericalCols []string               `json:"numerical_cols"`
}

type TabularTestRequest struct {
	TestID     int                      `json:"test_id"`
	ModelingID int                      `json:"modeling_id"`
	ModelName  string                   `json:"model_name"`
	XInputs    []map[string]interface{} `json:"x_input"`
	YInputs    []map[string]interface{} `json:"y_input"`
}

type TAPITestDTO struct {
	ModelingID int    `json:"modeling_id"`
	ModelName  string `json:"model_name"`
	ModelType  string `json:"model_type"`
	GpuId      int    `json:"gpu_id"`
}

type TAPIInferenceRequest struct {
	EngineType string                `form:"engine_type"`
	GPUID      int                   `form:"gpu_id"`
	File       *multipart.FileHeader `form:"file"`
	Heatmap    string                `form:"heatmap"`
	Threshold  float64               `form:"threshold"`
	ModelName  string                `form:"model_name"`
}

type TAPIInferenceTabularRequest struct {
	EngineType string                   `form:"engine_type"`
	GPUID      int                      `form:"gpu_id"`
	ModelingID int                      `form:"modeling_id"`
	ModelName  string                   `form:"model_name"`
	XInputs    string 					`form:"x_input"`
	// YInputs    []map[string]interface{} `form:"y_input"`
}

type UnloadRequest struct {
	ModelName string `json:"model_name"`
}

package modules

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	config_service "api_server/configuration/service"
)

type DatasetResBody struct {
	FilePath string `json:"file_path"`
}
type DatasetRequestClient struct {
	EnginePath string
	StaticPath string
}

func NewDatasetRequestClient(endpoint string, datasetId int) *DatasetRequestClient {
	cf := config_service.NewStatic()
	return &DatasetRequestClient{
		EnginePath: "http://localhost:5000" + endpoint,
		StaticPath: filepath.Join(cf.Get("PATH_STATIC_TEST"), "dataset"+fmt.Sprint(datasetId)),
	}
}

func (drc *DatasetRequestClient) SendPostRequest(jsonData []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", drc.EnginePath, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

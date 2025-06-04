package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	config_service "api_server/configuration/service"
	repo "api_server/device/repository"
	"api_server/logger"
	"api_server/utils"
)

func GatherDevicesInfo() {
	svc := New(repo.New())
	if devices, err := svc.ReadActive(); err != nil {
		return
	} else {
		for _, device := range devices {
			go getDeviceInfoThroughREST(device)
		}
	}
}

func getAPIParam() string {
	cf := config_service.NewStatic()
	fileconfig := utils.ReadConfigFromFile()
	return fmt.Sprintf(
		"{\"DB_HOST\":\"%s\", \"DB_PORT\":%d, \"DB_NAME\":\"%s\", \"TEST_PATH\":\"%s\"}",
		fileconfig.DBIP,
		fileconfig.DBPORT,
		fileconfig.DB,
		cf.Get("ROOT_PATH")+"/task/",
	)
}

func updateGPUInfo(respBody []byte, device_id int) {
	engineInfo := repo.EngineInfoDTO{}
	if errUnmarshal := json.Unmarshal(respBody, &engineInfo); errUnmarshal != nil {
		logger.Error(errUnmarshal)
		return
	} else {
		ctx := context.Background()
		dao := repo.NewGPUDAO()

		engineInfo.DeviceID = device_id
		dao.UpdateAllDisUse(ctx)
		dao.UpsertMany(ctx, engineInfo)
	}
}

func getDeviceInfoThroughREST(device *repo.DeviceDTO) {
	params := getAPIParam()

	reqBody := bytes.NewBufferString(params)

	if resp, err := http.Post(
		"http://"+device.IP+":"+strconv.Itoa(*device.Port)+"/api/sys",
		"application/json",
		reqBody); err != nil {
		logger.Error(err)
	} else if respBody, err := io.ReadAll(resp.Body); err != nil {
		logger.Error(err)
	} else {
		updateGPUInfo(respBody, device.ID)
	}
}

package service

import (
	"archive/zip"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	config_service "api_server/configuration/service"
	"api_server/ent"
	"api_server/ent/modeling"
	"api_server/ent/modelingmodels"
	"api_server/ent/task"
	"api_server/logger"
	"api_server/utils"
)

type DownloadService struct {
	entClient *ent.Client
}

var once sync.Once
var instance *DownloadService

func New() *DownloadService {
	once.Do(func() {
		logger.Debug("Download Service intance")
		instance = &DownloadService{
			entClient: utils.GetEntClient(),
		}
	})
	return instance
}

func (svc *DownloadService) GetModelZip(ctx context.Context, modeling_id int, model_name string) (string, *logger.Report) {
	downloadModel, err := svc.entClient.ModelingModels.
		Query().
		Where(modelingmodels.And(
			modelingmodels.ModelingID(modeling_id),
			modelingmodels.DataType("best_model_dict"),
		)).
		First(ctx)
	if err != nil {
		return "", logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	task, err := svc.entClient.Task.Query().
		Where(task.HasModelingsWith(modeling.ID(modeling_id))).
		Only(ctx)
	if err != nil {
		return "", logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	if downloadModel == nil || len(downloadModel.Data) == 0 {
		return "", logger.CreateReport(&logger.CODE_INVALID_METRIC, nil)
	}

	var bestModelDict map[string]map[string][]interface{}
	if err := json.Unmarshal([]byte(downloadModel.Data), &bestModelDict); err != nil {
		return "", logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, err)
	}

	modelPath := ""
	for _, models := range bestModelDict {
		for _, v := range models {
			if strings.HasSuffix(v[0].(string), model_name+".kaier") {
				modelPath = v[0].(string)
			}
		}
	}

	return svc.createZipFile(task.EngineType, modelPath)
}

func (svc *DownloadService) createZipFile(engineType string, modelPath string) (string, *logger.Report) {
	cf := config_service.NewStatic()

	modelPath = strings.ReplaceAll(modelPath, "\\", "/")
	modelPath = strings.ReplaceAll(modelPath, "//", "/")
	modelDir, modelFileName := filepath.Split(modelPath)

	zipPath := filepath.Join(cf.Get("PATH_STATIC_TEST"), modelFileName+".zip")

	archive, err := os.Create(zipPath)
	if err != nil {
		return "", logger.CreateReport(&logger.CODE_ADD_ZIP, err)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	if err := svc.addFileToZip(zipWriter, filepath.Join(modelDir, modelFileName), modelFileName); err != nil {
		return "", logger.CreateReport(&logger.CODE_ADD_ZIP, nil)
	}

	if engineType == utils.JOB_TYPE_VISION_CLS_SL || engineType == utils.JOB_TYPE_VISION_CLS_ML {
		configname := strings.ReplaceAll(modelFileName, ".kaier", ".kaml")
		if err := svc.addFileToZip(zipWriter, filepath.Join(modelDir, "config.kaml"), configname); err != nil {
			return "", logger.CreateReport(&logger.CODE_ADD_ZIP, nil)
		}
	} else if engineType == utils.JOB_TYPE_VISION_AD {
		logger.Debug("not implemented")
		// TODO
	} else if engineType == utils.JOB_TYPE_TABLE_CLS {
		logger.Debug("not implemented")
		// TODO
	} else if engineType == utils.JOB_TYPE_TABLE_REG {
		logger.Debug("not implemented")
		// TODO
	}

	return zipPath, nil
}

func (svc *DownloadService) addFileToZip(zipWriter *zip.Writer, filePath, fileName string) *logger.Report {
	// _filepath, err := svc.addDriveIfNeeded(filePath)
	// if err != nil {
	// 	return logger.CreateReport(&logger.CODE_FILE_OPEN, err)
	// }

	file, err := os.Open(filePath)
	if err != nil {
		return logger.CreateReport(&logger.CODE_FILE_OPEN, err)
	}
	defer file.Close()

	zipFile, err := zipWriter.Create(fileName)
	if err != nil {
		return logger.CreateReport(&logger.CODE_ENTRY_ZIP, err)
	}

	if _, err := io.Copy(zipFile, file); err != nil {
		return logger.CreateReport(&logger.CODE_COPY_ZIP, err)
	}

	return nil
}

// func (svc *DownloadService) addDriveIfNeeded(filePath string) (string, error) {
// 	/*
// 		hrs monkey patch
// 	*/
// 	if runtime.GOOS != "windows" {
// 		return filePath, nil
// 	}

// 	if len(filePath) >= 2 && filePath[1] == ':' {
// 		return filePath, nil
// 	}

// 	return "D:" + filePath, nil
// }

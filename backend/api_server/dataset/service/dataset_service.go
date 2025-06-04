package service

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"api_server/dataset/common"
	"api_server/dataset/modules"
	repo "api_server/dataset/repository"
	"api_server/logger"
	"api_server/utils"
)

type DatasetServiceInterface interface {
	ViewDatasets(datasetType []string, page int) (*repo.GetDatasetsDTO, *logger.Report)
	ViewDatasetForTAPI(engine string) ([]*repo.DatasetDTO, *logger.Report)
	ViewTableColumn(parent_id int) ([]string, *logger.Report)
	ViewClasses(parent_id int, engine_type string) (*repo.ImageClassModel, *logger.Report)
	GetIsTiff(path string) bool
	RemoveDataset(id int) *logger.Report
	GetDataStatistics(id int) (*repo.DatasetStatistics, *logger.Report)

	// GetDataStatByTypeFromJSON은 주어진 ID와 통계 타입(statType)에 해당하는 JSON 통계 파일 경로를 반환합니다.
	//   - id: 데이터셋의 고유 ID
	//   - statType: 가져오려는 통계의 타입 (예: "category", "summary" 등)
	//
	// 반환값:
	//   - string: 해당 통계 타입의 JSON 파일 경로
	//   - *logger.Report: 오류 발생 시 리포트, 없으면 nil
	GetDataStatByTypeFromJSON(id int, statType string) (string, *logger.Report)
	SelectTabularDatasetCompareNumerical(repo.CompareFeaturesStatics) (*repo.CompareNumericalFeaturesStatics, *logger.Report)
	SelectTabularDatasetCompareCategorical(repo.CompareFeaturesStatics) (*repo.CompareCategoricalFeaturesStatics, *logger.Report)
	SelectTabularDatasetCompareCategoricalNumerical(repo.CompareFeaturesStatics) (*repo.CompareCategoricalNumericalFeaturesStatics, *logger.Report)

	// GetDatasetType은 주어진 ID를 기반으로 데이터셋의 엔진 타입 리스트를 반환합니다.
	// 매개변수:
	//	- id: 데이터셋의 고유 ID
	//
	// 반환값:
	//	- []string: 데이터셋이 사용하는 엔진 타입 리스트
	//	- *logger.Report: 오류 발생 시 리포트, 없으면 nil
	GetDatasetType(id int) ([]string, *logger.Report)

	// GetDatasetType은 주어진 ID를 기반으로 데이터셋의 엔진 타입 리스트를 반환합니다.
	//	- id: 데이터셋의 고유 ID
	//
	// 반환값:
	//	- []string: 데이터셋이 사용하는 엔진 타입 리스트
	//	- *logger.Report: 오류 발생 시 리포트, 없으면 nil
	ReadDataset(id int) (*repo.DatasetDTO, *logger.Report)

	GetDatasetByName(name string) (*repo.DatasetDTO, *logger.Report)
}

type DatasetService struct {
	ctx             context.Context
	datasetWatcher  modules.DatasetWatcherInterface
	datasetAnalyzer modules.DatasetAnalyzerInterface
	datasetDAO      repo.DatasetDAOInterface
}

var datasetServiceInstance *DatasetService

func NewDatasetService(datasetWatcher modules.DatasetWatcherInterface, datasetAnalyzer modules.DatasetAnalyzerInterface, datasetDAO repo.DatasetDAOInterface) *DatasetService {
	if datasetServiceInstance == nil {
		datasetServiceInstance = &DatasetService{
			ctx:             context.Background(),
			datasetWatcher:  datasetWatcher,
			datasetAnalyzer: datasetAnalyzer,
			datasetDAO:      datasetDAO,
		}
	}

	return datasetServiceInstance
}

func (svc *DatasetService) ViewDatasets(datasetType []string, page int) (*repo.GetDatasetsDTO, *logger.Report) {
	// if isTest {
	// 	if datasets, pageCount, err := svc.datasetDAO.SelectTestableDatasets(svc.ctx, datasetType, page); err != nil {
	// 		return nil, err
	// 	} else {
	// 		return repo.ConvertDatasetEntsToGetDatasetsDTOs(datasets, pageCount), nil
	// 	}
	// } else {
	// 	if datasets, pageCount, err := svc.datasetDAO.SelectDatasets(svc.ctx, datasetType, page); err != nil {
	// 		return nil, err
	// 	} else {
	// 		return repo.ConvertDatasetEntsToGetDatasetsDTOs(datasets, pageCount), nil
	// 	}
	// }

	if datasets, pageCount, hasMore, nextPage, err := svc.datasetDAO.SelectDatasets(svc.ctx, datasetType, page); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertDatasetEntsToGetDatasetsDTOs(datasets, pageCount, hasMore, nextPage), nil
	}
}

func (svc *DatasetService) ViewDatasetForTAPI(engine string) ([]*repo.DatasetDTO, *logger.Report) {
	if engine != utils.JOB_TYPE_VISION_CLS_ML {
		return nil, logger.CreateReport(&logger.CODE_API_PARAM_ENGINE, nil)
	}

	svc.datasetWatcher.DetectDatasetModification()

	result, r := svc.datasetDAO.SelectDatasetForAPI(svc.ctx, 0, utils.DATA_TYPE_IMG)
	if r != nil {
		return nil, r
	}

	return repo.ConvertDatasetEntsToDTOs(result), nil
}

func (svc *DatasetService) ViewTableColumn(parent_id int) ([]string, *logger.Report) {
	datasets, r := svc.datasetDAO.SelectDataSetByParentID(svc.ctx, parent_id)
	if r != nil {
		return nil, r
	}

	if len(datasets) < 1 {
		datasets, r = svc.datasetDAO.SelectDataSetByID(svc.ctx, parent_id)
		if r != nil {
			return nil, r
		}

		if len(datasets) < 1 {
			return nil, logger.CreateReport(&logger.CODE_DATA_IMAGE_TYPE, nil)
		}
	}

	if datasets[0].DataType != utils.DATA_TYPE_TABLE {
		r := logger.CreateReport(&logger.CODE_DATA_TABLE_TYPE, nil)
		return nil, r
	}

	return svc.getTableColumns(datasets[0].Path)
}

func (svc *DatasetService) ViewClasses(parent_id int, engine_type string) (*repo.ImageClassModel, *logger.Report) {
	datasets, r := svc.datasetDAO.SelectDataSetByParentID(svc.ctx, parent_id)
	if r != nil {
		return nil, r
	}

	// folder test 인 경우 parent_id가 아님
	if len(datasets) < 1 {
		datasets, r = svc.datasetDAO.SelectDataSetByID(svc.ctx, parent_id)
		if r != nil {
			return nil, r
		}

		if len(datasets) < 1 {
			return nil, logger.CreateReport(&logger.CODE_DATA_IMAGE_TYPE, nil)
		}
	}

	if datasets[0].DataType != utils.DATA_TYPE_IMG {
		return nil, logger.CreateReport(&logger.CODE_DATA_IMAGE_TYPE, nil)
	}

	var classes []string
	if engine_type == utils.JOB_TYPE_VISION_CLS_ML {
		dataset, r := svc.datasetDAO.SelectDataSetByID(svc.ctx, parent_id)
		if r != nil {
			return nil, r
		} else if len(dataset) < 1 {
			r := logger.CreateReport(&logger.CODE_DATA_IMAGE_CLASS, nil)
			return nil, r
		}

		labels, err := common.GetLabelNames(dataset[0].Path)
		if err != nil || len(labels) < 1 {
			r := logger.CreateReport(&logger.CODE_DATA_IMAGE_CLASS, err)
			return nil, r
		}

		classes = labels[0:]
	} else {
		// Read Sub Folders
		dirs, err := common.GetChildDirNames(datasets[0].Path)
		if err != nil || len(dirs) < 1 {
			r := logger.CreateReport(&logger.CODE_DATA_IMAGE_CLASS, err)
			return nil, r
		}

		classes = dirs[0:]
	}

	imageClass := repo.ImageClassModel{Classes: classes, IsTiff: svc.GetIsTiff(datasets[0].Path)}
	return &imageClass, nil
}

func (svc *DatasetService) GetIsTiff(path string) bool {
	classDirs, errDir := common.GetChildDirDataset(path)
	if errDir != nil || len(classDirs) < 1 {
		return false
	}

	parent, err := os.Open(classDirs[0].Path)
	if err != nil {
		return false
	}
	defer parent.Close()

	childList, err := parent.Readdir(-1)
	if err != nil {
		return false
	}

	for _, child := range childList {
		filename := strings.ToLower(child.Name())
		if !child.IsDir() && (strings.HasSuffix(filename, ".tiff") || strings.HasSuffix(filename, ".tif")) {
			return true
		}
	}

	return false
}

func (svc *DatasetService) RemoveDataset(id int) *logger.Report {
	dsPath, err := svc.datasetDAO.SelectDataPathByDataSetId(svc.ctx, id)
	if err != nil {
		return err
	}

	errRemovePath := os.RemoveAll(dsPath)
	if errRemovePath != nil {
		return logger.CreateReport(&logger.CODE_DIR_NOT_EXIST, errRemovePath)
	}

	return svc.datasetDAO.DeleteDataset(svc.ctx, id)
}

func (svc *DatasetService) GetDataStatistics(id int) (*repo.DatasetStatistics, *logger.Report) {
	dataset, r := svc.datasetDAO.SelectStatistics(svc.ctx, id)
	if r != nil || len(dataset.Stat) < 1 {
		return nil, r
	}

	stat := &repo.DatasetStatistics{}
	if errJson := json.Unmarshal([]byte(dataset.Stat[0]), &stat); errJson != nil {
		return nil, logger.CreateReport(&logger.CODE_JSON_UNMARSHAL, errJson)
	}

	return stat, nil
}

func (svc *DatasetService) getTableColumns(rootpath string) ([]string, *logger.Report) {
	files, _ := utils.ReadFiles(rootpath, []string{utils.EXT_CSV, utils.EXT_XLS, utils.EXT_XLSX}, nil)
	if len(files) < 1 {
		return nil, logger.CreateReport(&logger.CODE_FILE_NOT_EXIST, nil)
	}

	for _, file := range files {
		rows, err := utils.ReadTabularFile(filepath.Join(rootpath, file.Name()))
		// TODO: 팀장님께 검사
		if err == nil && len(rows[0]) > 1 {
			return rows[0], nil
		}
	}

	return nil, logger.CreateReport(&logger.CODE_FILE_NOT_EXIST, nil)
}

func (svc *DatasetService) SelectTabularDatasetCompareNumerical(featureInfo repo.CompareFeaturesStatics) (*repo.CompareNumericalFeaturesStatics, *logger.Report) {
	compareNumericalFeaturesStatics, r := svc.datasetAnalyzer.CompareNumericalFeature(featureInfo.DatasetId, featureInfo.Feature1, featureInfo.Feature2)
	if r != nil {
		return nil, r
	}

	return compareNumericalFeaturesStatics, nil
}

func (svc *DatasetService) SelectTabularDatasetCompareCategorical(featureInfo repo.CompareFeaturesStatics) (*repo.CompareCategoricalFeaturesStatics, *logger.Report) {
	compareCategoricalFeatureStatics, r := svc.datasetAnalyzer.CompareCategoricalFeature(featureInfo.DatasetId, featureInfo.Feature1, featureInfo.Feature2)
	if r != nil {
		return nil, r
	}

	return compareCategoricalFeatureStatics, nil
}

func (svc *DatasetService) SelectTabularDatasetCompareCategoricalNumerical(featureInfo repo.CompareFeaturesStatics) (*repo.CompareCategoricalNumericalFeaturesStatics, *logger.Report) {
	compareCategoricalNumericalFeatureStatics, r := svc.datasetAnalyzer.CompareCategoricalNumericalFeature(featureInfo.DatasetId, featureInfo.Feature1, featureInfo.Feature2)
	if r != nil {
		return nil, r
	}

	return compareCategoricalNumericalFeatureStatics, nil
}

// GetDatasetType은 주어진 ID를 기반으로 데이터셋의 엔진 타입 리스트를 반환합니다.
func (svc *DatasetService) GetDatasetType(id int) ([]string, *logger.Report) {
	dataset, r := svc.datasetDAO.SelectDataSetByID(svc.ctx, id)
	if r != nil {
		return []string{}, r
	}

	return dataset[0].Engine, nil
}

// GetDataStatByTypeFromJSON은 주어진 ID와 통계 타입(statType)에 해당하는 JSON 통계 파일 경로를 반환합니다.
func (svc *DatasetService) GetDataStatByTypeFromJSON(id int, statType string) (string, *logger.Report) {
	dataset, r := svc.datasetDAO.SelectStatistics(svc.ctx, id)
	if r != nil || len(dataset.StatPath) < 1 {
		return "", r
	}

	file := filepath.Join(dataset.StatPath, statType+".json")

	return file, nil
}

// ReadDataset은 주어진 ID를 기반으로 데이터셋을 조회하고 DTO 형식으로 반환합니다.
func (svc *DatasetService) ReadDataset(id int) (*repo.DatasetDTO, *logger.Report) {
	dataset, r := svc.datasetDAO.SelectDataSetByID(svc.ctx, id)
	if r != nil {
		return nil, r
	}

	if len(dataset) < 1 {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, nil)
	}

	return repo.ConvertDatasetEntToDTO(dataset[0]), nil
}

func (svc *DatasetService) GetDatasetByName(unique_name string) (*repo.DatasetDTO, *logger.Report) {
	dataset, r := svc.datasetDAO.SelectDataSetByName(svc.ctx, unique_name)
	if r != nil {
		return nil, r
	}

	if len(dataset) < 1 {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, nil)
	}

	return repo.ConvertDatasetEntToDTO(dataset[0]), nil
}

package modules

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"

	repo "api_server/dataset/repository"

	"api_server/utils"
)

type DatasetValidatorInterface interface {
	Validate(dr_id int)
}

type DatasetValidator struct {
	ctx        context.Context
	datasetDAO repo.DatasetDAOInterface
	dataFormat string
}

func NewDatasetValidator(datasetDAO repo.DatasetDAOInterface) *DatasetValidator {
	return &DatasetValidator{
		ctx:        context.Background(),
		datasetDAO: datasetDAO,
	}
}

func (v *DatasetValidator) Validate(dr_id int) {
	datasetEnts := v.datasetDAO.SelectDatasetsByDRID(v.ctx, dr_id)
	datasets := repo.ConvertDatasetEntsToDTOs(datasetEnts)

	for _, dataset := range datasets {
		v.identifyDataType(dataset)

		v.dataFormat = v.identifyKaierFormat(dataset)
		dataset.IsValid = v.dataFormat != utils.DATA_FORMAT_NONE

		v.identifyEngineType(dataset)

		v.checkTestablePath(dataset)

		if !dataset.IsValid && !dataset.IsTrainable {
			dataset.Description = "The dataset structure is incomplete"
		} else if !dataset.IsTrainable {
			dataset.Description = "The dataset is incomplete"
		}

		v.updateDatasetValidation(dataset)
	}
}

func (v *DatasetValidator) identifyDataType(dataset *repo.DatasetDTO) string {
	dataset.DataType = v.checkDatatype(dataset.Path)
	if datasetEnts, r := v.datasetDAO.SelectDataSetByParentID(v.ctx, dataset.ID); r == nil {
		dataset.Childs = repo.ConvertDatasetEntsToDTOs(datasetEnts)
		for _, child := range dataset.Childs {
			dataType := v.identifyDataType(child)
			if dataset.DataType == utils.DATA_TYPE_INVALID {
				dataset.DataType = dataType
			}
		}
	}
	return dataset.DataType
}

func (v *DatasetValidator) checkDatatype(path string) string {
	files, err := utils.ReadFiles(path, nil, []string{".json", ".txt"})
	if err != nil {
		return utils.DATA_TYPE_INVALID
	}
	if len(files) == 0 {
		return utils.DATA_TYPE_INVALID
	}

	for _, file := range files {
		file_path := filepath.Join(path, file.Name())
		if utils.IsImageFile(file_path) {
			return utils.DATA_TYPE_IMG
		}
		if utils.IsTabularFile(file_path) {
			return utils.DATA_TYPE_TABLE
		}
	}

	return utils.DATA_TYPE_INVALID
}

func (v *DatasetValidator) identifyKaierFormat(dataset *repo.DatasetDTO) string {
	if dirs, err := utils.ReadDirs(dataset.Path); err != nil {
		return utils.DATA_FORMAT_NONE
	} else {
		isTrain := false
		isValid := false
		isTest := false
		for _, d := range dirs {
			if d.Name() == "train" {
				isTrain = true
			} else if d.Name() == "valid" {
				isValid = true
			} else if d.Name() == "test" {
				isTest = true
			}
		}

		if !isTrain {
			return utils.DATA_FORMAT_NONE
		}
		if isValid && isTest {
			return utils.DATA_FORMAT_KAIER_TVT
		}
		if isValid && !isTest {
			return utils.DATA_FORMAT_KAIER_TV
		}
		if !isValid && isTest {
			return utils.DATA_FORMAT_KAIER_TT
		}
		if !isValid && !isTest {
			return utils.DATA_FORMAT_KAIER_TRAIN
		}
	}

	return utils.DATA_FORMAT_NONE
}

func (v *DatasetValidator) identifyEngineType(dataset *repo.DatasetDTO) {
	engineType := []string{}

	if v.validateVisionClsSlDataset(dataset) {
		engineType = append(engineType, utils.JOB_TYPE_VISION_CLS_SL)
	}

	if v.validateVisionClsMlDataset(dataset) {
		engineType = append(engineType, utils.JOB_TYPE_VISION_CLS_ML)
	}

	if v.validateVisionADDataset(dataset) {
		engineType = append(engineType, utils.JOB_TYPE_VISION_AD)
	}

	if v.validateTabularDataset(dataset) {
		engineType = append(engineType, utils.JOB_TYPE_TABLE_CLS, utils.JOB_TYPE_TABLE_REG)
	}

	if len(engineType) < 1 {
		engineType = append(engineType, utils.JOB_TYPE_INVALID)
	} else {
		dataset.IsTrainable = true
	}

	dataset.Engine = engineType
}

func (v *DatasetValidator) validateVisionClsSlDataset(dataset *repo.DatasetDTO) bool {
	if v.dataFormat == utils.DATA_FORMAT_NONE {
		return false
	}

	if ok := v.hasSufficientFiles(filepath.Join(dataset.Path, "train")); !ok {
		return false
	}

	if v.dataFormat == utils.DATA_FORMAT_KAIER_TV || v.dataFormat == utils.DATA_FORMAT_KAIER_TVT {
		if ok := v.hasSufficientFiles(filepath.Join(dataset.Path, "valid")); !ok {
			return false
		}
	}

	if v.dataFormat == utils.DATA_FORMAT_KAIER_TT || v.dataFormat == utils.DATA_FORMAT_KAIER_TVT {
		if ok := v.hasSufficientFiles(filepath.Join(dataset.Path, "test")); !ok {
			return false
		}
	}

	return true
}

func (v *DatasetValidator) validateVisionClsMlDataset(dataset *repo.DatasetDTO) bool {
	if v.dataFormat == utils.DATA_FORMAT_NONE {
		return false
	}

	file, err := os.Open(filepath.Join(dataset.Path, "label.txt"))
	if err != nil {
		return false
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		labelInfo := strings.Split(line, " ")
		if _, err := os.Stat(filepath.Join(dataset.Path, labelInfo[0])); os.IsNotExist(err) {
			return false
		}

		if len(strings.Join(labelInfo[1:], "")) == 0 {
			return false
		}

		count++
		if count > 3 {
			return true
		}
	}

	return false
}

func (v *DatasetValidator) validateVisionADDataset(dataset *repo.DatasetDTO) bool {
	if v.dataFormat == utils.DATA_FORMAT_NONE {
		return false
	}

	dirs, _ := utils.ReadDirs(filepath.Join(dataset.Path, "train"))

	for _, d := range dirs {
		if d.Name() == "normal" {
			if v.hasSufficientFiles(filepath.Join(dataset.Path, "train")) {
				return true
			}
		}
	}

	return false
}

func (v *DatasetValidator) validateTabularDataset(dataset *repo.DatasetDTO) bool {
	if v.dataFormat == utils.DATA_FORMAT_NONE {
		return false
	}

	files, _ := utils.ReadFiles(filepath.Join(dataset.Path, "train"), []string{utils.EXT_CSV, utils.EXT_XLS, utils.EXT_XLSX}, nil)
	if len(files) < 1 {
		return false
	}

	for _, file := range files {
		path := filepath.Join(dataset.Path, "train", file.Name())
		rows, err := utils.ReadTabularFile(path)

		for i := 0; i < 5; i++ {
			if err != nil || len(rows[i]) < 2 {
				return false
			}
		}
	}

	return true
}

func (v *DatasetValidator) hasSufficientFiles(path string) bool {
	dirs, _ := utils.ReadDirs(path)
	if len(dirs) < 1 {
		return false
	}

	for _, d := range dirs {
		files, err := utils.ReadFiles(filepath.Join(path, d.Name()), nil, nil)
		if err != nil {
			return false
		}

		if len(files) < 3 {
			return false
		}
	}

	return true
}

func (v *DatasetValidator) checkTestablePath(dataset *repo.DatasetDTO) {
	if dataset.IsTrainable {
		dataset.IsTestable = true
	} else {
		dataset.IsTestable = v.isTestablePath(dataset.Path, dataset.DataType)
	}

	for _, ds := range dataset.Childs {
		v.checkTestablePath(ds)
	}
}

func (v *DatasetValidator) isTestablePath(path string, dataType string) bool {
	if dataType == utils.DATA_TYPE_IMG {
		return v.isTestableImagePath(path)
	}
	if dataType == utils.DATA_TYPE_TABLE {
		result := v.isTestableTabularPath(path)
		return result
	}

	return false
}

func (v *DatasetValidator) isTestableImagePath(path string) bool {
	files, err := utils.ReadFiles(path, nil, nil)
	if err != nil {
		return false
	}

	if len(files) < 1 {
		return false
	}

	return true
}

func (v *DatasetValidator) isTestableTabularPath(path string) bool {
	files, _ := utils.ReadFiles(path, []string{utils.EXT_CSV, utils.EXT_XLS, utils.EXT_XLSX}, nil)
	if len(files) < 1 {
		return false
	}

	for _, file := range files {
		path := filepath.Join(path, file.Name())
		rows, err := utils.ReadTabularFile(path)
		if err != nil {
			return false
		}

		if len(rows) < 2 {
			return false
		}

		size := min(5, len(rows[0]))
		for i := 0; i < size; i++ {
			if len(rows[i]) < 2 {
				return false
			}
		}
	}

	return true
}

func (v *DatasetValidator) updateDatasetValidation(dataset *repo.DatasetDTO) {
	v.datasetDAO.UpdateValidation(v.ctx, *dataset)
	for _, child := range dataset.Childs {
		v.updateDatasetValidation(child)
	}
}

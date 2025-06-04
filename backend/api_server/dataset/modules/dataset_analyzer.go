package modules

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"

	go_stats "github.com/aclements/go-moremath/stats"
	stats "github.com/montanaflynn/stats"

	repo "api_server/dataset/repository"
	"api_server/logger"
	"api_server/utils"
)

type DatasetAnalyzerInterface interface {
	Analyze(dr_id int)

	CompareNumericalFeature(datasetId int, feature1 string, feature2 string) (*repo.CompareNumericalFeaturesStatics, *logger.Report)
	CompareCategoricalFeature(datasetId int, feature1 string, feature2 string) (*repo.CompareCategoricalFeaturesStatics, *logger.Report)
	CompareCategoricalNumericalFeature(datasetId int, feature1 string, feature2 string) (*repo.CompareCategoricalNumericalFeaturesStatics, *logger.Report)
}

type DatasetAnalyzer struct {
	ctx        context.Context
	datasetDAO repo.DatasetDAOInterface

	// using count categorical feature
	// detailCount              map[string]map[string]map[string]int
	feature map[string][]string
	// categoricalFeatureStat   []map[string]string
	// categoricalFeatureDetail map[string][]map[string]string
}

var TVT_NAMES = []string{"train", "valid", "test"}
var STAT_CATEGORY_NAMES = []string{"mean", "median", "min", "max", "stdev"}
var TABULAR_EXTENSIONS = []string{utils.EXT_CSV, utils.EXT_XLS, utils.EXT_XLSX}

type Sample struct {
	Xs      []float64
	Weights []float64
	Sorted  bool
}

type KDEKernel int
type KDE struct {
	Sample []float64
	Kernel KDEKernel
}

func NewDatasetAnalyzer(datasetDAO repo.DatasetDAOInterface) *DatasetAnalyzer {
	return &DatasetAnalyzer{
		ctx:        context.Background(),
		datasetDAO: datasetDAO,
	}
}

func (da *DatasetAnalyzer) Analyze(dr_id int) {
	datasetEnts := da.datasetDAO.SelectDatasetsByDRID(da.ctx, dr_id)
	datasets := repo.ConvertDatasetEntsToDTOs(datasetEnts)

	for _, dataset := range datasets {
		da.analyzeDataset(dataset)
	}
}

func (da *DatasetAnalyzer) analyzeDataset(dataset *repo.DatasetDTO) {
	var classStatics *repo.ClassStatics
	var resolutionStatics *repo.ResolutionStatics
	var noneTypeStat *repo.NoneTypeStat

	if slices.Contains(dataset.Engine, utils.JOB_TYPE_VISION_CLS_ML) {
		classStatics = da.multiClass(dataset.Path)
		resolutionStatics = da.multilabelResolution(dataset.Path)
	} else if slices.Contains(dataset.Engine, utils.JOB_TYPE_VISION_CLS_SL) {
		classStatics = da.singleClass(dataset.Path)
		resolutionStatics = da.singlelabelResolution(dataset.Path)
	} else if len(dataset.Engine) < 1 || slices.Contains(dataset.Engine, utils.JOB_TYPE_INVALID) {
		noneTypeStat = da.countNonetypeDataset(dataset.Path, dataset.DataType)
	}

	da.feature = da.readFeatures(filepath.Join(dataset.Path, "train"))

	// TODO : compress data
	// fmt.Println(dataset.Name)
	// r1 := da.analyzeNumericalFeature(dataset.Path)
	// if jsonBytes, err := json.Marshal(r1); err == nil {
	// 	fmt.Println(len(jsonBytes))
	// }

	// r2 := da.computeNumericalHeatmap(dataset.Path)
	// if jsonBytes, err := json.Marshal(r2); err == nil {
	// 	fmt.Println(len(jsonBytes))
	// }
	// fmt.Println("----------")

	stat := repo.DatasetStatistics{
		ClassStatics:                classStatics,
		MultiLabelClassStatics:      da.multilabelStat(dataset.Path),
		ResolutionStatics:           resolutionStatics,
		Features:                    da.feature,
		CategoricalFeatureStatics:   da.analyzeCategoricalFeatureOnEngine(dataset.Path, dataset.ID),
		CategoricalHeatmap:          da.computeCategoricalHeatmapOnEngine(dataset.Path, dataset.ID),
		CategoricalNumericalHeatmap: da.computeCategoricalNumericalHeatmapOnEngine(dataset.Path, dataset.ID),
		NumericalHeatmap:            da.computeNumericalHeatmapOnEngine(dataset.Path, dataset.ID),
		NumericalFeatureStatics:     da.analyzeNumericalFeatureOnEngine(dataset.Path, dataset.ID),

		NoneTypeStat: noneTypeStat,
	}

	reqClient := NewDatasetRequestClient("", dataset.ID)

	if jsonBytes, err := json.Marshal(stat); err == nil {
		da.datasetDAO.UpdateStat(da.ctx, dataset.ID, string(jsonBytes))
		da.datasetDAO.UpdateStatPath(da.ctx, dataset.ID, reqClient.StaticPath)
	}
}

func (da *DatasetAnalyzer) multiClass(path string) *repo.ClassStatics {
	file, err := os.Open(filepath.Join(path, "label.txt"))
	if err != nil {
		return nil
	}
	defer file.Close()

	class := make(map[string]map[string]int)
	count := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		labelInfo := strings.Split(line, " ")

		refinePath := utils.RefinePathSeparator(labelInfo[0])
		directory := strings.Split(refinePath, string(os.PathSeparator))[0]

		if class[directory] == nil {
			class[directory] = make(map[string]int)
		}

		for _, label := range labelInfo[1:] {
			if label == "" {
				continue
			}
			class[directory][label]++
		}
		count[directory]++
	}

	return &repo.ClassStatics{Class: class, Count: count}
}

func (da *DatasetAnalyzer) singleClass(path string) *repo.ClassStatics {
	tvtDirs, err := utils.ReadDirs(path)
	if err != nil {
		return nil
	}

	class := make(map[string]map[string]int)
	count := make(map[string]int)

	for _, tvt := range tvtDirs {
		class[tvt.Name()] = make(map[string]int)
		count[tvt.Name()] = 0

		directories, err := utils.ReadDirs(filepath.Join(path, tvt.Name()))
		if err != nil {
			continue
		}

		for _, directory := range directories {
			files, err := utils.ReadFiles(filepath.Join(path, tvt.Name(), directory.Name()), nil, nil)
			if err != nil {
				continue
			}

			class[tvt.Name()][directory.Name()] = len(files)
			count[tvt.Name()] += len(files)
		}
	}

	return &repo.ClassStatics{Class: class, Count: count}
}

func (da *DatasetAnalyzer) multilabelResolution(path string) *repo.ResolutionStatics {
	file, err := os.Open(filepath.Join(path, "label.txt"))
	if err != nil {
		return nil
	}
	defer file.Close()

	resolution := make(map[string]map[int]map[int]int)
	count := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		labelInfo := strings.Split(line, " ")

		refinePath := utils.RefinePathSeparator(labelInfo[0])
		directory := strings.Split(refinePath, string(os.PathSeparator))[0]

		if resolution[directory] == nil {
			resolution[directory] = make(map[int]map[int]int)
		}

		openedFile, err := os.Open(filepath.Join(path, labelInfo[0]))
		if err != nil {
			continue
		}

		image, _, err := image.DecodeConfig(openedFile)
		if err != nil {
			continue
		}

		if resolution[directory][image.Width] == nil {
			resolution[directory][image.Width] = make(map[int]int)
		}
		resolution[directory][image.Width][image.Height]++
		count[directory]++
	}

	return &repo.ResolutionStatics{Resolution: resolution, Count: count}
}

func (da *DatasetAnalyzer) singlelabelResolution(path string) *repo.ResolutionStatics {
	tvtDirs, err := utils.ReadDirs(path)
	if err != nil {
		return nil
	}

	resolution := make(map[string]map[int]map[int]int)
	count := make(map[string]int)

	for _, tvt := range tvtDirs {
		resolution[tvt.Name()] = make(map[int]map[int]int)
		count[tvt.Name()] = 0

		tvtPath := filepath.Join(path, tvt.Name())
		directories, err := utils.ReadDirs(tvtPath)
		if err != nil {
			logger.Error("Failed to read dir : ", err)
		}

		for _, directory := range directories {
			da.singlelabelResolutionPerDir(filepath.Join(tvtPath, directory.Name()), tvt.Name(), resolution, count)
		}
	}

	return &repo.ResolutionStatics{Resolution: resolution, Count: count}
}

func (da *DatasetAnalyzer) singlelabelResolutionPerDir(path string, dirName string, resolution map[string]map[int]map[int]int, count map[string]int) {
	files, err := utils.ReadFiles(path, nil, []string{".txt", ".json"})
	if err != nil {
		return
	}

	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		openedFile, err := os.Open(filePath)
		if err != nil {
			logger.Error("Failed to open file: ", err)
			continue
		}
		defer openedFile.Close()

		// 파일 확장자 확인
		ext := strings.ToLower(filepath.Ext(file.Name()))

		var width, height int

		if ext == ".tiff" || ext == ".tif" {
			width, height, err = da.getTiffDimensions(filePath)
			if err != nil {
				logger.Error("Failed to get TIFF dimensions: ", err)
				continue
			}
		} else {
			imgConfig, _, err := image.DecodeConfig(openedFile)
			if err != nil {
				logger.Error("Failed to decode image: ", err)
				continue
			}
			width = imgConfig.Width
			height = imgConfig.Height
		}

		if resolution[dirName] == nil {
			resolution[dirName] = make(map[int]map[int]int)
		}

		if resolution[dirName][width] == nil {
			resolution[dirName][width] = make(map[int]int)
		}

		resolution[dirName][width][height]++
		count[dirName]++
	}
}

func (da *DatasetAnalyzer) multilabelStat(path string) *repo.MultiLabelClassStatics {
	file, err := os.Open(filepath.Join(path, "label.txt"))
	if err != nil {
		return nil
	}
	defer file.Close()

	class := make(map[string]map[int]int)
	count := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		labelInfo := strings.Split(line, " ")

		refinePath := utils.RefinePathSeparator(labelInfo[0])
		directory := strings.Split(refinePath, string(os.PathSeparator))[0]

		if class[directory] == nil {
			class[directory] = make(map[int]int)
		}
		class[directory][len(labelInfo[1:])]++
		count[directory]++
	}

	return &repo.MultiLabelClassStatics{Class: class, Count: count}
}

func (da *DatasetAnalyzer) readFeatures(path string) map[string][]string {
	files, err := utils.ReadFiles(path, TABULAR_EXTENSIONS, nil)
	if err != nil || len(files) < 1 {
		return nil
	}

	return da.readFeaturesFromFile(filepath.Join(path, files[0].Name()))
}

func (da *DatasetAnalyzer) readFeaturesFromFile(filePath string) map[string][]string {
	featureInfo := make(map[string][]string)

	rows, err := utils.ReadTabularFile(filePath)
	if err != nil {
		return nil
	}

	layouts := []string{
		"2006-01-02",
		"2006/01/02",
		"02-01-2006",
		time.RFC3339,
		"2006-01-02 15:04:05",
	}

	for index, col := range rows[0] {
		value := rows[1][index]

		_, err := strconv.ParseFloat(value, 64)
		if err == nil {
			featureInfo["numerical"] = append(featureInfo["numerical"], col)
			continue
		}

		isNumerical := false
		for _, layout := range layouts {
			_, err := time.Parse(layout, value)
			if err == nil {
				featureInfo["numerical"] = append(featureInfo["numerical"], col)
				isNumerical = true
				break
			}
		}
		if isNumerical {
			continue
		}

		featureInfo["categorical"] = append(featureInfo["categorical"], col)
	}

	featureInfo["numerical_count"] = []string{strconv.Itoa(len(featureInfo["numerical"]))}
	featureInfo["categorical_count"] = []string{strconv.Itoa(len(featureInfo["categorical"]))}
	featureInfo["total"] = []string{strconv.Itoa(len(featureInfo["categorical"]) + len(featureInfo["numerical"]))}

	return featureInfo
}

// func (da *DatasetAnalyzer) analyzeCategoricalFeature(path string) *repo.CategoricalFeatureStatics {
// 	if ok := da.readCategoricalFeature(path); !ok {
// 		return nil
// 	}

// 	da.countCategoricalFeature()
// 	da.generateFeatureDetail()

// 	return &repo.CategoricalFeatureStatics{
// 		FeatureStat:   da.categoricalFeatureStat,   // categorical table
// 		FeatureDetail: da.categoricalFeatureDetail, // as-is categorical detail table
// 	}
// }

// func (da *DatasetAnalyzer) readCategoricalFeature(path string) bool {
// 	da.detailCount = make(map[string]map[string]map[string]int)
// 	// da.categoricalFeatureDetailType = make(map[string][]string)

// 	tvtdirs, err := utils.ReadDirs(path)
// 	if err != nil {
// 		return false
// 	}

// 	for _, tvt := range tvtdirs {
// 		files, _ := utils.ReadFiles(filepath.Join(path, tvt.Name()), TABULAR_EXTENSIONS, nil)
// 		if len(files) < 1 {
// 			return false
// 		}

// 		for _, file := range files {
// 			data_path := filepath.Join(path, tvt.Name(), file.Name())
// 			da.readCategoricalFeaturePerFile(data_path, tvt.Name())
// 		}
// 	}

// 	return true
// }

// func (da *DatasetAnalyzer) readCategoricalFeaturePerFile(filePath string, dirName string) {
// 	rows, err := utils.ReadTabularFile(filePath)
// 	if err != nil {
// 		return
// 	}

// 	for index, col := range rows[0] {
// 		if !slices.Contains(da.feature["categorical"], col) {
// 			continue
// 		}

// 		// count categorical detail count
// 		if _, exists := da.detailCount[col]; !exists {
// 			da.detailCount[col] = make(map[string]map[string]int)
// 		}

// 		for _, row := range rows[1:] {
// 			if _, exists := da.detailCount[col][row[index]]; !exists {
// 				da.detailCount[col][row[index]] = make(map[string]int)
// 				da.detailCount[col][row[index]]["train"] = 0
// 				da.detailCount[col][row[index]]["valid"] = 0
// 				da.detailCount[col][row[index]]["test"] = 0
// 			}
// 			da.detailCount[col][row[index]][dirName]++
// 			da.detailCount[col][row[index]]["total"]++
// 		}

// 		// // count categorical feature detail type
// 		// uniqueValues := make(map[string]struct{})
// 		// for _, row := range rows {
// 		// 	value := row[index]
// 		// 	if _, err := strconv.ParseFloat(value, 64); err != nil {
// 		// 		uniqueValues[value] = struct{}{}
// 		// 	}
// 		// }
// 		// existingValues := make(map[string]struct{})
// 		// for _, val := range da.categoricalFeatureDetailType[col] {
// 		// 	existingValues[val] = struct{}{}
// 		// }

// 		// for val := range uniqueValues {
// 		// 	if _, exists := existingValues[val]; !exists {
// 		// 		da.categoricalFeatureDetailType[col] = append(da.categoricalFeatureDetailType[col], val)
// 		// 	}
// 		// }
// 	}
// }

// func (da *DatasetAnalyzer) countCategoricalFeature() {
// 	da.categoricalFeatureStat = []map[string]string{}
// 	for _, category := range da.feature["categorical"] {
// 		row := make(map[string]string)
// 		row["feature"] = category

// 		for _, tvt := range TVT_NAMES {
// 			min := math.MaxInt32
// 			max := 0
// 			cate_min := []string{}
// 			cate_max := []string{}
// 			size := 0
// 			for key, value := range da.detailCount[category] {
// 				if value[tvt] < min {
// 					min = value[tvt]
// 					cate_min = []string{key}
// 				} else if value[tvt] == min {
// 					cate_min = append(cate_min, key)
// 				}

// 				if value[tvt] > max {
// 					max = value[tvt]
// 					cate_max = []string{key}
// 				} else if value[tvt] == max {
// 					cate_max = append(cate_max, key)
// 				}

// 				size += value[tvt]
// 			}
// 			row[fmt.Sprintf(`%s_size`, tvt)] = strconv.Itoa(size)
// 			row[fmt.Sprintf(`%s_minMode_size`, tvt)] = strings.Join(cate_min, ",")
// 			row[fmt.Sprintf(`%s_maxMode_size`, tvt)] = strings.Join(cate_max, ",")
// 		}
// 		trainSize, _ := strconv.Atoi(row["train_size"])
// 		validSize, _ := strconv.Atoi(row["valid_size"])
// 		testSize, _ := strconv.Atoi(row["test_size"])

// 		row["total"] = strconv.Itoa(trainSize + validSize + testSize)

// 		da.categoricalFeatureStat = append(da.categoricalFeatureStat, row)
// 	}
// }

// func (da *DatasetAnalyzer) generateFeatureDetail() {
// 	da.categoricalFeatureDetail = map[string][]map[string]string{}

// 	for category, item := range da.detailCount {
// 		for key, value := range item {
// 			detail := make(map[string]string)
// 			detail["feature"] = key
// 			for k, v := range value {
// 				detail[k] = strconv.Itoa(v)
// 			}
// 			da.categoricalFeatureDetail[category] = append(da.categoricalFeatureDetail[category], detail)
// 		}
// 	}
// }

// func (da *DatasetAnalyzer) analyzeNumericalFeature(path string) *repo.NumericalFeatureStatics {
// 	count := make(map[string]map[string]int)
// 	values := make(map[string]map[string][]float64) // boxplot

// 	tvtdirs, err := utils.ReadDirs(path)
// 	if err != nil {
// 		return nil
// 	}

// 	for _, tvt := range tvtdirs {
// 		files, _ := utils.ReadFiles(filepath.Join(path, tvt.Name()), TABULAR_EXTENSIONS, nil)
// 		if len(files) < 1 {
// 			continue
// 		}

// 		for _, file := range files {
// 			data_path := filepath.Join(path, tvt.Name(), file.Name())
// 			da.analyzeNumericalFeaturePerFile(data_path, tvt.Name(), count, values)
// 		}
// 	}

// 	boxPlot := da.countBolxplot(values, "numerical")

// 	return &repo.NumericalFeatureStatics{
// 		FeatureStat:   da.countNumericalFeature(count),
// 		FeatureDetail: da.countNumericalFeatureDetail(boxPlot),
// 		BoxPlot:       boxPlot,
// 		PDF:           da.countPDF(boxPlot, "numerical"),
// 	}
// }

// func (da *DatasetAnalyzer) analyzeNumericalFeaturePerFile(data_path string, dirName string, count map[string]map[string]int, values map[string]map[string][]float64) {
// 	rows, err := utils.ReadTabularFile(data_path)
// 	if err != nil {
// 		return
// 	}

// 	for index, col := range rows[0] {
// 		if !slices.Contains(da.feature["numerical"], col) {
// 			continue
// 		}

// 		// count features
// 		value := []float64{}
// 		for _, row := range rows[1:] {
// 			if val, err := strconv.ParseFloat(row[index], 64); err == nil {
// 				if _, exists := count[col]; !exists {
// 					count[col] = make(map[string]int)
// 				}
// 				if _, exists := count[col][dirName]; !exists {
// 					count[col][dirName] = 0
// 				}
// 				count[col][dirName]++
// 				value = append(value, val)
// 			}
// 		}

// 		if _, exists := values[col]; !exists {
// 			values[col] = make(map[string][]float64)
// 		}
// 		if _, exists := values[col][dirName]; !exists {
// 			values[col][dirName] = make([]float64, 0)
// 		}
// 		values[col][dirName] = value
// 	}
// }

// func (da *DatasetAnalyzer) countNumericalFeature(count map[string]map[string]int) []map[string]string {
// 	featureStat := []map[string]string{}
// 	for k, v := range count {
// 		stat := make(map[string]string)
// 		stat["feature"] = k

// 		total := 0
// 		for k2, v2 := range v {
// 			stat[k2] = strconv.Itoa(v2)
// 			total += v2
// 		}
// 		stat["total"] = strconv.Itoa(total)

// 		featureStat = append(featureStat, stat)
// 	}
// 	return featureStat
// }

// func (da *DatasetAnalyzer) countNumericalFeatureDetail(values map[string]map[string][]float64) map[string][]map[string]string {
// 	featureDetail := make(map[string][]map[string]string)

// 	for _, feature := range da.feature["numerical"] {
// 		categoryStats := da.calcNumericalStats(values[feature])
// 		for _, category := range STAT_CATEGORY_NAMES {
// 			detail := make(map[string]string)
// 			detail["category"] = category
// 			detail["train"] = categoryStats[category]["train"]
// 			detail["valid"] = categoryStats[category]["valid"]
// 			detail["test"] = categoryStats[category]["test"]

// 			featureDetail[feature] = append(featureDetail[feature], detail)
// 		}
// 	}

// 	return featureDetail
// }

// func (da *DatasetAnalyzer) calcNumericalStats(values map[string][]float64) map[string]map[string]string {
// 	categoryStats := make(map[string]map[string]string)
// 	categoryStats["min"] = make(map[string]string)
// 	categoryStats["max"] = make(map[string]string)
// 	categoryStats["mean"] = make(map[string]string)
// 	categoryStats["median"] = make(map[string]string)
// 	categoryStats["stdev"] = make(map[string]string)

// 	for _, tvt := range TVT_NAMES {
// 		min, _ := stats.Min(values[tvt])
// 		max, _ := stats.Max(values[tvt])
// 		mean, _ := stats.Mean(values[tvt])
// 		median, _ := stats.Median(values[tvt])
// 		stdev, _ := stats.StandardDeviation(values[tvt])

// 		categoryStats["min"][tvt] = strconv.FormatFloat(math.Round(min*100)/100, 'f', -1, 64)
// 		categoryStats["max"][tvt] = strconv.FormatFloat(math.Round(max*100)/100, 'f', -1, 64)
// 		categoryStats["mean"][tvt] = strconv.FormatFloat(math.Round(mean*100)/100, 'f', -1, 64)
// 		categoryStats["median"][tvt] = strconv.FormatFloat(math.Round(median*100)/100, 'f', -1, 64)
// 		categoryStats["stdev"][tvt] = strconv.FormatFloat(math.Round(stdev*100)/100, 'f', -1, 64)
// 	}

// 	return categoryStats
// }

// func (da *DatasetAnalyzer) countBolxplot(values map[string]map[string][]float64, targetFeature repo.FeatureType) map[string]map[string][]float64 {
// 	bolxPlot := make(map[string]map[string][]float64)

// 	for _, tvt := range TVT_NAMES {
// 		for _, feature := range da.feature[string(targetFeature)] {
// 			if _, exists := bolxPlot[feature]; !exists {
// 				bolxPlot[feature] = make(map[string][]float64)
// 			}
// 			bolxPlot[feature][tvt] = stats.LoadRawData(values[feature][tvt])
// 		}
// 	}

// 	return bolxPlot
// }

// func (da *DatasetAnalyzer) countPDF(values map[string]map[string][]float64, targetFeature repo.FeatureType) map[string]map[string]repo.DataPDF {
// 	pdf := make(map[string]map[string]repo.DataPDF) // gaussian

// 	for _, tvt := range TVT_NAMES {
// 		for _, feature := range da.feature[string(targetFeature)] {
// 			if _, exists := pdf[feature]; !exists {
// 				pdf[feature] = make(map[string]repo.DataPDF)
// 			}
// 			pdf[feature][tvt] = da.ComputePDF(values[feature][tvt])
// 		}
// 	}

// 	return pdf
// }

func (da *DatasetAnalyzer) ComputePDF(d []float64) repo.DataPDF {
	sample := go_stats.Sample{
		Xs:      d,
		Weights: nil,
		Sorted:  false,
	}

	kde := go_stats.KDE{
		Sample: sample,
	}

	points := make([]repo.DataPoint, len(d))
	for i := range points {
		x := d[i]
		points[i] = repo.DataPoint{
			X: x,
			Y: kde.PDF(x),
		}
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].X < points[j].X
	})

	if len(points) > 0 {
		newPoints := points[:1]
		for i := 1; i < len(points); i++ {
			if points[i].X != points[i-1].X {
				newPoints = append(newPoints, points[i])
			}
		}
		points = newPoints
	}

	newLength := len(points)
	xData := make([]float64, newLength)
	yData := make([]float64, newLength)

	for i, point := range points {
		xData[i] = point.X
		yData[i] = point.Y
	}

	return repo.DataPDF{
		XData: xData,
		YData: yData,
	}
}

func (da *DatasetAnalyzer) CompareNumericalFeature(datasetId int, feature1 string, feature2 string) (*repo.CompareNumericalFeaturesStatics, *logger.Report) {
	features := make(map[string]map[string][]float64)
	pairWise := make(map[string][]stats.Coordinate)
	regression := make(map[string]stats.Series)

	path, r := da.datasetDAO.SelectDataPathByDataSetId(da.ctx, datasetId)
	if r != nil {
		return nil, r
	}

	tvtdirs, err := utils.ReadDirs(path)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	for _, tvt := range tvtdirs {
		files, _ := utils.ReadFiles(filepath.Join(path, tvt.Name()), TABULAR_EXTENSIONS, nil)
		if len(files) < 1 {
			continue
		}

		for _, file := range files {
			data_path := filepath.Join(path, tvt.Name(), file.Name())
			rows, err := utils.ReadTabularFile(data_path)
			if err != nil {
				continue
			}

			for index, col := range rows[0] {
				if _, exists := features[col]; !exists {
					features[col] = make(map[string][]float64)
				}
				if _, exists := features[col][tvt.Name()]; !exists {
					features[col][tvt.Name()] = make([]float64, 0)
				}

				if col == feature1 || col == feature2 {
					features[col][tvt.Name()] = []float64{}
					for _, row := range rows[1:] {
						value, err := strconv.ParseFloat(row[index], 64)
						if err == nil {
							features[col][tvt.Name()] = append(features[col][tvt.Name()], value)
						}
					}
				}
			}

			for index, value := range features[feature1][tvt.Name()] {

				pairWise[tvt.Name()] = append(pairWise[tvt.Name()], stats.Coordinate{
					X: value,
					Y: features[feature2][tvt.Name()][index],
				})
			}
			regression[tvt.Name()], _ = stats.LinearRegression(pairWise[tvt.Name()])
		}
	}

	return &repo.CompareNumericalFeaturesStatics{
		CompareResult: pairWise,
		Regression:    regression,
	}, nil
}

// func (da *DatasetAnalyzer) computeNumericalHeatmap(path string) *repo.NumericalHeatmap {
// 	if len(da.feature["numerical"]) < 1 {
// 		return nil
// 	}

// 	args := []string{
// 		"--data_path", path,
// 		//"--selected_cat_features",
// 	}
// 	//args = append(args, da.feature["categorical"]...)

// 	cmd := exec.Command("./scripts/numnum_corr", args...)

// 	var output bytes.Buffer
// 	var stderr bytes.Buffer
// 	cmd.Stdout = &output
// 	cmd.Stderr = &stderr

// 	err := cmd.Run()
// 	if err != nil {
// 		logger.Debug(fmt.Sprint(err) + ": " + stderr.String())
// 		return nil
// 	}

// 	out := string(output.String())
// 	out = strings.ReplaceAll(out, "'", "\"")
// 	out = strings.ReplaceAll(out, "nan", "0")

// 	result := repo.NumericalHeatmap{}
// 	if err := json.Unmarshal([]byte(out), &result); err != nil {
// 		logger.Debug(err)
// 		return nil
// 	}

// 	return &result
// }

// func (da *DatasetAnalyzer) computeCategoricalHeatmap(path string) *repo.CategoricalHeatmap {
// 	if len(da.feature["categorical"]) < 1 {
// 		return nil
// 	}

// 	args := []string{
// 		"--data_path", path,
// 		"--selected_cat_features",
// 	}
// 	args = append(args, da.feature["categorical"]...)

// 	cmd := exec.Command("./scripts/catcat_corr", args...)

// 	var output bytes.Buffer
// 	var stderr bytes.Buffer
// 	cmd.Stdout = &output
// 	cmd.Stderr = &stderr

// 	err := cmd.Run()
// 	if err != nil {
// 		logger.Debug(fmt.Sprint(err) + ": " + stderr.String())
// 		return nil
// 	}

// 	out := string(output.String())
// 	out = strings.ReplaceAll(out, "'", "\"")
// 	out = strings.ReplaceAll(out, "nan", "0")

// 	result := repo.CategoricalHeatmap{}
// 	if err := json.Unmarshal([]byte(out), &result); err != nil {
// 		logger.Debug(err)
// 		return nil
// 	}

//		return &result
//	}
func (da *DatasetAnalyzer) computeNumericalHeatmap(path string) *repo.NumericalHeatmap {
	features := make(map[string][]float64)
	correlation := make(map[string]map[string]map[string]float64)

	tvtdirs, err := utils.ReadDirs(path)
	if err != nil {
		return nil
	}

	for _, tvt := range tvtdirs {
		files, _ := utils.ReadFiles(filepath.Join(path, tvt.Name()), TABULAR_EXTENSIONS, nil)
		if len(files) < 1 {
			return nil
		}

		for _, file := range files {
			data_path := filepath.Join(path, tvt.Name(), file.Name())
			rows, err := utils.ReadTabularFile(data_path)
			if err != nil {
				continue
			}

			for index, col := range rows[0] {
				features[col] = []float64{}
				for _, row := range rows[1:] {
					value, err := strconv.ParseFloat(row[index], 64)
					if err == nil {
						features[col] = append(features[col], value)
					}
				}
			}

			for feature1, feature1Values := range features {
				if _, exists := correlation[tvt.Name()]; !exists {
					correlation[tvt.Name()] = make(map[string]map[string]float64)
				}
				if _, exists := correlation[tvt.Name()][feature1]; !exists {
					correlation[tvt.Name()][feature1] = make(map[string]float64, 0)
				}
				correlation[tvt.Name()][feature1] = make(map[string]float64)
				for feature2, feature2Values := range features {
					if feature1 > feature2 {
						continue
					}

					if len(feature1Values) != len(feature2Values) {
						continue
					}

					correlationValue, err := stats.Correlation(feature1Values, feature2Values)
					if err != nil {
						continue
					}
					roundedCorrelationValue, err := stats.Round(correlationValue, 2)
					if err != nil {
						continue
					}
					correlation[tvt.Name()][feature1][feature2] = roundedCorrelationValue
				}
			}
		}
	}

	return &repo.NumericalHeatmap{
		Feature:     da.feature["numerical"],
		Correlation: correlation,
	}
}

func (da *DatasetAnalyzer) computeCategoricalHeatmap(path string) *repo.CategoricalHeatmap {
	if len(da.feature["categorical"]) < 1 {
		return nil
	}

	args := []string{
		"--data_path", path,
		"--selected_cat_features",
	}
	args = append(args, da.feature["categorical"]...)

	cmd := exec.Command("./scripts/catcat_corr", args...)

	var output bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logger.Debug(cmd, da.feature["categorical"])
		logger.Debug(fmt.Sprint(err) + ": " + stderr.String())
		return nil
	}

	out := string(output.String())
	out = strings.ReplaceAll(out, "'", "\"")
	out = strings.ReplaceAll(out, "nan", "0")

	result := repo.CategoricalHeatmap{}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		logger.Debug(err)
		return nil
	}

	return &result
}

func (da *DatasetAnalyzer) CompareCategoricalFeature(datasetId int, feature1 string, feature2 string) (*repo.CompareCategoricalFeaturesStatics, *logger.Report) {
	features := []string{}
	cateFeatures := []string{}
	cateDetailFeatures := make(map[string][]string)
	compareResult := make(map[string]map[string]map[string]int)

	feature1Index := 0
	feature2Index := 0

	path, err := da.datasetDAO.SelectDataPathByDataSetId(da.ctx, datasetId)
	if err != nil {
		return nil, err
	}

	tvtdirs, _ := utils.ReadDirs(path)
	childrenDatasets, _ := da.datasetDAO.SelectDataSetByParentID(da.ctx, datasetId)
	files, _ := utils.ReadFiles(childrenDatasets[0].Path, nil, nil)
	_path := filepath.Join(childrenDatasets[0].Path, files[0].Name())
	rows, _ := utils.ReadTabularFile(_path)
	// categorical feature
	for index, col := range rows[0] {
		features = append(features, col)

		value := rows[1][index]
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			cateFeatures = append(cateFeatures, col)
		}
	}
	// cateDetailFeature
	for _, tvt := range tvtdirs {
		files, err := utils.ReadFiles(filepath.Join(path, tvt.Name()), TABULAR_EXTENSIONS, nil)
		if len(files) < 1 {
			return nil, logger.CreateReport(&logger.CODE_FILE_NOT_EXIST, err)
		}

		for _, file := range files {
			data_path := filepath.Join(path, tvt.Name(), file.Name())
			rows, err := utils.ReadTabularFile(data_path)
			if err != nil {
				return nil, logger.CreateReport(&logger.CODE_FAILE, err)
			}

			for index, col := range rows[0] {
				uniqueValues := make(map[string]struct{})
				for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
					value := rows[rowIndex][index]
					if _, err := strconv.ParseFloat(value, 64); err != nil {
						uniqueValues[value] = struct{}{}
					}
				}

				existingValues := make(map[string]struct{})
				for _, val := range cateDetailFeatures[col] {
					existingValues[val] = struct{}{}
				}

				for val := range uniqueValues {
					if _, exists := existingValues[val]; !exists {
						cateDetailFeatures[col] = append(cateDetailFeatures[col], val)
					}
				}

			}
		}
	}

	for _, tvt := range tvtdirs {
		files, err := utils.ReadFiles(filepath.Join(path, tvt.Name()), TABULAR_EXTENSIONS, nil)
		if len(files) < 1 {
			return nil, logger.CreateReport(&logger.CODE_FILE_NOT_EXIST, err)
		}

		for _, file := range files {
			data_path := filepath.Join(path, tvt.Name(), file.Name())
			rows, err := utils.ReadTabularFile(data_path)
			if err != nil {
				return nil, logger.CreateReport(&logger.CODE_FAILE, err)
			}

			for index, feature := range features {
				if feature == feature1 {
					feature1Index = index
				}
				if feature == feature2 {
					feature2Index = index
				}
			}

			for _, value1 := range cateDetailFeatures[feature1] {
				for _, row := range rows {
					if _, exists := compareResult[tvt.Name()]; !exists {
						compareResult[tvt.Name()] = make(map[string]map[string]int)
					}
					if _, exists := compareResult[tvt.Name()][value1]; !exists {
						compareResult[tvt.Name()][value1] = make(map[string]int)
					}

					for _, value2 := range cateDetailFeatures[feature2] {
						if _, exists := compareResult[tvt.Name()][value1][value2]; !exists {
							compareResult[tvt.Name()][value1][value2] = 0
						}
						if row[feature1Index] == value1 && row[feature2Index] == value2 {
							compareResult[tvt.Name()][value1][value2] += 1
						}

					}
				}

				sum := 0
				for _, count := range compareResult[tvt.Name()][value1] {
					sum += count
				}
				compareResult[tvt.Name()][value1]["sum"] = sum
			}
		}
	}

	return &repo.CompareCategoricalFeaturesStatics{
		CategoricalFeatures:       cateFeatures,
		CategoricalDetailFeatures: cateDetailFeatures,
		CompareResult:             compareResult,
	}, nil
}

func (da *DatasetAnalyzer) CompareCategoricalNumericalFeature(datasetId int, cate_feature string, nume_feature string) (*repo.CompareCategoricalNumericalFeaturesStatics, *logger.Report) {
	features := []string{}
	cateFeatures := []string{}
	cateDetailFeatures := make(map[string][]string)

	numeIndex := -1
	cateIndex := -1

	// categorical features
	childrenDatasets, _ := da.datasetDAO.SelectDataSetByParentID(da.ctx, datasetId)
	files, _ := utils.ReadFiles(childrenDatasets[0].Path, nil, nil)
	_path := filepath.Join(childrenDatasets[0].Path, files[0].Name())
	rows, _ := utils.ReadTabularFile(_path)

	for index, col := range rows[0] {
		features = append(features, col)

		value := rows[1][index]
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			cateFeatures = append(cateFeatures, col)
		}
	}

	// cateDetailFeatures
	path, err := da.datasetDAO.SelectDataPathByDataSetId(da.ctx, datasetId)
	if err != nil {
		return nil, err
	}

	tvtdirs, _ := utils.ReadDirs(path)
	for _, tvt := range tvtdirs {
		files, err := utils.ReadFiles(filepath.Join(path, tvt.Name()), TABULAR_EXTENSIONS, nil)
		if len(files) < 1 {
			return nil, logger.CreateReport(&logger.CODE_FILE_NOT_EXIST, err)
		}

		for _, file := range files {
			data_path := filepath.Join(path, tvt.Name(), file.Name())
			rows, err := utils.ReadTabularFile(data_path)
			if err != nil {
				return nil, logger.CreateReport(&logger.CODE_FAILE, err)
			}

			for index, col := range rows[0] {
				uniqueValues := make(map[string]struct{})
				for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
					_value := rows[rowIndex][index]
					if _, err := strconv.ParseFloat(_value, 64); err != nil {
						uniqueValues[_value] = struct{}{}
					}
				}

				existingValues := make(map[string]struct{})
				for _, val := range cateDetailFeatures[col] {
					existingValues[val] = struct{}{}
				}

				for val := range uniqueValues {
					if _, exists := existingValues[val]; !exists {
						cateDetailFeatures[col] = append(cateDetailFeatures[col], val)
					}
				}

			}
		}
	}

	// values
	values := make(map[string]map[string][]float64)

	path, r := da.datasetDAO.SelectDataPathByDataSetId(da.ctx, datasetId)
	if r != nil {
		return nil, r
	}

	for _, tvt := range tvtdirs {
		files, _ := utils.ReadFiles(filepath.Join(path, tvt.Name()), TABULAR_EXTENSIONS, nil)
		if len(files) < 1 {
			continue
		}

		for _, file := range files {
			data_path := filepath.Join(path, tvt.Name(), file.Name())
			rows, err := utils.ReadTabularFile(data_path)
			if err != nil {
				continue
			}

			for index, feature := range features {
				if feature == nume_feature {
					numeIndex = index
				}
				if feature == cate_feature {
					cateIndex = index
				}
			}
			if numeIndex != -1 && cateIndex != -1 {

				for _, row := range rows[1:] {
					numeValue := row[numeIndex]
					cateValue := row[cateIndex]

					if _, exists := values[cateValue]; !exists {
						values[cateValue] = make(map[string][]float64)
					}

					if _, exists := values[cateValue][tvt.Name()]; !exists {
						values[cateValue][tvt.Name()] = make([]float64, 0)
					}

					if val, err := strconv.ParseFloat(numeValue, 64); err == nil {
						values[cateValue][tvt.Name()] = append(values[cateValue][tvt.Name()], val)
					}
				}

			}
		}
	}
	// boxplot
	boxPlot := make(map[string]map[string][]float64)

	for _, tvt := range TVT_NAMES {
		for _, feature := range cateDetailFeatures[cate_feature] {
			if _, exists := boxPlot[feature]; !exists {
				boxPlot[feature] = make(map[string][]float64)
			}
			boxPlot[feature][tvt] = stats.LoadRawData(values[feature][tvt])
		}
	}
	// pdf = gaussian
	pdf := make(map[string]map[string]repo.DataPDF)

	for _, tvt := range TVT_NAMES {
		for _, feature := range cateDetailFeatures[cate_feature] {
			if _, exists := pdf[feature]; !exists {
				pdf[feature] = make(map[string]repo.DataPDF)
			}
			pdf[feature][tvt] = da.ComputePDF(values[feature][tvt])
		}
	}

	return &repo.CompareCategoricalNumericalFeaturesStatics{
		Feature:                   cateFeatures,
		CategoricalDetailFeatures: cateDetailFeatures,
		PDF:                       pdf,
	}, nil
}

func (da *DatasetAnalyzer) ComputeCategoricalNumericalHeatmap(path string) *repo.CategoricalNumericalHeatmap {
	if len(da.feature["categorical"]) < 1 || len(da.feature["numerical"]) < 1 {
		return nil
	}

	args := []string{
		"--data_path", path,
		"--selected_cat_features",
	}
	args = append(args, da.feature["categorical"]...)
	args = append(args, "--selected_nume_features")
	args = append(args, da.feature["numerical"]...)

	cmd := exec.Command("./scripts/catnum_corr", args...)

	var output bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logger.Debug(fmt.Sprint(err) + ": " + stderr.String())
		return nil
	}

	out := output.String()
	// UTF-8로 디코딩
	// decoder := korean.EUCKR.NewDecoder()
	// out, _, err := transform.String(decoder, output.String())
	// if err != nil {
	// 	logger.Debug("Failed to transform encoding:", err)
	// 	return nil
	// }
	out = strings.ReplaceAll(out, "'", "\"")
	out = strings.ReplaceAll(out, "nan", "0")

	result := repo.CategoricalNumericalHeatmap{}
	if err = json.Unmarshal([]byte(out), &result); err != nil {
		logger.Debug(err)
		return nil
	}

	return &result
}

func (da *DatasetAnalyzer) countNonetypeDataset(path string, dataType string) *repo.NoneTypeStat {
	if dataType == utils.DATA_TYPE_IMG {
		imageStat := make(map[string]*repo.ImageTypeStat)
		da.countNonetypeImage(path, imageStat)
		return &repo.NoneTypeStat{
			ImageStat: imageStat,
		}
	} else if dataType == utils.DATA_TYPE_TABLE {
		tabularStat := make(map[string]map[string]*repo.TabularStat)
		da.countNonetypeTabular(path, tabularStat)
		return &repo.NoneTypeStat{
			TabularStat: tabularStat,
		}
	}

	return nil
}

func (da *DatasetAnalyzer) countNonetypeImage(path string, statMap map[string]*repo.ImageTypeStat) {
	if children, err := utils.ReadDirs(path); err == nil {
		for _, child := range children {
			da.countNonetypeImage(filepath.Join(path, child.Name()), statMap)
		}
	}

	except := []string{}
	copy(TABULAR_EXTENSIONS, except)
	except = append(except, ".txt", ".json")
	if files, err := utils.ReadFiles(path, nil, except); err == nil {
		dirName := path[strings.LastIndex(path, string(filepath.Separator))+1:]

		resolution := make(map[string]map[int]map[int]int)
		resolution[dirName] = make(map[int]map[int]int)

		count := make(map[string]int)
		count[dirName] = 0

		da.singlelabelResolutionPerDir(path, dirName, resolution, count)

		statMap[path] = &repo.ImageTypeStat{
			Count:             len(files),
			ResolutionStatics: &repo.ResolutionStatics{Resolution: resolution, Count: count},
		}
	}
}

func (da *DatasetAnalyzer) countNonetypeTabular(path string, statMap map[string]map[string]*repo.TabularStat) {
	if children, err := utils.ReadDirs(path); err == nil {
		for _, child := range children {
			da.countNonetypeTabular(filepath.Join(path, child.Name()), statMap)
		}
	}

	if files, err := utils.ReadFiles(path, TABULAR_EXTENSIONS, []string{".txt", ".json"}); err == nil {
		if _, exists := statMap[path]; !exists {
			statMap[path] = make(map[string]*repo.TabularStat)
		}

		for _, file := range files {
			filePath := filepath.Join(path, file.Name())

			statMap[path][file.Name()] = &repo.TabularStat{
				Count:    da.countRows(filePath),
				Features: da.readFeaturesFromFile(filePath),
			}
		}
	}
}

func (d *DatasetAnalyzer) countRows(filePath string) int {
	if rows, err := utils.ReadTabularFile(filePath); err != nil {
		return 0
	} else {
		return len(rows)
	}
}

func (da *DatasetAnalyzer) computeCategoricalHeatmapOnEngine(path string, datasetId int) *repo.CategoricalHeatmap {
	catFeatures := da.feature["categorical"]
	if len(catFeatures) < 1 {
		return nil
	}

	reqClient := NewDatasetRequestClient("/api/dataset/categorical_heatmap", datasetId)

	type CatHeatReqBody struct {
		StaticPath          string   `json:"static_path"`
		DataPath            string   `json:"dataset_path"`
		CategoricalFeatures []string `json:"cat_features"`
	}
	reqData := CatHeatReqBody{
		StaticPath:          reqClient.StaticPath,
		DataPath:            path,
		CategoricalFeatures: catFeatures,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		logger.Debug(err)
		return nil
	}

	//req to python engine
	_, err = reqClient.SendPostRequest(jsonData)
	if err != nil {
		logger.Debug(err)
		return nil
	}

	//return empty value
	result := repo.CategoricalHeatmap{}
	return &result
}

func (da *DatasetAnalyzer) computeCategoricalNumericalHeatmapOnEngine(path string, datasetId int) *repo.CategoricalNumericalHeatmap {
	catFeatures := da.feature["categorical"]
	numFeatures := da.feature["numerical"]
	if len(catFeatures) < 1 || len(numFeatures) < 1 {
		return nil
	}

	reqClient := NewDatasetRequestClient("/api/dataset/categorical_numerical_heatmap", datasetId)

	type CatNumHeatReqBody struct {
		StaticPath          string   `json:"static_path"`
		DataPath            string   `json:"dataset_path"`
		CategoricalFeatures []string `json:"cat_features"`
		NumericalFeatures   []string `json:"num_features"`
	}
	reqData := CatNumHeatReqBody{
		StaticPath:          reqClient.StaticPath,
		DataPath:            path,
		CategoricalFeatures: catFeatures,
		NumericalFeatures:   numFeatures,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		logger.Debug(err)
		return nil
	}

	//req to python engine
	_, err = reqClient.SendPostRequest(jsonData)
	if err != nil {
		logger.Debug(err)
		return nil
	}

	//return empty
	result := repo.CategoricalNumericalHeatmap{}
	return &result
}

func (da *DatasetAnalyzer) analyzeCategoricalFeatureOnEngine(path string, datasetId int) *repo.CategoricalFeatureStatics {
	catFeatures := da.feature["categorical"]
	if len(catFeatures) < 1 {
		return nil
	}

	reqClient := NewDatasetRequestClient("/api/dataset/categorical_feature", datasetId)

	type CatFeatReqBody struct {
		StaticPath          string   `json:"static_path"`
		DataPath            string   `json:"dataset_path"`
		CategoricalFeatures []string `json:"cat_features"`
	}
	reqData := CatFeatReqBody{
		StaticPath:          reqClient.StaticPath,
		DataPath:            path,
		CategoricalFeatures: catFeatures,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		logger.Debug(err)
		return nil
	}

	//req to python engine
	//ignore response value(file_path)
	_, err = reqClient.SendPostRequest(jsonData)
	if err != nil {
		logger.Debug(err)
		return nil
	}

	//return empty value
	result := repo.CategoricalFeatureStatics{}
	return &result
}

func (da *DatasetAnalyzer) analyzeNumericalFeatureOnEngine(path string, datasetId int) *repo.NumericalFeatureStatics {
	numFeatures := da.feature["numerical"]
	if len(numFeatures) < 1 {
		return nil
	}

	reqClient := NewDatasetRequestClient("/api/dataset/numerical_feature", datasetId)

	type NumFeatReqBody struct {
		StaticPath        string   `json:"static_path"`
		DataPath          string   `json:"dataset_path"`
		NumericalFeatures []string `json:"num_features"`
	}
	reqData := NumFeatReqBody{
		StaticPath:        reqClient.StaticPath,
		DataPath:          path,
		NumericalFeatures: numFeatures,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		logger.Debug(err)
		return nil
	}
	_, err = reqClient.SendPostRequest(jsonData)
	if err != nil {
		logger.Debug(err)
		return nil
	}

	result := repo.NumericalFeatureStatics{}
	return &result
}

func (da *DatasetAnalyzer) computeNumericalHeatmapOnEngine(path string, datasetId int) *repo.NumericalHeatmap {
	numFeatures := da.feature["numerical"]
	if len(numFeatures) < 1 {
		return nil
	}

	reqClient := NewDatasetRequestClient("/api/dataset/numerical_heatmap", datasetId)

	type NumHeatReqBody struct {
		StaticPath        string   `json:"static_path"`
		DataPath          string   `json:"dataset_path"`
		NumericalFeatures []string `json:"num_features"`
	}
	reqData := NumHeatReqBody{
		StaticPath:        reqClient.StaticPath,
		DataPath:          path,
		NumericalFeatures: numFeatures,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		logger.Debug(err)
		return nil
	}

	_, err = reqClient.SendPostRequest(jsonData)
	if err != nil {
		logger.Debug(err)
		return nil
	}

	result := repo.NumericalHeatmap{}
	return &result
}

// TIFF 이미지의 크기 정보만 직접 추출
func (da *DatasetAnalyzer) getTiffDimensions(filePath string) (width int, height int, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	// 바이너리 데이터에서 필요한 정보만 직접 추출
	// TIFF 헤더 건너뛰기 (8바이트)
	header := make([]byte, 8)
	if _, err := file.Read(header); err != nil {
		return 0, 0, err
	}

	// 리틀 엔디안인지 확인
	isLittleEndian := header[0] == 'I' && header[1] == 'I'

	// IFD 오프셋 읽기
	var ifdOffset uint32
	if isLittleEndian {
		ifdOffset = uint32(header[4]) | uint32(header[5])<<8 | uint32(header[6])<<16 | uint32(header[7])<<24
	} else {
		ifdOffset = uint32(header[7]) | uint32(header[6])<<8 | uint32(header[5])<<16 | uint32(header[4])<<24
	}

	// IFD 위치로 이동
	if _, err := file.Seek(int64(ifdOffset), 0); err != nil {
		return 0, 0, err
	}

	// IFD 엔트리 수 읽기
	entryCount := make([]byte, 2)
	if _, err := file.Read(entryCount); err != nil {
		return 0, 0, err
	}

	var numEntries uint16
	if isLittleEndian {
		numEntries = uint16(entryCount[0]) | uint16(entryCount[1])<<8
	} else {
		numEntries = uint16(entryCount[1]) | uint16(entryCount[0])<<8
	}

	// 각 IFD 엔트리 읽기 (각 엔트리는 12바이트)
	foundWidth := false
	foundHeight := false
	var imgWidth, imgHeight uint32

	for i := 0; i < int(numEntries) && (!foundWidth || !foundHeight); i++ {
		entry := make([]byte, 12)
		if _, err := file.Read(entry); err != nil {
			return 0, 0, err
		}

		var tagID uint16
		if isLittleEndian {
			tagID = uint16(entry[0]) | uint16(entry[1])<<8
		} else {
			tagID = uint16(entry[1]) | uint16(entry[0])<<8
		}

		// ImageWidth는 태그 256 (0x0100)
		// ImageLength(높이)는 태그 257 (0x0101)
		if tagID == 256 { // ImageWidth
			if isLittleEndian {
				imgWidth = uint32(entry[8]) | uint32(entry[9])<<8 | uint32(entry[10])<<16 | uint32(entry[11])<<24
			} else {
				imgWidth = uint32(entry[11]) | uint32(entry[10])<<8 | uint32(entry[9])<<16 | uint32(entry[8])<<24
			}
			foundWidth = true
		} else if tagID == 257 { // ImageLength (Height)
			if isLittleEndian {
				imgHeight = uint32(entry[8]) | uint32(entry[9])<<8 | uint32(entry[10])<<16 | uint32(entry[11])<<24
			} else {
				imgHeight = uint32(entry[11]) | uint32(entry[10])<<8 | uint32(entry[9])<<16 | uint32(entry[8])<<24
			}
			foundHeight = true
		}
	}

	if !foundWidth || !foundHeight {
		return 0, 0, fmt.Errorf("width or height information not found in TIFF")
	}

	return int(imgWidth), int(imgHeight), nil
}

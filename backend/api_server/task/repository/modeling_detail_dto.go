package repository

import (
	"api_server/logger"
	"api_server/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"sort"

	config_service "api_server/configuration/service"
)

type ModelingDetailDTO struct {
	ID         int       `json:"id"`
	Model      string    `json:"model"`
	DataType   string    `json:"data_type"`
	Data       []string  `json:"data"`
	CreatedAt  time.Time `json:"created_at"`
	ModelingID int       `json:"modeling_id"`
}

type StatusJsonEntity struct {
	TrialID      int
	TrialNo      int
	ModelingStep string
	StatusJson   sql.NullString
}

type LossPerfChartResponse struct {
	LastTrial    int         `json:"last_trial"`
	ModelingStep string      `json:"modeling_step"`
	TargetMetric string      `json:"target_metric"`
	LossItems    []LossChart `json:"loss"`
	LossDims     []string    `json:"loss_dims"`
	PerfItems    []PerfChart `json:"perf"`
	PerfDims     []string    `json:"perf_dims"`
}

type LossChart struct {
	TrialNo   int       `json:"trial_no"`
	Epoch     string    `json:"epoch"`
	TrainLoss *float64  `json:"train_loss"`
	ValidLoss *float64  `json:"valid_loss"`
	LogTime   time.Time `json:"log_time,omitempty"`
}

type PerfChart struct {
	TrialNo    int       `json:"trial_no"`
	Epoch      string    `json:"epoch"`
	TrainScore *float64  `json:"train_score"`
	ValidScore *float64  `json:"valid_score"`
	LogTime    time.Time `json:"log_time,omitempty"`
}

type TabularChartResponse struct {
	ModelingStep string             `json:"modeling_step"`
	TargetMetric string             `json:"target_metric"`
	Dims         []string           `json:"dims"`
	Items        []TabularChartItem `json:"items"`
}

type TabularChartItem struct {
	TrialNo string  `json:"trial_no"`
	Score   float64 `json:"score,omitempty"`
}

func (lp *LossPerfChartResponse) AddChartData(entity StatusJsonEntity) int {
	if !entity.StatusJson.Valid {
		return -1
	}

	status := make(map[string]interface{})
	if err := json.Unmarshal([]byte(entity.StatusJson.String), &status); err != nil {
		return -1
	}

	_epoch, _ := status["epoch"].(float64)
	epoch := strconv.Itoa(entity.TrialNo) + "-" + strconv.Itoa(int(_epoch)) + "epoch"

	targetMetric, _ := status["target_metric"].(string)

	trainLoss, _ := status["train_loss"].(float64)
	validLoss, _ := status["validation_loss"].(float64)
	trainScore, _ := status["train_target_score"].(float64)
	validScore, _ := status["validation_target_score"].(float64)

	timeFloat, _ := status["log_time"].(float64)
	sec, dec := math.Modf(timeFloat)
	logTime := time.Unix(int64(sec), int64(dec*(1e9)))

	lp.LastTrial = entity.TrialID
	lp.TargetMetric = targetMetric
	lp.ModelingStep = entity.ModelingStep

	loss := LossChart{
		TrialNo:   entity.TrialNo,
		Epoch:     epoch,
		TrainLoss: &trainLoss,
		ValidLoss: &validLoss,
		LogTime:   logTime,
	}
	lp.LossItems = append(lp.LossItems, loss)

	perf := PerfChart{
		TrialNo:    entity.TrialNo,
		Epoch:      epoch,
		TrainScore: &trainScore,
		ValidScore: &validScore,
		LogTime:    logTime,
	}
	lp.PerfItems = append(lp.PerfItems, perf)

	return int(_epoch)
}

func (lp *LossPerfChartResponse) AddEmptyData(trial_no int, epoch int) {
	_epoch := strconv.Itoa(trial_no) + "-"

	loss := LossChart{TrialNo: trial_no, Epoch: _epoch}
	lp.LossItems = append(lp.LossItems, loss)

	perf := PerfChart{TrialNo: trial_no, Epoch: _epoch}
	lp.PerfItems = append(lp.PerfItems, perf)
}

type ModelPerformanceEntity struct {
	ModelingID int            `json:"modeling_id"`
	Model      string         `json:"model"`
	DBModel    string         `json:"db_model"`
	Score      sql.NullString `json:"score"`
	InfTime    sql.NullString `json:"inf_time"`
	CfMatrix   sql.NullString `json:"cf_matrix"`
	PredResult sql.NullString `json:"pred_result"`
	LabelDict  sql.NullString `json:"label_dict"`
}

type ModelPerfResponse struct {
	Rows      []map[string]interface{}          `json:"rows"`
	KeyRows   map[string]map[string]interface{} `json:"keyRows"`
	Summaries []map[string]interface{}          `json:"summaries"`
	Sum       int                               `json:"sum"`
	Acc       string                            `json:"acc"`
}

type PredResultEntity struct {
	Pred      map[string]interface{} `json:"pred"` // TODO : VCLS.ML 은 여러 개일 걸?
	PredProb  map[string][]float64   `json:"pred_prob"`
	Label     map[string]interface{} `json:"label"`      // TODO : VCLS.ML 은 여러 개일 걸?
	ImagePath map[string]string      `json:"image_path"` // TODO : Tabular는 이름이 다를 걸?
}

func (mp *ModelPerfResponse) AddRow(entity ModelPerformanceEntity) {
	row := make(map[string]interface{})

	row["modeling_id"] = entity.ModelingID
	row["model"] = entity.Model

	if score, err := mp.generateScore(entity.Score.String); err == nil {
		maps.Copy(row, score)
	}

	if inf_time, err := mp.generateInfTime(entity.InfTime.String); err == nil {
		maps.Copy(row, inf_time)
	}

	if acc, err := mp.generateAcc(entity.CfMatrix.String); err == nil {
		maps.Copy(row, acc)
	}

	mp.Rows = append(mp.Rows, row)
}

func (mp *ModelPerfResponse) AddTabularRow(entity ModelPerformanceEntity) {
	row := make(map[string]interface{})

	row["modeling_id"] = entity.ModelingID
	row["model"] = entity.Model
	row["db_model"] = entity.DBModel

	if score, err := utils.StringToMap(entity.Score.String); err == nil {
		maps.Copy(row, score)
	}

	if inf_time, err := strconv.ParseFloat(entity.InfTime.String, 64); err != nil {
		row["inf_time"] = 0.0
	} else {
		row["inf_time"] = inf_time
	}

	_str, _ := strings.CutPrefix(entity.CfMatrix.String, "\"")
	_str, _ = strings.CutSuffix(_str, "\"")
	_str = strings.ReplaceAll(_str, "\\", "")
	if acc, err := mp.generateAcc(_str); err == nil {
		maps.Copy(row, acc)
	}

	mp.Rows = append(mp.Rows, row)
}

func (mp *ModelPerfResponse) AddVisionSLConfusionMatrix(entity ModelPerformanceEntity) {
	var result map[string]map[string]interface{}
	rawString := entity.CfMatrix.String
	// 따옴표를 제거 (문자열의 첫 번째와 마지막 문자가 따옴표인 경우)
	if len(rawString) >= 2 && rawString[0] == '"' && rawString[len(rawString)-1] == '"' {
		// 따옴표 제거 및 이스케이프된 문자를 언이스케이프
		unquoted, err := strconv.Unquote(rawString)
		if err != nil {
			fmt.Println("Failed to unquote string:", err)
			return
		}
		rawString = unquoted
	}
	if err := json.Unmarshal([]byte(rawString), &result); err != nil {
		fmt.Println("Result unmarshal:", err)
		return
	}

	// 레이블 처리 - LabelDict가 없을 경우 혼동 행렬에서 키를 추출
	var labels []string
	dictString := entity.LabelDict.String

	if dictString == "" {
		// LabelDict가 없는 경우 혼동 행렬에서 키를 순서대로 추출
		// result의 키 순서를 유지하기 위해 직접 처리
		seenLabels := make(map[string]bool) // 중복 체크용
		
		// 먼저 result에서 키 순서대로 순회하면서 배열에 추가
		for outerKey, innerMap := range result {
			if !seenLabels[outerKey] {
				labels = append(labels, outerKey)
				seenLabels[outerKey] = true
			}
			
			// 내부 맵의 키도 확인
			for innerKey := range innerMap {
				if !seenLabels[innerKey] {
					labels = append(labels, innerKey)
					seenLabels[innerKey] = true
				}
			}
		}
	} else {
		// LabelDict가 있는 경우
		if len(dictString) >= 2 && dictString[0] == '"' && dictString[len(dictString)-1] == '"' {
			unquoted, err := strconv.Unquote(dictString)
			if err != nil {
				fmt.Println("Failed to unquote string:", err)
				return
			}
			dictString = unquoted
		}
		
		// JSON 디코더를 사용하여 JSON 구조를 유지하면서 파싱
		decoder := json.NewDecoder(strings.NewReader(dictString))
		// UseNumber()를 사용하여 숫자를 json.Number로 파싱
		decoder.UseNumber()
		
		// JSON 객체 파싱
		var labelDict map[string]json.Number
		if err := decoder.Decode(&labelDict); err != nil {
			fmt.Println("Failed to parse label dictionary:", err)
			return
		}
		
		// 인덱스 기반 정렬을 위한 임시 구조체
		type labelIndex struct {
			label string
			index int
		}
		
		// 레이블과 인덱스 추출
		indexedLabels := make([]labelIndex, 0, len(labelDict))
		for k, v := range labelDict {
			index, _ := v.Int64()
			indexedLabels = append(indexedLabels, labelIndex{label: k, index: int(index)})
		}
		
		// 인덱스 기준으로 정렬
		sort.Slice(indexedLabels, func(i, j int) bool {
			return indexedLabels[i].index < indexedLabels[j].index
		})
		
		// 정렬된 순서대로 레이블만 추출
		labels = make([]string, len(indexedLabels))
		for i, item := range indexedLabels {
			labels[i] = item.label
		}
	}

	totalSum, equalSum := 0, 0
	for _, v1 := range labels {
		row := map[string]interface{}{
			"label": v1,
		}

		rowSum := 0
		for _, v2 := range labels {
			value := getNumericValue(result[v2][v1])
			rowSum += value
			row[v2] = result[v2][v1]

			if v1 == v2 {
				equalSum += value
			}
			totalSum += value
		}

		row["sum"] = rowSum
		row["rec"] = divideAndRound(getNumericValue(result[v1][v1]), rowSum)
		mp.Rows = append(mp.Rows, row)
	}

	mp.Summaries = []map[string]interface{}{
		generateSumSummary(labels, result),
		generatePrecisionSummary(labels, result),
	}
	mp.Sum = totalSum
	mp.Acc = divideAndRound(equalSum, totalSum)
}

func (mp *ModelPerfResponse) AddVisionMLConfusionMatrix(entity ModelPerformanceEntity) {
	var result map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(entity.CfMatrix.String), &result); err != nil {
		fmt.Println("Result unmarshal:", err)
		return
	}

	labels, err := utils.StringToReversMap(entity.LabelDict.String)
	if err != nil {
		fmt.Println("Result tomap:", err)
		return
	}

	// generate rows
	for _, label := range labels {
		row := map[string]interface{}{
			"label": label,
			"id":    label,
		}
		for key, col := range result {
			if col[label] == nil {
				row[key] = 0
			} else {
				row[key] = col[label]
			}
		}
		mp.Rows = append(mp.Rows, row)
	}
	// generate summaries
	summary := map[string]interface{}{
		"label": "sum",
		"id":    "sum",
	}
	for key, col := range result {
		if key == "fn" || key == "fp" || key == "tn" || key == "tp" || key == "sum" {
			if col["sum"] == nil {
				summary[key] = 0
			} else {
				summary[key] = col["sum"]
			}
		}
	}
	mp.Summaries = append(mp.Summaries, summary)
}

func (mp *ModelPerfResponse) AddPredResult(entity ModelPerformanceEntity) {
	predStruct := PredResultEntity{}
	if err := json.Unmarshal([]byte(entity.PredResult.String), &predStruct); err != nil {
		fmt.Println("AddPredResult unmarshal", err)
	} else if labels, err := utils.StringToReversMap(entity.LabelDict.String); err != nil {
		fmt.Println("AddPredResult tomap", err)
	} else {
		for i := 0; i < len(predStruct.Pred); i++ {
			k := strconv.Itoa(i)
			row := make(map[string]interface{})

			row["sample"] = predStruct.ImagePath[k]
			row["isCorrect"] = (predStruct.Label[k] == predStruct.Pred[k])

			tmp := fmt.Sprintf("%v", predStruct.Label[k])
			anss := strings.Split(tmp, ",")
			ansStr := ""
			for _, v := range anss {
				ansStr = ansStr + labels[v] + ","
			}
			ansStr = strings.TrimSuffix(ansStr, ",")
			row["ans"] = ansStr

			tmp2 := fmt.Sprintf("%v", predStruct.Pred[k])
			infers := strings.Split(tmp2, ",")
			inferStr := ""
			for _, v := range infers {
				inferStr = inferStr + labels[v] + ","
			}
			inferStr = strings.TrimSuffix(inferStr, ",")
			row["infer"] = inferStr
			for k2, v2 := range labels {
				index, _ := strconv.Atoi(k2)
				row[v2] = predStruct.PredProb[k][index]
			}

			mp.Rows = append(mp.Rows, row)
		}
	}
}

func (mp *ModelPerfResponse) AddHeatmapResult(entity ModelPerformanceEntity, task_id int, engine_type string, modeling_id int, dataset_type string, model_name string) {
	predStruct := PredResultEntity{}
	if err := json.Unmarshal([]byte(entity.PredResult.String), &predStruct); err != nil {
		fmt.Println("AddPredResult unmarshal", err)
	} else {
		mp.KeyRows = make(map[string]map[string]interface{})
		for i := 0; i < len(predStruct.Pred); i++ {
			k := strconv.Itoa(i)
			row := make(map[string]interface{})
			originPath := predStruct.ImagePath[k]
			overlayPath := generatePath("overlay", originPath, task_id, engine_type, modeling_id, dataset_type, model_name)
			heatmapPath := generatePath("heatmap", originPath, task_id, engine_type, modeling_id, dataset_type, model_name)
			originURL, err := utils.ImagePathToURL(originPath)
			if err != nil {
				logger.Debug("Failed to change origin image path to URL")
			}

			heatmapURL, err := utils.ImagePathToURL(heatmapPath)
			if err != nil {
				logger.Debug("Failed to change heatmap image path to URL")
			}

			overlayURL, err := utils.ImagePathToURL(overlayPath)
			if err != nil {
				logger.Debug("Failed to change overlay image path to URL")
			}

			row["origin"] = originURL
			row["heatmap"] = heatmapURL
			row["overlay"] = overlayURL

			mp.KeyRows[originPath] = row
		}
	}
}

func (mp *ModelPerfResponse) generateScore(source string) (map[string]interface{}, error) {
	if score, err := utils.StringToMapSlice(source); err != nil {
		return nil, err
	} else {
		result := make(map[string]interface{})
		for k, v := range score {
			if _v, ok := v[0].(float64); ok {
				result[k] = _v
			}
		}

		return result, nil
	}
}

func (mp *ModelPerfResponse) generateInfTime(source string) (map[string]interface{}, error) {
	if inf_time, err := utils.StringToMap(source); err != nil {
		return nil, err
	} else {
		result := make(map[string]interface{})
		result["mean"] = inf_time["avg inference time"]
		result["min"] = inf_time["min inference time"]
		result["max"] = inf_time["max inference time"]
		result["stdev"] = inf_time["std inference time"]

		return result, nil
	}
}

func (mp *ModelPerfResponse) generateAcc(source string) (map[string]interface{}, error) {
	if acc, err := utils.StringToMapMap(source); err != nil {
		return nil, err
	} else if strings.HasPrefix(source, "{\"fn\"") {
		result := make(map[string]interface{})

		sum := make(map[string]int)
		correct := make(map[string]int)
		for k, v := range acc {
			if k == "label_accuracy" || k == "label_precision" || k == "label_f1_score" || k == "label_recall" || k == "sum" {
				continue
			}
			for k2, v2 := range v {
				if _v, ok := v2.(float64); ok {
					if k == "tp" || k == "tn" {
						correct[k2] = correct[k2] + int(_v)
					}
					sum[k2] = sum[k2] + int(_v)
				}
			}
		}

		for k := range sum {
			result[k+"_acc"] = float64(correct[k]) / float64(sum[k])
			result[k+"_ratio"] = fmt.Sprintf("%d/%d", correct[k], sum[k])
		}

		return result, nil
	} else {
		result := make(map[string]interface{})

		sum := make(map[string]int)
		correct := make(map[string]int)
		for k, v := range acc {
			for k2, v2 := range v {
				if _v, ok := v2.(float64); ok {
					sum[k2] = sum[k2] + int(_v)
					if k == k2 {
						correct[k2] = correct[k2] + int(_v)
					}
				}
			}
		}

		for k := range acc {
			result[k+"_acc"] = float64(correct[k]) / float64(sum[k])
			result[k+"_ratio"] = fmt.Sprintf("%d/%d", correct[k], sum[k])
		}

		return result, nil
	}
}

func generateSumSummary(labels []string, result map[string]map[string]interface{}) map[string]interface{} {
	summary := map[string]interface{}{
		"label": "sum",
		"id":    "sum",
	}

	for _, v1 := range labels {
		sum := 0
		for _, v2 := range labels {
			sum += getNumericValue(result[v1][v2])
		}
		summary[v1] = sum
	}

	return summary
}

func generatePrecisionSummary(labels []string, result map[string]map[string]interface{}) map[string]interface{} {
	summary := map[string]interface{}{
		"label": "precision",
		"id":    "precision",
	}

	for _, v1 := range labels {
		sum := 0
		diagValue := 0
		for _, v2 := range labels {
			value := getNumericValue(result[v1][v2])
			sum += value
			if v1 == v2 {
				diagValue = value
			}
		}
		summary[v1] = divideAndRound(diagValue, sum)
	}

	return summary
}

func getNumericValue(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	default:
		return 0
	}
}

func divideAndRound(v1 int, v2 int) string {
	if v2 == 0 {
		return "-"
	}
	percentage := float64(v1) / float64(v2) * 100
	return fmt.Sprintf("%.2f", percentage)
}

func generatePath(imgType string, originPath string, task_id int, engine_type string, modeling_id int, dataset_type string, model_name string) string {
	filename := filepath.Base(originPath)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	newFileName := fmt.Sprintf("%s_%s.jpg", name, imgType)

	cf := config_service.NewStatic()
	root_path := cf.Get("ROOT_PATH")

	finalHeatmapPath := filepath.Join(root_path, "task", engine_type,
		fmt.Sprintf("%d", task_id), fmt.Sprintf("%d", modeling_id),
		dataset_type, model_name, newFileName)

	return finalHeatmapPath
}

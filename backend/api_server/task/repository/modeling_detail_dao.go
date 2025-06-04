package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"api_server/ent"
	"api_server/ent/modeling"
	"api_server/ent/modelingdetails"
	"api_server/logger"
	"api_server/utils"
)

type ModelingDetailDAO struct {
	dbms *ent.Client
	ctx  context.Context
}

type IModelinDetailDAO interface {
	SelectModelingType(modeling_id int) (string, error)
	SelectLossChart(modeling_id int) (*LossPerfChartResponse, error)
	SelectTabularChart(modeling_id int) (*TabularChartResponse, error)
	SelectModelPerformance(modeling_id int, dataset_type string) (*ModelPerfResponse, error)
	SelectTabularModelPerformance(modeling_id int, dataset_type string) (*ModelPerfResponse, error)
	//SelectThreshold(modeling_id int, model_name string) ([]map[string]interface{}, error)
	SelectThreshold(modeling_id int, dataset_type string, model_name string) (map[string][]interface{}, error)
	//SelectFeatureImportanceChart(modeling_id int) (map[string][]interface{}, error)
	SelectFeatureImportanceChart(modeling_id int, dataset_type string, model_name string) (map[string]interface{}, error)

	// SelectOne은 주어진 ID에 해당하는 단일 ModelingDetails 데이터를 조회하는 함수입니다.
	//
	// 매개변수:
	//   - ctx: 실행 컨텍스트
	//   - id: 조회할 ModelingDetails 데이터의 ID
	//
	// 반환 값:
	//   - *ent.ModelingDetails: 조회된 ModelingDetails 객체
	//   - error: 조회 중 발생한 오류
	SelectOne(ctx context.Context, id int) (*ent.ModelingDetails, error)
	// SelectMany는 주어진 IDs에 해당하는 여러 ModelingDetails 데이터를 조회하는 함수입니다.
	//
	// 매개변수:
	//   - ctx: 실행 컨텍스트
	//   - ids: 조회할 ModelingDetails 데이터의 ID 목록
	//
	// 반환 값:
	//   - []*ent.ModelingDetails: 조회된 ModelingDetails 객체들의 슬라이스
	//   - error: 조회 중 발생한 오류
	SelectMany(ctx context.Context, ids []int) ([]*ent.ModelingDetails, error)
	// SelectIDsByModelingIDs는 주어진 modelingIDs에 해당하는 ModelingDetails의 ID 목록을 조회하는 함수입니다.
	//
	// 매개변수:
	//   - ctx: 실행 컨텍스트
	//   - modelingIDs: 조회할 Modeling의 ID 목록
	//
	// 반환 값:
	//   - []int: 조회된 ModelingDetails의 ID 목록
	//   - error: 조회 중 발생한 오류
	SelectIDsByModelingIDs(ctx context.Context, modelingIDs []int) ([]int, error)
	// InsertOne은 새로운 ModelingDetails 데이터를 데이터베이스에 삽입하는 함수입니다.
	//
	// 매개변수:
	//   - ctx: 실행 컨텍스트
	//   - req: 삽입할 ModelingDetails 데이터 객체 (DTO 형식)
	//
	// 반환 값:
	//   - *ent.ModelingDetails: 삽입된 ModelingDetails 객체
	//   - error: 삽입 중 발생한 오류
	InsertOne(ctx context.Context, req ModelingDetailDTO) (*ent.ModelingDetails, error)

	// SelectManyByModelingIDAndModelName 함수는 특정 모델링 ID와 모델 이름에 해당하는 상세 모델링 데이터를 조회합니다.
	//
	// 매개변수:
	//   - ctx: 요청의 컨텍스트 (타임아웃, 취소 등을 제어)
	//   - modelingID: 조회할 모델링 ID
	//   - modelName: 조회할 모델 이름
	//
	// 반환값:
	//   - []*ent.ModelingDetails: 조회된 모델링 상세 데이터 리스트
	//   - error: 쿼리 또는 처리 중 발생한 오류
	SelectManyByModelingIDAndModelName(ctx context.Context, modelingID int, modelName string) ([]*ent.ModelingDetails, error)
}

func NewModelingDetailDAO() *ModelingDetailDAO {
	return &ModelingDetailDAO{
		dbms: utils.GetEntClient(),
		ctx:  context.Background(),
	}
}

func (dao *ModelingDetailDAO) SelectModelingType(modeling_id int) (string, error) {
	if modeling, err := dao.dbms.Modeling.Query().
		Select(modeling.FieldID, modeling.FieldTaskID, modeling.FieldTaskID).
		WithTask().
		Where(modeling.ID(modeling_id)).
		Only(dao.ctx); err != nil {
		return "", err
	} else {
		return modeling.Edges.Task.EngineType, nil
	}
}

func (dao *ModelingDetailDAO) SelectLossChart(modeling_id int) (*LossPerfChartResponse, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d}`, modeling_id))
	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`SELECT TR.id, TS.status_json, M.modeling_step
		FROM trial_status TS
			JOIN trial TR ON TR.uuid = TS.trial_uuid
			JOIN modeling M ON M.id = TR.modeling_id
		WHERE TR.modeling_id = %d
		ORDER BY TR.id ASC, TS.id ASC;
		`,
			modeling_id),
	)

	if err != nil {
		return nil, err
	}

	trialNo := 0
	lastTrial := 0
	lastEpoch := 0
	results := LossPerfChartResponse{
		LossDims: []string{"epoch", "train_loss", "valid_loss"},
		PerfDims: []string{"epoch", "train_score", "valid_score"},
	}
	for rows.Next() {
		result := StatusJsonEntity{}
		if err := rows.Scan(
			&result.TrialID, &result.StatusJson, &result.ModelingStep,
		); err != nil {
			fmt.Println(err)
			continue
		}

		if lastTrial < result.TrialID {
			if trialNo > 0 {
				results.AddEmptyData(trialNo, lastEpoch)
			}
			trialNo++
			lastTrial = result.TrialID
			lastEpoch = 0
		}

		result.TrialNo = trialNo
		lastEpoch = results.AddChartData(result)
	}

	return &results, nil
}

func (dao *ModelingDetailDAO) SelectTabularChart(modeling_id int) (*TabularChartResponse, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d}`, modeling_id))
	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`SELECT TS.status_json->'metrics'->T.target_metric, M.modeling_step, T.target_metric
		FROM trial_status TS
			JOIN trial TR ON TR.uuid = TS.trial_uuid
			JOIN modeling M ON M.id = TR.modeling_id
			JOIN task T ON T.id = M.task_id 
		WHERE TR.modeling_id = %d
		ORDER BY TR.id ASC, TS.id ASC;
		`,
			modeling_id),
	)

	if err != nil {
		return nil, err
	}

	results := TabularChartResponse{
		Dims: []string{"trial_no", "score"},
	}
	trialNo := 1
	var targetMetric string
	var modelingStep string
	for rows.Next() {
		var value float64

		if err := rows.Scan(&value, &modelingStep, &targetMetric); err != nil {
			fmt.Println(err)
			continue
		}

		item := TabularChartItem{TrialNo: strconv.Itoa(trialNo), Score: value}

		results.Items = append(results.Items, item)

		trialNo++
	}

	results.Dims = append(results.Dims, targetMetric)
	results.TargetMetric = targetMetric
	results.ModelingStep = modelingStep

	return &results, nil
}

func (dao *ModelingDetailDAO) SelectModelPerformance(modeling_id int, dataset_type string) (*ModelPerfResponse, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d, "dataset_type": %s}`, modeling_id, dataset_type))
	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`with rawdata as (
			select model, data, data_type 
			from modeling_details 
			where modeling_id = %d and model <> 'total'
		)
		select distinct R.model
		, (SELECT data from rawdata where model = R.model and data_type = '%sset_score') as score
		, (select data from rawdata where model = R.model and data_type = '%s_avg_inference_time') as inf_time
		, (select data from rawdata where model = R.model and data_type = '%sset_cf_matrix') as cf_matrix
	from rawdata as R
	order by R.model ASC;
		`,
			modeling_id, dataset_type, dataset_type, dataset_type),
	)

	if err != nil {
		return nil, err
	}

	results := ModelPerfResponse{}
	for rows.Next() {
		result := ModelPerformanceEntity{ModelingID: modeling_id}
		if err := rows.Scan(
			&result.Model, &result.Score, &result.InfTime, &result.CfMatrix,
		); err != nil {
			fmt.Println(err)
			continue
		}

		results.AddRow(result)
	}

	return &results, nil
}

func (dao *ModelingDetailDAO) SelectTabularModelPerformance(modeling_id int, dataset_type string) (*ModelPerfResponse, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d, "dataset_type": %s}`, modeling_id, dataset_type))
	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`WITH model_data AS (
    SELECT mm.data::json AS best_model
    FROM modeling_models mm
    WHERE mm.data_type = 'best_model_dict' AND mm.modeling_id = %d
),
extracted_models AS (
    SELECT 
        metric_key,
        rank_key,
        best_model->metric_key->rank_key->0 AS file_path,
        best_model->metric_key->rank_key->2 AS model_name,
        TRIM('.kaier"' FROM SUBSTRING((best_model->metric_key->rank_key->0)::text, '([^/\\\\]+)$')) as extracted_filename,
        TRIM('"' FROM (best_model->metric_key->rank_key->2)::text) AS original_model
    FROM 
        model_data,
        json_object_keys(best_model) AS metric_key,
        json_object_keys(best_model->metric_key) AS rank_key
),
model_ids AS (
    SELECT distinct 
        original_model,
        extracted_filename,
        TRIM('"' FROM model_name::text) || '_' || 
        SUBSTRING(TRIM('"' FROM file_path::text) FROM 'model_([0-9]+)\.kaier$') AS combined_model_id
    FROM 
        extracted_models
    WHERE 
        original_model IS NOT NULL AND extracted_filename IS NOT NULL
),
rawdata AS (
    SELECT 
        m.extracted_filename AS model_name,
        md.model,
        md.data, 
        md.data_type 
    FROM 
        modeling_details md
        JOIN model_ids m ON md.model = m.combined_model_id
    WHERE 
        md.modeling_id = %d AND model IN (SELECT combined_model_id FROM model_ids)
) 
SELECT DISTINCT 
    R.model_name as model, R.model as dbname,
    (SELECT data FROM rawdata WHERE model = R.model AND data_type = '%sset_score') AS score,
    (SELECT data FROM rawdata WHERE model = R.model AND data_type = '%s_inference_time') AS inf_time,
    (SELECT data FROM rawdata WHERE model = R.model AND data_type = '%sset_cf_matrix') AS cf_matrix
FROM rawdata AS R;`,
			modeling_id, modeling_id, dataset_type, dataset_type, dataset_type),
	)

	if err != nil {
		return nil, err
	}

	results := ModelPerfResponse{}
	for rows.Next() {
		result := ModelPerformanceEntity{ModelingID: modeling_id}
		if err := rows.Scan(
			&result.Model, &result.DBModel, &result.Score, &result.InfTime, &result.CfMatrix,
		); err != nil {
			fmt.Println(err)
			continue
		}

		results.AddTabularRow(result)
	}

	return &results, nil
}

func (dao *ModelingDetailDAO) SelectVisionSLModelConfusionMatrix(modeling_id int, dataset_type string, model_name string) (*ModelPerfResponse, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d, "dataset_type": %s, "model_name": %s}`, modeling_id, dataset_type, model_name))

	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`select md.data, coalesce(m.dataset_stat->'label_dict', m.dataset_stat->'label') as label_dict
		from modeling_details md
			join modeling m on m.id = md.modeling_id
		where modeling_id = %d and data_type = '%sset_cf_matrix' and model = '%s';
		`,
			modeling_id, dataset_type, model_name),
	)

	if err != nil {
		return nil, err
	}

	results := ModelPerfResponse{}
	for rows.Next() {
		result := ModelPerformanceEntity{}
		if err := rows.Scan(
			&result.CfMatrix, &result.LabelDict,
		); err != nil {
			continue
		}
		results.AddVisionSLConfusionMatrix(result)
	}

	return &results, nil
}

func (dao *ModelingDetailDAO) SelectVisionMLModelConfusionMatrix(modeling_id int, dataset_type string, model_name string) (*ModelPerfResponse, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d, "dataset_type": %s, "model_name": %s}`, modeling_id, dataset_type, model_name))

	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`select md.data, coalesce(m.dataset_stat->'label_dict', m.dataset_stat->'label') as label_dict
		from modeling_details md
			join modeling m on m.id = md.modeling_id
		where modeling_id = %d and data_type = '%sset_cf_matrix' and model = '%s';
		`,
			modeling_id, dataset_type, model_name),
	)

	if err != nil {
		return nil, err
	}

	results := ModelPerfResponse{}
	for rows.Next() {
		result := ModelPerformanceEntity{}
		if err := rows.Scan(
			&result.CfMatrix, &result.LabelDict,
		); err != nil {
			continue
		}
		results.AddVisionMLConfusionMatrix(result)
	}

	return &results, nil
}

func (dao *ModelingDetailDAO) SelectModelSampleTest(modeling_id int, dataset_type string, model_name string) (*ModelPerfResponse, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d, "dataset_type": %s, "model_name": %s}`, modeling_id, dataset_type, model_name))
	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`select md.model, md.data, coalesce(m.dataset_stat->'label_dict', m.dataset_stat->'label') as label_dict
		from modeling_details md
			join modeling m on m.id = md.modeling_id
		where modeling_id = %d and data_type = '%s_pred_results' and model = '%s';
		`,
			modeling_id, dataset_type, model_name),
	)

	if err != nil {
		return nil, err
	}

	results := ModelPerfResponse{}
	for rows.Next() {
		result := ModelPerformanceEntity{}
		if err := rows.Scan(
			&result.Model, &result.PredResult, &result.LabelDict,
		); err != nil {
			fmt.Println(err)
			continue
		}

		results.AddPredResult(result)
	}

	return &results, nil
}

func (dao *ModelingDetailDAO) SelectHeatmapImage(task_id int, engine_type string, modeling_id int, dataset_type string, model_name string) (*ModelPerfResponse, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d, "dataset_type": %s, "model_name": %s}`, modeling_id, dataset_type, model_name))

	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`select md.data
		from modeling_details md
			join modeling m on m.id = md.modeling_id
		where modeling_id = %d and data_type = '%s_pred_results' and model = '%s';
		`,
			modeling_id, dataset_type, model_name),
	)

	if err != nil {
		return nil, err
	}

	results := ModelPerfResponse{}
	for rows.Next() {
		result := ModelPerformanceEntity{}
		if err := rows.Scan(&result.PredResult); err != nil {
			fmt.Println(err)
			continue
		}

		results.AddHeatmapResult(result, task_id, engine_type, modeling_id, dataset_type, model_name)
	}

	return &results, nil

}

func (dao *ModelingDetailDAO) SelectFeatureImportanceChart(modeling_id int, dataset_type string, model_name string) (map[string]interface{}, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d, "dataset_type": %s, "model_name": %s}`, modeling_id, dataset_type, model_name))

	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`SELECT md.data
		FROM modeling_details md
		WHERE md.modeling_id = %d and data_type = '%s_feature_importance' and model = '%s';
		`,
			modeling_id, dataset_type, model_name),
	)
	if err != nil {
		return nil, err
	}

	var row []uint8
	if rows.Next() {
		if err := rows.Scan(&row); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	var jsonString string
	if err := json.Unmarshal(row, &jsonString); err != nil {
		return nil, err
	}

	rawData := make(map[string]map[string]interface{})
	if err := json.Unmarshal([]byte(jsonString), &rawData); err != nil {
		return nil, err
	}

	sortedImportance, err := utils.SortByValue(rawData["importance"])
	if err != nil {
		logger.Debug(err)
	}

	feature := []string{}
	importance := []map[string]interface{}{}
	for _, v := range sortedImportance {
		leftSide := map[string]string{"position": "left"}
		rightSide := map[string]string{"position": "right"}

		featureValue, ok1 := rawData["feature"][v.Index].(string)
		if ok1 {
			feature = append(feature, featureValue)
		}

		data, ok2 := rawData["importance"][v.Index].(float64)
		if ok2 {
			if data > 0 {
				appendImportance := map[string]interface{}{"value": data, "label": rightSide}
				importance = append(importance, appendImportance)
			} else {
				appendImportance := map[string]interface{}{"value": data, "label": leftSide}
				importance = append(importance, appendImportance)
			}
		}
	}

	data := map[string]interface{}{
		"features":   feature,
		"importance": importance,
	}
	return data, nil

}

func (dao *ModelingDetailDAO) SelectThreshold(modeling_id int, dataset_type string, model_name string) (map[string][]interface{}, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d, "dataset_type": %s, "model_name": %s}`, modeling_id, dataset_type, model_name))

	results := make(map[string][]interface{})

	lineResults, err := dao.SelectLineThresholdData(modeling_id, model_name, results)
	if err != nil {
		return nil, fmt.Errorf("select line threshold data error: %w", err)
	}

	barResults, err := dao.SelectBarThresholdData(modeling_id, dataset_type, model_name, lineResults)
	if err != nil {
		return nil, fmt.Errorf("select bar threshold data error: %w", err)
	}

	return barResults, nil
}

func (dao *ModelingDetailDAO) SelectLineThresholdData(modeling_id int, model_name string, results map[string][]interface{}) (map[string][]interface{}, error) {
	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`select md.data
		from modeling_details md
			join modeling m on m.id = md.modeling_id
		where modeling_id = %d and data_type like 'threshold_%%' and model = '%s';`,
			modeling_id, model_name),
	)

	if err != nil {
		return nil, err
	}

	thresholdAdd := true

	row := []uint8{}
	for rows.Next() {
		result := make(map[string]map[string]interface{})
		if err := rows.Scan(&row); err != nil {
			fmt.Println(err)
			continue
		}

		if errJson := json.Unmarshal(row, &result); errJson != nil {
			fmt.Println(errJson)
		}

		for key, value_map := range result {
			if key == "max_performance" {
				continue
			}
			if key == "threshold" {
				if thresholdAdd {
					results[key], err = utils.SortObjectToMap(value_map)
					if err != nil {
						fmt.Println(err)
					}
					thresholdAdd = false
				} else {
					continue
				}
			}
			results[key], err = utils.SortObjectToMap(value_map)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return results, nil
}

func (dao *ModelingDetailDAO) SelectBarThresholdData(modeling_id int, dataset_type string, model_name string, results map[string][]interface{}) (map[string][]interface{}, error) {
	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`select md.data
		from modeling_details md
			join modeling m on m.id = md.modeling_id
		where modeling_id = %d and data_type = '%s_pred_results' and model = '%s';
		`,
			modeling_id, dataset_type, model_name),
	)

	if err != nil {
		return nil, err
	}

	row := []uint8{}
	for rows.Next() {
		result := make(map[string]interface{})
		if err := rows.Scan(&row); err != nil {
			fmt.Println(err)
			continue
		}

		if errJson := json.Unmarshal(row, &result); errJson != nil {
			fmt.Println(errJson)
		}

		correctLabels := extractCorrectLabel(result)
		thresholdResults := calculateThresholdResults(correctLabels, result)

		for _key, _value := range thresholdResults {
			interfaceSlice := make([]interface{}, len(_value))
			for i, v := range _value {
				interfaceSlice[i] = v
			}
			results[_key] = interfaceSlice
		}
	}

	return results, nil
}

func extractCorrectLabel(result map[string]interface{}) [][]string {
	labels := result["label"].(map[string]interface{})
	correctLabels := make([][]string, len(labels))
	for i := 0; i < len(labels); i++ {
		key := strconv.Itoa(i)
		if val, exists := labels[key]; exists {
			// cut between ","
			labelStr := fmt.Sprint(val)
			if strings.Contains(labelStr, ",") {
				correctLabels[i] = strings.Split(labelStr, ",")
				// trim space
				for j := range correctLabels[i] {
					correctLabels[i][j] = strings.TrimSpace(correctLabels[i][j])
				}
			} else {
				correctLabels[i] = []string{labelStr}
			}
		}
	}
	return correctLabels
}

func generatePredictionSample(correctLabels [][]string, threshold float64, predProbs map[string]interface{}, correctSample []int) []int {
	predictionSample := make([]int, len(correctLabels))

	for i := 0; i < len(correctLabels); i++ {
		key := strconv.Itoa(i)
		probArray, ok := predProbs[key].([]interface{})
		if !ok {
			continue
		}

		predictionLabels := make([]string, 0)
		for idx, prob := range probArray {
			if p, ok := prob.(float64); ok && p >= threshold {
				predictionLabels = append(predictionLabels, strconv.Itoa(idx))
			}
		}

		correctCount := 0
		for _, cLabel := range correctLabels[i] {
			for _, pLabel := range predictionLabels {
				if cLabel == pLabel {
					correctCount++
				}
			}
		}

		correctSample[i] = correctCount
		predictionSample[i] = len(predictionLabels)
	}

	return predictionSample

}

func calculateThresholdResults(correctLabels [][]string, result map[string]interface{}) map[string][]float64 {
	thresholds := make([]float64, 0)
	for t := 0.05; t <= 0.90; t += 0.05 {
		thresholds = append(thresholds, math.Round(t*100)/100)
	}

	thresholdResults := make(map[string][]float64)
	thresholdResults["corr_label_avg"] = make([]float64, len(thresholds))
	thresholdResults["pred_label_avg"] = make([]float64, len(thresholds))

	predProbs := result["pred_prob"].(map[string]interface{})

	for thresholdIndex, threshold := range thresholds {
		correctSample := make([]int, len(correctLabels))
		predictionSample := generatePredictionSample(correctLabels, threshold, predProbs, correctSample)

		var corrSum, predSum float64
		for i := 0; i < len(correctSample); i++ {
			corrSum += float64(correctSample[i])
			predSum += float64(predictionSample[i])
		}

		thresholdResults["corr_label_avg"][thresholdIndex] = corrSum / float64(len(correctSample))
		thresholdResults["pred_label_avg"][thresholdIndex] = predSum / float64(len(predictionSample))
	}
	return thresholdResults
}

// SelectOne은 주어진 ID에 해당하는 단일 ModelingDetails 데이터를 조회하는 함수입니다.
func (dao *ModelingDetailDAO) SelectOne(ctx context.Context, id int) (*ent.ModelingDetails, error) {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))

	// 주어진 ID에 해당하는 단일 ModelingDetails 데이터 조회
	return dao.dbms.ModelingDetails.Query().
		Select(
			modelingdetails.FieldID,         // ID 필드 선택
			modelingdetails.FieldModel,      // Model 필드 선택
			modelingdetails.FieldDataType,   // DataType 필드 선택
			modelingdetails.FieldData,       // Data 필드 선택
			modelingdetails.FieldCreatedAt,  // CreatedAt 필드 선택
			modelingdetails.FieldModelingID, // ModelingID 필드 선택
		).
		Where(modelingdetails.ID(id)). // 주어진 ID에 해당하는 데이터를 필터링
		Only(ctx)                      // 단일 결과만 반환
}

// SelectMany는 주어진 IDs에 해당하는 여러 ModelingDetails 데이터를 조회하는 함수입니다.
func (dao *ModelingDetailDAO) SelectMany(ctx context.Context, ids []int) ([]*ent.ModelingDetails, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_detail_ids": %v}`, ids))

	// 주어진 IDs에 해당하는 여러 ModelingDetails 데이터를 조회
	return dao.dbms.ModelingDetails.Query().
		Select().                            // 선택할 필드 지정 (모든 필드 선택)
		Where(modelingdetails.IDIn(ids...)). // 주어진 IDs에 해당하는 데이터를 필터링
		All(ctx)                             // 결과를 모두 조회
}

// SelectIDsByModelingIDs는 주어진 modelingIDs에 해당하는 ModelingDetails의 ID 목록을 조회하는 함수입니다.
func (dao *ModelingDetailDAO) SelectIDsByModelingIDs(ctx context.Context, modelingIDs []int) ([]int, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_ids": %v}`, modelingIDs))

	// 주어진 modelingIDs에 해당하는 ModelingDetails의 ID 목록을 조회
	return dao.dbms.ModelingDetails.Query().
		Where(modelingdetails.ModelingIDIn(modelingIDs...)). // 주어진 modelingIDs에 해당하는 데이터를 필터링
		Select(modelingdetails.FieldID).                     // ID 필드만 선택
		Ints(ctx)                                            // 결과를 정수형 ID 목록으로 반환
}

// InsertOne은 새로운 ModelingDetails 데이터를 데이터베이스에 삽입하는 함수입니다.
func (dao *ModelingDetailDAO) InsertOne(ctx context.Context, req ModelingDetailDTO) (*ent.ModelingDetails, error) {
	logger.Debug(fmt.Sprintf("%+v", req))

	// 새로운 ModelingDetails 레코드를 데이터베이스에 삽입
	return dao.dbms.ModelingDetails.Create().
		SetModel(req.Model).           // Model 필드 설정
		SetDataType(req.DataType).     // DataType 필드 설정
		SetData(req.Data).             // Data 필드 설정
		SetModelingID(req.ModelingID). // ModelingID 필드 설정
		SetCreatedAt(time.Now()).      // CreatedAt 필드 설정 (현재 시간)
		Save(ctx)                      // 데이터베이스에 저장
}

// SelectManyByModelingIDAndModelName 함수는 특정 모델링 ID와 모델 이름에 해당하는 상세 모델링 데이터를 조회합니다.
func (dao *ModelingDetailDAO) SelectManyByModelingIDAndModelName(
	ctx context.Context,
	modelingID int,
	modelName string,
) ([]*ent.ModelingDetails, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d, "model_name": "%s"}`, modelingID, modelName))

	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`
		SELECT id, modeling_id, model, data_type, data, created_at
		FROM modeling_details
		WHERE modeling_id = %d AND model = '%s'
		ORDER BY data_type ASC;
		`,
			modelingID, modelName,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("querying modeling_details: %w", err)
	}
	defer rows.Close()

	// 쿼리 결과를 순회하면서 데이터 파싱
	var results []*ent.ModelingDetails
	for rows.Next() {
		var md ent.ModelingDetails
		var data string // DB에는 string 형태로 저장되어 있음

		// 결과 행을 구조체에 매핑
		err := rows.Scan(
			&md.ID,
			&md.ModelingID,
			&md.Model,
			&md.DataType,
			&data,
			&md.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		// data 컬럼은 []string 형태로 반환해야 하므로 슬라이스로 변환
		md.Data = []string{data}
		results = append(results, &md)
	}

	// 순회 중 에러가 있었는지 확인
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return results, nil
}

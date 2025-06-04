package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"api_server/ent"
	"api_server/ent/modeling"
	"api_server/ent/modelingmodels"
	"api_server/logger"
	"api_server/utils"

	"entgo.io/ent/dialect/sql"
)

type IModelingDAO interface {
	InsertOne(ctx context.Context, req ModelingDTO) (*ent.Modeling, error)
	SelectManyIdle(ctx context.Context) ([]*ent.Modeling, error)
	SelectManyFinish(ctx context.Context) ([]*ent.Modeling, error)
	SelectByTask(ctx context.Context, task_id int) ([]*ModelingDB, error)
	SelectModelingType(ctx context.Context, task_id int) ([]*ModelingDB, error)
	SelectOne(ctx context.Context, id int) (*ent.Modeling, error)
	SelectFull(ctx context.Context, id int) (*ModelingDB, error)
	UpdateCancel(ctx context.Context) error
	UpdateState(ctx context.Context, modeling_id int, state string) error
	UpdateParams(ctx context.Context, modeling_id int, params []string) error
	SelectBestModelByModeling(ctx context.Context, modeling_id int) (string, error)
	SelectBestModelsByModelingId(ctx context.Context, modeling_id int) (string, error)
	SelectTestScoreByModelAndModeling(ctx context.Context, model string, modeling_id int) map[string]float64
	SelectTestInfTimeByModelAndModeling(ctx context.Context, model string, modeling_id int) float64
	DeleteOne(ctx context.Context, id int) error
	DeleteModelingAndTask(ctx context.Context, id int, taskId int) error

	// SelectMany는 주어진 IDs에 해당하는 Modeling 데이터를 조회하는 함수입니다.
	//
	// 매개변수:
	//   - ctx: 실행 컨텍스트
	//   - ids: 조회할 Modeling 데이터의 ID 목록
	//
	// 반환 값:
	//   - []*ent.Modeling: 조회된 Modeling 데이터의 슬라이스
	//   - error: 조회 중 발생한 오류
	SelectMany(ctx context.Context, ids []int) ([]*ent.Modeling, error)
	// SelectManyModelingModels는 주어진 IDs에 해당하는 ModelingModels 데이터를 조회하는 함수입니다.
	//
	// 매개변수:
	//   - ctx: 실행 컨텍스트
	//   - ids: 조회할 ModelingModels 데이터의 ID 목록
	//
	// 반환 값:
	//   - []*ent.ModelingModels: 조회된 ModelingModels 데이터의 슬라이스
	//   - error: 조회 중 발생한 오류
	SelectManyModelingModels(ctx context.Context, ids []int) ([]*ent.ModelingModels, error)
	// SelectIDsByTaskIDs는 주어진 taskIDs에 해당하는 Modeling의 ID 목록을 조회하는 함수입니다.
	//
	// 매개변수:
	//   - ctx: 실행 컨텍스트
	//   - taskIDs: 조회할 task의 ID 목록
	//
	// 반환 값:
	//   - []int: 조회된 Modeling의 ID 목록
	//   - error: 조회 중 발생한 오류
	SelectIDsByTaskIDs(ctx context.Context, taskIDs []int) ([]int, error)
	// SelectModelingModelsIDsByModelingIDs는 주어진 modelingIDs에 해당하는 ModelingModels의 ID 목록을 조회하는 함수입니다.
	//
	// 매개변수:
	//   - ctx: 실행 컨텍스트
	//   - modelingIDs: 조회할 Modeling의 ID 목록
	//
	// 반환 값:
	//   - []int: 조회된 ModelingModels의 ID 목록
	//   - error: 조회 중 발생한 오류
	SelectModelingModelsIDsByModelingIDs(ctx context.Context, modelingIDs []int) ([]int, error)
	// InsertModelingModels는 새로운 ModelingModels 데이터를 데이터베이스에 삽입하는 함수입니다.
	//
	// 매개변수:
	//   - ctx: 실행 컨텍스트
	//   - req: 삽입할 ModelingModels 데이터 객체
	//
	// 반환 값:
	//   - *ent.ModelingModels: 삽입된 ModelingModels 객체
	//   - error: 삽입 중 발생한 오류
	InsertModelingModels(ctx context.Context, req ModelingModels) (*ent.ModelingModels, error)
	SelectModelingModelsByTypeAndModelingID(ctx context.Context, modeling_id int, data_type string) (*ent.ModelingModels, error)
}

type ModelingDAO struct {
	dbms *ent.Client
}

var onceModeling sync.Once
var instanceModeling *ModelingDAO

func NewModelingDAO() *ModelingDAO {
	onceModeling.Do(func() {
		logger.Debug("Modeling DAO instance")
		instanceModeling = &ModelingDAO{
			dbms: utils.GetEntClient(),
		}
	})

	return instanceModeling
}

func (dao *ModelingDAO) InsertOne(ctx context.Context, req ModelingDTO) (*ent.Modeling, error) {
	logger.Debug(fmt.Sprintf("%+v", req))
	parent, err := dao.dbms.Modeling.Query().
		Select(modeling.FieldID, modeling.FieldLocalID).
		Where(modeling.ID(req.ParentID)).
		Only(ctx)
	if err != nil {
		logger.Debug(err)
	}
	if parent == nil {
		parent = &ent.Modeling{LocalID: 0}
	}

	return dao.dbms.Modeling.Create().
		SetParams(req.Params).
		SetModelingType(req.ModelingType).
		SetModelingStep(req.ModelingStep).
		SetProgress(0).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetTaskID(req.TaskID).
		SetParentID(req.ParentID).
		SetLocalID(dao.dbms.Modeling.
			Query().
			Where(modeling.TaskID(req.TaskID)).
			CountX(ctx) + 1).
		SetParentLocalID(parent.LocalID).
		Save(ctx)
}

func (dao *ModelingDAO) SelectManyIdle(ctx context.Context) ([]*ent.Modeling, error) {
	// logger.Debug("Select idle modelings")
	return dao.dbms.Modeling.Query().
		Select(modeling.FieldID, modeling.FieldParams, modeling.FieldModelingType).
		Where(modeling.ModelingStep(utils.MODELING_STEP_IDLE)).
		Order(modeling.ByCreatedAt(sql.OrderAsc())).
		All(ctx)
}

func (dao *ModelingDAO) SelectManyFinish(ctx context.Context) ([]*ent.Modeling, error) {
	// logger.Debug("Select finish modelings")
	return dao.dbms.Modeling.Query().
		Select(modeling.FieldID, modeling.FieldParams, modeling.FieldModelingType).
		Where(modeling.ModelingStep(utils.MODELING_STEP_FINISH)).
		Order(modeling.ByUpdatedAt(sql.OrderAsc())).
		All(ctx)
}

func (dao *ModelingDAO) SelectByTask(ctx context.Context, task_id int) ([]*ModelingDB, error) {
	logger.Debug(fmt.Sprintf(`{"task_id": %d}`, task_id))
	rows, err := dao.dbms.QueryContext(
		ctx,
		fmt.Sprintf(
			`SELECT id, local_id, task_id, parent_id, parent_local_id, dataset_id
				, modeling_type, modeling_step, params
				, dataset_stat, performance, progress
				, created_at, updated_at
			FROM modeling
			WHERE task_id = %d
			ORDER BY id DESC
			`,
			task_id,
		),
	)

	if err != nil {
		return nil, err
	}

	results := []*ModelingDB{}
	for rows.Next() {
		result := ModelingDB{}
		if err := rows.Scan(
			&result.ID, &result.LocalID, &result.TaskID, &result.ParentID, &result.ParentLocalID, &result.DatasetID,
			&result.ModelingType, &result.ModelingStep, &result.Params,
			&result.DatasetStat, &result.Performance, &result.Progress,
			&result.CreatedAt, &result.UpdatedAt,
		); err != nil {
			fmt.Println(err)
			continue
		}

		results = append(results, &result)
	}

	return results, nil
}

func (dao *ModelingDAO) SelectModelingType(ctx context.Context, task_id int) ([]*ModelingDB, error) {
	logger.Debug(fmt.Sprintf(`{"task_id": %d}`, task_id))
	rows, err := dao.dbms.QueryContext(
		ctx,
		fmt.Sprintf(
			`SELECT id, local_id, task_id, parent_id, parent_local_id, dataset_id
				, modeling_type, modeling_step, params
				, dataset_stat, performance, progress
				, created_at ,updated_at, started_at
			FROM modeling
			WHERE task_id = %d AND modeling_type IN ('%s')
			`,
			task_id,
			utils.MODELING_TYPE_INITIAL,
			//utils.MODELING_TYPE_UPDATE,
		),
	)

	if err != nil {
		return nil, err
	}

	results := []*ModelingDB{}
	for rows.Next() {
		result := ModelingDB{}
		if err := rows.Scan(
			&result.ID, &result.LocalID, &result.TaskID, &result.ParentID, &result.ParentLocalID, &result.DatasetID,
			&result.ModelingType, &result.ModelingStep, &result.Params,
			&result.DatasetStat, &result.Performance, &result.Progress,
			&result.CreatedAt, &result.UpdatedAt,
		); err != nil {
			fmt.Println(err)
			continue
		}

		results = append(results, &result)
	}

	return results, nil
}

func (dao *ModelingDAO) SelectOne(ctx context.Context, id int) (*ent.Modeling, error) {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	return dao.dbms.Modeling.Query().
		Select(
			modeling.FieldID,
			modeling.FieldLocalID,
			modeling.FieldParams,
			modeling.FieldModelingStep,
			modeling.FieldModelingType,
			modeling.FieldTaskID,
			modeling.FieldParentID,
			modeling.FieldParentLocalID,
			modeling.FieldDatasetID,
		).
		Where(modeling.ID(id)).
		Only(ctx)
}

func (dao *ModelingDAO) SelectFull(ctx context.Context, id int) (*ModelingDB, error) {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	rows, err := dao.dbms.QueryContext(
		ctx,
		fmt.Sprintf(
			`SELECT id, local_id, task_id, parent_id, parent_local_id, dataset_id
				, modeling_type, modeling_step, params
				, dataset_stat, performance, progress
				, created_at, updated_at, started_at
			FROM modeling
			WHERE id = %d
			`,
			id,
		),
	)

	if err != nil {
		return nil, err
	}

	results := []*ModelingDB{}
	for rows.Next() {
		result := ModelingDB{}
		if err := rows.Scan(
			&result.ID, &result.LocalID, &result.TaskID, &result.ParentID, &result.ParentLocalID, &result.DatasetID,
			&result.ModelingType, &result.ModelingStep, &result.Params,
			&result.DatasetStat, &result.Performance, &result.Progress,
			&result.CreatedAt, &result.UpdatedAt,
		); err != nil {
			fmt.Println(err)
			continue
		}

		results = append(results, &result)
	}

	return results[0], nil
}

func (dao *ModelingDAO) UpdateCancel(ctx context.Context) error {
	// logger.Debug("cancel modeling tasks")
	return dao.dbms.Modeling.Update().
		Where(modeling.ModelingStep(utils.MODELING_STEP_RUN)).
		SetModelingStep(utils.MODELING_STEP_CANCEL).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
}

func (dao *ModelingDAO) UpdateState(ctx context.Context, modeling_id int, state string) error {
	logger.Debug(fmt.Sprintf(`{"id": %d, "state": %s}`, modeling_id, state))
	return dao.dbms.Modeling.Update().
		Where(modeling.ID(modeling_id)).
		SetModelingStep(state).
		SetUpdatedAt(time.Now()).
		SetStartedAt(time.Now()).
		Exec(ctx)
}

func (dao *ModelingDAO) UpdateParams(ctx context.Context, modeling_id int, params []string) error {
	return dao.dbms.Modeling.Update().
		Where(modeling.ID(modeling_id)).
		SetParams(params).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
}

func (dao *ModelingDAO) SelectBestModelByModeling(ctx context.Context, modeling_id int) (string, error) {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, modeling_id))
	rows, err := dao.dbms.QueryContext(
		ctx,
		fmt.Sprintf(`select mm.data::json->t.target_metric->'1' as best_model, t.engine_type
		from modeling_models mm
		join modeling m on m.id = mm.modeling_id 
		join task t on t.id = m.task_id 
		where mm.data_type = 'best_model_dict' and mm.modeling_id = %d;`, modeling_id),
	)

	if err != nil {
		return "", err
	}

	best_model := ""
	for rows.Next() {
		row := []uint8{}
		engine := ""
		if err := rows.Scan(&row, &engine); err != nil {
			fmt.Println(err)
		} else {
			model := strings.ReplaceAll(string(row), "\\", "/")
			model = strings.ReplaceAll(model, "//", "/")
			model = strings.ReplaceAll(model, "\"", "")
			_row := strings.Split(model, ",")

			if engine == utils.JOB_TYPE_TABLE_CLS || engine == utils.JOB_TYPE_TABLE_REG {
				model = _row[2]
				_paths := strings.Split(_row[0], "/")
				no := _paths[len(_paths)-2]
				model = strings.ReplaceAll(model, "[", "")
				model = strings.ReplaceAll(model, "]", "")
				model = strings.ReplaceAll(model, " ", "")
				best_model = model + "_" + no
			} else {
				_paths := strings.Split(_row[0], "/")
				model = _paths[len(_paths)-1]
				best_model = strings.Split(model, ".kaier")[0]
			}
		}
	}

	return best_model, nil
}

func (dao *ModelingDAO) SelectBestModelsByModelingId(ctx context.Context, modeling_id int) (string, error) {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, modeling_id))
	rows, err := dao.dbms.QueryContext(
		ctx,
		fmt.Sprintf(`select mm.data::json as best_model
			from modeling_models mm
			where mm.data_type = 'best_model_dict' and mm.modeling_id = %d;`, modeling_id),
	)

	if err != nil {
		return "", err
	}

	best_model_dict := ""
	for rows.Next() {
		row := []uint8{}
		if err := rows.Scan(&row); err != nil {
			fmt.Println(err)
		}

		best_model_dict = string(row)
	}

	return best_model_dict, nil
}

func (dao *ModelingDAO) SelectTestScoreByModelAndModeling(ctx context.Context, model string, modeling_id int) map[string]float64 {
	logger.Debug(fmt.Sprintf(`{"id": %d, "model": %s}`, modeling_id, model))
	rows, err := dao.dbms.QueryContext(
		ctx,
		fmt.Sprintf(`select data 
			from modeling_details 
			where data_type='testset_score' and model = '%s' and modeling_id = %d`,
			model, modeling_id,
		),
	)

	if err != nil {
		return nil
	}

	row := []uint8{}
	for rows.Next() {
		if err := rows.Scan(&row); err != nil {
			fmt.Println(err)
		}
	}

	score_map := make(map[string][]float64)
	if errJson := json.Unmarshal(row, &score_map); errJson == nil {
		result := make(map[string]float64)
		for k, v := range score_map {
			result[k] = v[0]
		}

		return result
	} else {
		score_m := make(map[string]float64)
		if errJson := json.Unmarshal(row, &score_m); errJson == nil {
			return score_m
		}
	}

	return nil
}

func (dao *ModelingDAO) SelectTestInfTimeByModelAndModeling(ctx context.Context, model string, modeling_id int) float64 {
	logger.Debug(fmt.Sprintf(`{"id": %d, "model": %s}`, modeling_id, model))
	rows, err := dao.dbms.QueryContext(
		ctx,
		fmt.Sprintf(`select data 
			from modeling_details 
			where data_type like concat('test', '%%', 'inference_time') and model = '%s' and modeling_id = %d;`,
			model, modeling_id,
		),
	)

	if err != nil {
		return 0.0
	}

	row := []uint8{}
	for rows.Next() {
		if err := rows.Scan(&row); err != nil {
			fmt.Println(err)
		}
	}

	score_map := make(map[string]interface{})
	if errJson := json.Unmarshal(row, &score_map); errJson == nil {
		inf_time, _ := score_map["avg inference time"].(float64)
		return inf_time
	} else {
		inf_time, err := strconv.ParseFloat(string(row), 64)
		if err == nil {
			return inf_time
		}
	}

	return 0.0
}

func (dao *ModelingDAO) DeleteOne(ctx context.Context, id int) error {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	return dao.dbms.Modeling.DeleteOneID(id).Exec(ctx)
}

func (dao *ModelingDAO) DeleteModelingAndTask(ctx context.Context, id int, taskId int) error {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	err := dao.dbms.Modeling.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	exists, err := dao.SelectByTask(ctx, taskId)
	if len(exists) < 1 && err == nil {
		return dao.dbms.Task.DeleteOneID(taskId).Exec(ctx)
	}

	return nil
}

// SelectMany는 주어진 IDs에 해당하는 Modeling 데이터를 조회하는 함수입니다.
// 매개변수:
//   - ctx: 실행 컨텍스트
//   - ids: 조회할 Modeling 데이터의 ID 목록
//
// 반환 값:
//   - []*ent.Modeling: 조회된 Modeling 데이터의 슬라이스
//   - error: 조회 중 발생한 오류
func (dao *ModelingDAO) SelectMany(ctx context.Context, ids []int) ([]*ent.Modeling, error) {
	logger.Debug(fmt.Sprintf(`{"ids": %v}`, ids))

	// 주어진 IDs에 해당하는 Modeling 데이터를 조회
	return dao.dbms.Modeling.Query().
		Select().                     // 선택할 필드 지정 (모든 필드 선택)
		Where(modeling.IDIn(ids...)). // 주어진 IDs에 해당하는 데이터를 필터링
		All(ctx)                      // 결과를 모두 조회
}

// SelectManyModelingModels는 주어진 IDs에 해당하는 ModelingModels 데이터를 조회하는 함수입니다.
// 매개변수:
//   - ctx: 실행 컨텍스트
//   - ids: 조회할 ModelingModels 데이터의 ID 목록
//
// 반환 값:
//   - []*ent.ModelingModels: 조회된 ModelingModels 데이터의 슬라이스
//   - error: 조회 중 발생한 오류
func (dao *ModelingDAO) SelectManyModelingModels(ctx context.Context, ids []int) ([]*ent.ModelingModels, error) {
	logger.Debug(fmt.Sprintf(`{"ids": %v}`, ids))

	// 주어진 IDs에 해당하는 ModelingModels 데이터를 조회
	return dao.dbms.ModelingModels.Query().
		Select().                           // 선택할 필드 지정 (모든 필드 선택)
		Where(modelingmodels.IDIn(ids...)). // 주어진 IDs에 해당하는 데이터를 필터링
		All(ctx)                            // 결과를 모두 조회
}

// SelectIDsByTaskIDs는 주어진 taskIDs에 해당하는 Modeling의 ID 목록을 조회하는 함수입니다.
// 매개변수:
//   - ctx: 실행 컨텍스트
//   - taskIDs: 조회할 task의 ID 목록
//
// 반환 값:
//   - []int: 조회된 Modeling의 ID 목록
//   - error: 조회 중 발생한 오류
func (dao *ModelingDAO) SelectIDsByTaskIDs(ctx context.Context, taskIDs []int) ([]int, error) {
	logger.Debug(fmt.Sprintf(`{"task_ids": %v}`, taskIDs))

	// 주어진 taskIDs에 해당하는 Modeling의 ID 목록을 조회
	return dao.dbms.Modeling.Query().
		Where(modeling.TaskIDIn(taskIDs...)). // 주어진 taskIDs에 해당하는 데이터를 필터링
		Select(modeling.FieldID).             // ID 필드만 선택
		Ints(ctx)                             // 결과를 정수형 ID 목록으로 반환
}

// SelectModelingModelsIDsByModelingIDs는 주어진 modelingIDs에 해당하는 ModelingModels의 ID 목록을 조회하는 함수입니다.
func (dao *ModelingDAO) SelectModelingModelsIDsByModelingIDs(ctx context.Context, modelingIDs []int) ([]int, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_ids": %v}`, modelingIDs))

	// 주어진 modelingIDs에 해당하는 ModelingModels의 ID 목록을 조회
	return dao.dbms.ModelingModels.Query().
		Where(modelingmodels.ModelingIDIn(modelingIDs...)). // 주어진 modelingIDs에 해당하는 데이터를 필터링
		Select(modelingmodels.FieldID).                     // ID 필드만 선택
		Ints(ctx)                                           // 결과를 정수형 ID 목록으로 반환
}

// InsertModelingModels는 새로운 ModelingModels 데이터를 데이터베이스에 삽입하는 함수입니다.
func (dao *ModelingDAO) InsertModelingModels(ctx context.Context, req ModelingModels) (*ent.ModelingModels, error) {
	logger.Debug(fmt.Sprintf("%+v", req))

	// 새로운 ModelingModels 레코드를 삽입
	return dao.dbms.ModelingModels.Create().
		SetDataType(req.DataType).     // DataType 필드 설정
		SetData(req.Data).             // Data 필드 설정
		SetModelingID(req.ModelingID). // ModelingID 필드 설정
		SetCreatedAt(time.Now()).      // CreatedAt 필드 설정 (현재 시간)
		Save(ctx)                      // 데이터베이스에 저장
}

func (dao *ModelingDAO) SelectModelingModelsByTypeAndModelingID(ctx context.Context, modeling_id int, data_type string) (*ent.ModelingModels, error) {
	logger.Debug(fmt.Sprintf(`{"data_type": %s, "modeling_id": %d}`, data_type, modeling_id))
	return dao.dbms.ModelingModels.Query().
		Where(modelingmodels.DataType(data_type)).
		Where(modelingmodels.ModelingID(modeling_id)).
		Only(ctx)
}

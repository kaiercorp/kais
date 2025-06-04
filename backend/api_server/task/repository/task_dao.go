package repository

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"api_server/ent"
	"api_server/ent/task"
	"api_server/logger"
	"api_server/utils"

	"entgo.io/ent/dialect/sql"
)

type ITaskDAO interface {
	InsertOne(ctx context.Context, req TaskDTO) (*ent.Task, error)
	SelectOne(ctx context.Context, id int) (*ent.Task, error)
	SelectAll(ctx context.Context) ([]*ent.Task, error)
	SelectByProject(ctx context.Context, project_id int) ([]*ent.Task, int, error)
	UpdateOne(ctx context.Context, req TaskDTO) (*ent.Task, error)
	DeleteOne(ctx context.Context, id int) error
	DeleteByProject(ctx context.Context, project_id int) error
	WithTx(ctx context.Context, fn func(tx *ent.Tx) error) error

	// SelectMany는 주어진 IDs에 해당하는 여러 Task 데이터를 조회하는 함수입니다.
	//
	// 매개변수:
	//   - ctx: 실행 컨텍스트
	//   - ids: 조회할 Task 데이터의 ID 목록
	//
	// 반환 값:
	//   - []*ent.Task: 조회된 Task 객체들의 슬라이스
	//   - error: 조회 중 발생한 오류
	SelectMany(ctx context.Context, ids []int) ([]*ent.Task, error)
	// SelectIDsByProjectIDs는 주어진 projectIDs에 해당하는 Task의 ID 목록을 조회하는 함수입니다.
	//
	// 매개변수:
	//   - ctx: 실행 컨텍스트
	//   - projectIDs: 조회할 프로젝트의 ID 목록
	//
	// 반환 값:
	//   - []int: 조회된 Task의 ID 목록
	//   - error: 조회 중 발생한 오류
	SelectIDsByProjectIDs(ctx context.Context, projectIDs []int) ([]int, error)

	SelectByEngineType(ctx context.Context, engine_type string) ([]*ent.Task, error)
}

type TaskDAO struct {
	dbms *ent.Client
}

var onceTask sync.Once
var instanceTask *TaskDAO

func NewTaskDAO() *TaskDAO {
	onceTask.Do(func() {
		logger.Debug("Task DAO instance")
		instanceTask = &TaskDAO{
			dbms: utils.GetEntClient(),
		}
	})

	return instanceTask
}

func (dao *TaskDAO) InsertOne(ctx context.Context, req TaskDTO) (*ent.Task, error) {
	logger.Debug(fmt.Sprintf("%+v", req))
	return dao.dbms.Task.Create().
		SetProjectID(*req.ProjectID).
		SetDatasetID(*req.DatasetID).
		SetTitle(req.Title).
		SetDescription(req.Description).
		SetEngineType(req.EngineType).
		SetTargetMetric(req.TargetMetric).
		SetParams(req.Params).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)
}

func (dao *TaskDAO) SelectOne(ctx context.Context, id int) (*ent.Task, error) {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	return dao.dbms.Task.Query().Select().
		Where(task.ID(id)).
		Only(ctx)
}

func (dao *TaskDAO) SelectByProject(ctx context.Context, project_id int) ([]*ent.Task, int, error) {
	logger.Debug(fmt.Sprintf(`{"page": %d}`, project_id))
	if all, err := dao.dbms.Task.Query().Select().
		Where(task.ProjectIDEQ(project_id)).
		Order(task.ByID(sql.OrderDesc())).
		All(ctx); err != nil {
		return all, 0, err
	} else {
		count := dao.dbms.Task.Query().
			Where(task.ProjectIDEQ(project_id)).
			CountX(ctx)

		return all, int(math.Ceil(float64(count) / float64(25))), nil
	}
}

func (dao *TaskDAO) UpdateOne(ctx context.Context, req TaskDTO) (*ent.Task, error) {
	logger.Debug(fmt.Sprintf("%+v", req))
	return dao.dbms.Task.
		UpdateOneID(req.ID).
		SetTitle(req.Title).
		SetDescription(req.Description).
		SetUpdatedAt(time.Now()).
		Save(ctx)
}

func (dao *TaskDAO) DeleteOne(ctx context.Context, id int) error {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	return dao.dbms.Task.
		DeleteOneID(id).
		Exec(ctx)
}

func (dao *TaskDAO) DeleteByProject(ctx context.Context, project_id int) error {
	logger.Debug(fmt.Sprintf(`{"project_id": %d}`, project_id))
	_, err := dao.dbms.Task.Delete().
		Where(task.ProjectID(project_id)).
		Exec(ctx)

	return err
}

func (dao *TaskDAO) WithTx(ctx context.Context, fn func(tx *ent.Tx) error) error {
	logger.Debug("excuete sql with transaction")
	tx, err := dao.dbms.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// SelectMany는 주어진 IDs에 해당하는 여러 Task 데이터를 조회하는 함수입니다.
func (dao *TaskDAO) SelectMany(ctx context.Context, ids []int) ([]*ent.Task, error) {
	logger.Debug(fmt.Sprintf(`{"task_ids": %v}`, ids))

	// 주어진 IDs에 해당하는 여러 Task 데이터를 조회
	return dao.dbms.Task.Query().Select().
		Where(task.IDIn(ids...)). // 주어진 IDs에 해당하는 데이터를 필터링
		All(ctx)                  // 결과를 모두 조회
}

// SelectIDsByProjectIDs는 주어진 projectIDs에 해당하는 Task의 ID 목록을 조회하는 함수입니다.
func (dao *TaskDAO) SelectIDsByProjectIDs(ctx context.Context, projectIDs []int) ([]int, error) {
	logger.Debug(fmt.Sprintf(`{"project_ids": %v}`, projectIDs))

	// 주어진 projectIDs에 해당하는 Task의 ID 목록을 조회
	return dao.dbms.Task.Query().
		Where(task.ProjectIDIn(projectIDs...)). // 주어진 projectIDs에 해당하는 데이터를 필터링
		Select(task.FieldID).                   // ID 필드만 선택
		Ints(ctx)                               // 결과를 정수형 ID 목록으로 반환
}

func (dao *TaskDAO) SelectAll(ctx context.Context) ([]*ent.Task, error) {
	logger.Debug("select all task")
	return dao.dbms.Task.Query().Select().All(ctx)
}

func (dao *TaskDAO) SelectByEngineType(ctx context.Context, engine_type string) ([]*ent.Task, error) {
	logger.Debug(fmt.Sprintf(`{"engine_type": %s}`, engine_type))
	return dao.dbms.Task.Query().Select().
		Where(task.EngineTypeEQ(engine_type)).
		All(ctx)
}

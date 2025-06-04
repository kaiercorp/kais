package repository

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"api_server/ent"
	"api_server/ent/project"
	"api_server/logger"
	"api_server/utils"

	"entgo.io/ent/dialect/sql"
)

type IProjectDAO interface {
	InsertOne(ctx context.Context, req ProjectDTO) (*ent.Project, error)
	SelectByPage(ctx context.Context, page int) ([]*ent.Project, int, bool, int, error)
	SelectAll(ctx context.Context) ([]*ent.Project, error)
	// SelectMany 함수는 주어진 프로젝트 ID 목록을 기반으로 프로젝트를 조회합니다.
	// ID 조건을 적용하여 해당하는 프로젝트를 필터링하고, ID 기준 내림차순으로 정렬하여 반환합니다.
	//
	// 만약 projectIDs가 비어있다면 오류를 반환합니다.
	//
	// 매개변수:
	//   - ctx: 컨텍스트 객체 (요청의 실행 흐름 관리)
	//   - projectIDs: 조회할 프로젝트의 ID 목록
	//
	// 반환값:
	//   - []*ent.Project: 조회된 프로젝트 목록
	//   - error: 오류 발생 시 반환되는 에러 객체
	SelectMany(ctx context.Context, projectIDs []int) ([]*ent.Project, error)
	// SelectOne는 주어진 ID에 해당하는 단일 프로젝트를 조회하는 함수입니다.
	//
	// 이 함수는 데이터베이스에서 특정 프로젝트 ID를 검색하여 반환합니다.
	// ID가 0 이하이거나, 해당 ID의 프로젝트가 존재하지 않는 경우 오류를 반환합니다.
	//
	// 매개변수:
	//   - ctx: 데이터베이스 쿼리를 위한 컨텍스트 (취소 및 타임아웃을 관리).
	//   - projectID: 조회할 프로젝트의 ID.
	//
	// 반환값:
	//   - *ent.Project: 조회된 프로젝트 엔터티의 포인터.
	//   - error: 프로젝트를 찾을 수 없거나 쿼리에 실패한 경우 오류.
	SelectOne(ctx context.Context, projectID int) (*ent.Project, error)

	UpdateOne(ctx context.Context, req ProjectDTO) (*ent.Project, error)
	DeleteOne(ctx context.Context, id int) error

	// SelectByPageWithIDs는 주어진 projectIDs에 해당하는 프로젝트만 페이지네이션하여 조회합니다.
	//
	// 매개변수:
	//   - ctx: 컨텍스트
	//   - page: 조회할 페이지 (1부터 시작)
	//   - projectIDs: 조회할 프로젝트 ID 목록
	//
	// 반환값:
	//   - []*ent.Project: 조회된 프로젝트 목록
	//   - int: 전체 페이지 수
	//   - bool: 다음 페이지 존재 여부
	//   - int: 다음 페이지 번호
	//   - error: 오류 (있다면)
	SelectByPageWithIDs(ctx context.Context, page int, projectIDs []int) ([]*ent.Project, int, bool, int, error)
}

type ProjectDAO struct {
	dbms *ent.Client
}

var once sync.Once
var instance *ProjectDAO

func New() *ProjectDAO {
	once.Do(func() { // atomic, does not allow repeating
		logger.Debug("Project DAO instance")
		instance = &ProjectDAO{
			dbms: utils.GetEntClient(),
		}
	})

	return instance
}

func (dao *ProjectDAO) InsertOne(ctx context.Context, req ProjectDTO) (*ent.Project, error) {
	logger.Debug(fmt.Sprintf("%+v", req))
	return dao.dbms.Project.Create().
		SetTitle(req.Title).
		SetDescription(req.Description).
		Save(ctx)
}

func (dao *ProjectDAO) SelectByPage(ctx context.Context, page int) ([]*ent.Project, int, bool, int, error) {
	logger.Debug(fmt.Sprintf(`{"page": %d}`, page))
	if all, err := dao.dbms.Project.Query().Select().
		Order(project.ByID(sql.OrderDesc())).
		Offset((page - 1) * 25).
		Limit(25).
		All(ctx); err != nil {
		return all, 0, false, 0, err
	} else {
		count := dao.dbms.Project.
			Query().
			CountX(ctx)
		curPage := int(math.Ceil(float64(count) / float64(25)))
		return all, curPage, (curPage > page), page + 1, nil
	}
}

func (dao *ProjectDAO) SelectAll(ctx context.Context) ([]*ent.Project, error) {
	return dao.dbms.Project.Query().
		Order(project.ByID(sql.OrderDesc())).
		All(ctx)
}

func (dao *ProjectDAO) UpdateOne(ctx context.Context, req ProjectDTO) (*ent.Project, error) {
	logger.Debug(fmt.Sprintf("%+v", req))
	return dao.dbms.Project.
		UpdateOneID(req.ID).
		SetTitle(req.Title).
		SetDescription(req.Description).
		SetUpdatedAt(time.Now()).
		Save(ctx)
}

func (dao *ProjectDAO) DeleteOne(ctx context.Context, id int) error {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	return dao.dbms.Project.
		DeleteOneID(id).
		Exec(ctx)
}

// SelectMany 함수는 주어진 프로젝트 ID 목록을 기반으로 프로젝트를 조회합니다.
// ID 조건을 적용하여 해당하는 프로젝트를 필터링하고, ID 기준 내림차순으로 정렬하여 반환합니다.
//
// 만약 projectIDs가 비어있다면 오류를 반환합니다.
//
// 매개변수:
//   - ctx: 컨텍스트 객체 (요청의 실행 흐름 관리)
//   - projectIDs: 조회할 프로젝트의 ID 목록
//
// 반환값:
//   - []*ent.Project: 조회된 프로젝트 목록
//   - error: 오류 발생 시 반환되는 에러 객체
func (dao *ProjectDAO) SelectMany(ctx context.Context, projectIDs []int) ([]*ent.Project, error) {

	if len(projectIDs) == 0 {
		return nil, fmt.Errorf("projectIDs cannot be empty")
	}

	return dao.dbms.Project.Query().
		Where(project.IDIn(projectIDs...)).
		Order(project.ByID(sql.OrderDesc())).
		All(ctx)
}

// SelectOne는 주어진 ID에 해당하는 단일 프로젝트를 조회하는 함수입니다.
func (dao *ProjectDAO) SelectOne(ctx context.Context, projectID int) (*ent.Project, error) {
	// 유효성 검사: ID가 0 이하이면 에러 반환
	if projectID <= 0 {
		return nil, fmt.Errorf("invalid project ID")
	}

	// 프로젝트 조회
	return dao.dbms.Project.Query().
		Where(project.IDEQ(projectID)).
		Only(ctx) // 단일 결과만 반환
}

// SelectByPageWithIDs는 주어진 projectIDs에 해당하는 프로젝트만 페이지네이션하여 조회합니다.
func (dao *ProjectDAO) SelectByPageWithIDs(ctx context.Context, page int, projectIDs []int) ([]*ent.Project, int, bool, int, error) {
	logger.Debug(fmt.Sprintf(`{"page": %d, "project_ids": %v}`, page, projectIDs))

	const limit = 25
	offset := (page - 1) * limit

	query := dao.dbms.Project.Query().
		Where(project.IDIn(projectIDs...))

	// 데이터 쿼리
	projects, err := query.Clone().
		Order(project.ByID(sql.OrderDesc())).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, 0, false, 0, err
	}

	// 전체 개수 쿼리
	count := query.CountX(ctx)
	totalPages := int(math.Ceil(float64(count) / float64(limit)))
	hasNext := totalPages > page

	return projects, totalPages, hasNext, page + 1, nil
}

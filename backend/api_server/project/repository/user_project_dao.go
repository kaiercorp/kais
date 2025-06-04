package repository

import (
	"api_server/ent"
	"api_server/ent/userproject"
	"api_server/logger"
	"api_server/utils"
	"context"
	"fmt"
	"sync"

	"entgo.io/ent/dialect/sql"
)

type IUserProjectDAO interface {
	// InsertOne은 새로운 사용자-프로젝트 매핑 데이터를 삽입합니다.
	//
	// 매개변수:
	//   - ctx: context.Context 객체 (데이터베이스 작업에 대한 컨텍스트)
	//   - req: 사용자-프로젝트 매핑에 대한 DTO (UserProjectDTO)
	//
	// 반환값:
	//   - *ent.UserProject: 삽입된 사용자-프로젝트 매핑 객체
	//   - error: 작업 중 발생한 오류, 없으면 nil
	InsertOne(ctx context.Context, req UserProjectDTO) (*ent.UserProject, error)

	// DeleteOne은 특정 프로젝트 ID와 사용자명을 기준으로 매핑 데이터를 삭제합니다.
	//
	// 매개변수:
	//   - ctx: context.Context 객체 (데이터베이스 작업에 대한 컨텍스트)
	//   - id: 삭제할 프로젝트의 ID
	//   - username: 삭제할 사용자명
	//
	// 반환값:
	//   - int: 삭제된 레코드의 수
	//   - error: 작업 중 발생한 오류, 없으면 nil
	DeleteOne(ctx context.Context, id int, username string) (int, error)

	// DeleteByProjectId는 특정 프로젝트 ID에 해당하는 모든 사용자-프로젝트 매핑 데이터를 삭제합니다.
	//
	// 매개변수:
	//   - ctx: context.Context 객체 (데이터베이스 작업에 대한 컨텍스트)
	//   - id: 삭제할 프로젝트의 ID
	//
	// 반환값:
	//   - int: 삭제된 레코드의 수
	//   - error: 작업 중 발생한 오류, 없으면 nil
	DeleteByProjectId(ctx context.Context, id int) (int, error)

	// DeleteByUsername은 특정 사용자명에 해당하는 모든 사용자-프로젝트 매핑 데이터를 삭제합니다.
	//
	// 매개변수:
	//   - ctx: context.Context 객체 (데이터베이스 작업에 대한 컨텍스트)
	//   - username: 삭제할 사용자명
	//
	// 반환값:
	//   - int: 삭제된 레코드의 수
	//   - error: 작업 중 발생한 오류, 없으면 nil
	DeleteByUsername(ctx context.Context, username string) (int, error)

	// SelectByUsername은 특정 사용자명에 해당하는 사용자-프로젝트 매핑 데이터를 조회합니다.
	//
	// 매개변수:
	//   - ctx: context.Context 객체 (데이터베이스 작업에 대한 컨텍스트)
	//   - username: 조회할 사용자명
	//
	// 반환값:
	//   - []*ent.UserProject: 사용자-프로젝트 매핑 객체 목록
	//   - error: 작업 중 발생한 오류, 없으면 nil
	SelectByUsername(ctx context.Context, username string) ([]*ent.UserProject, error)

	// SelectByProjectID는 특정 프로젝트 ID에 해당하는 사용자-프로젝트 매핑 데이터를 조회합니다.
	//
	// 매개변수:
	//   - ctx: context.Context 객체 (데이터베이스 작업에 대한 컨텍스트)
	//   - projectID: 조회할 프로젝트의 ID
	//
	// 반환값:
	//   - []*ent.UserProject: 사용자-프로젝트 매핑 객체 목록
	//   - error: 작업 중 발생한 오류, 없으면 nil
	SelectByProjectID(ctx context.Context, projectID int) ([]*ent.UserProject, error)

	// SelectProjectIdsByUsername은 특정 사용자명에 해당하는 프로젝트 ID 목록을 조회합니다.
	//
	// 매개변수:
	//   - ctx: context.Context 객체 (데이터베이스 작업에 대한 컨텍스트)
	//   - username: 조회할 사용자명
	//
	// 반환값:
	//   - []*ent.UserProject: 사용자-프로젝트 매핑 객체 목록 (각 객체에서 프로젝트 ID 추출)
	//   - error: 작업 중 발생한 오류, 없으면 nil
	SelectProjectIdsByUsername(ctx context.Context, username string) ([]*ent.UserProject, error)
}

// UserProjectDAO는 IUserProjectDAO 인터페이스를 구현하는 구조체입니다.
// 이 구조체는 데이터베이스 작업을 실제로 수행합니다.
type UserProjectDAO struct {
	dbms *ent.Client
}

var onceUserProejct sync.Once
var instanceUserProject *UserProjectDAO

func NewUserProject() *UserProjectDAO {
	onceUserProejct.Do(func() { // atomic, does not allow repeating
		logger.Debug("UserProject DAO instance")
		instanceUserProject = &UserProjectDAO{
			dbms: utils.GetEntClient(),
		}
	})

	return instanceUserProject
}

// InsertOne은 새로운 사용자-프로젝트 매핑 데이터를 삽입합니다.
func (dao *UserProjectDAO) InsertOne(ctx context.Context, req UserProjectDTO) (*ent.UserProject, error) {
	logger.Debug(fmt.Sprintf("%+v", req))
	return dao.dbms.UserProject.Create().
		SetUsername(req.Username).
		SetProjectID(req.ProjectId).
		Save(ctx)
}

// DeleteOne은 특정 프로젝트 ID와 사용자명을 기준으로 매핑 데이터를 삭제합니다.
func (dao *UserProjectDAO) DeleteOne(ctx context.Context, id int, username string) (int, error) {
	logger.Debug(fmt.Sprintf(`{"project_id": %d, "username": "%s"}`, id, username))
	return dao.dbms.UserProject.
		Delete().
		Where(userproject.ProjectID(id)).
		Where(userproject.Username(username)).
		Exec(ctx)
}

// DeleteByProjectId는 특정 프로젝트 ID에 해당하는 모든 사용자-프로젝트 매핑 데이터를 삭제합니다.
func (dao *UserProjectDAO) DeleteByProjectId(ctx context.Context, id int) (int, error) {
	logger.Debug(fmt.Sprintf(`{"project_id": %d}`, id))
	return dao.dbms.UserProject.
		Delete().
		Where(userproject.ProjectID(id)).
		Exec(ctx)
}

// DeleteByUsername은 특정 사용자명에 해당하는 모든 사용자-프로젝트 매핑 데이터를 삭제합니다.
func (dao *UserProjectDAO) DeleteByUsername(ctx context.Context, username string) (int, error) {
	logger.Debug(fmt.Sprintf(`{"username": %s}`, username))
	return dao.dbms.UserProject.
		Delete().
		Where(userproject.Username(username)).
		Exec(ctx)
}

// SelectByUsername은 특정 사용자명에 해당하는 사용자-프로젝트 매핑 데이터를 조회합니다.
func (dao *UserProjectDAO) SelectByUsername(ctx context.Context, username string) ([]*ent.UserProject, error) {
	return dao.dbms.UserProject.Query().
		Where(userproject.Username(username)).
		Order(userproject.ByProjectID(sql.OrderDesc())).
		All(ctx)
}

// SelectByProjectID는 특정 프로젝트 ID에 해당하는 사용자-프로젝트 매핑 데이터를 조회합니다.
func (dao *UserProjectDAO) SelectByProjectID(ctx context.Context, projectID int) ([]*ent.UserProject, error) {
	return dao.dbms.UserProject.Query().
		Where(userproject.ProjectID(projectID)).
		Order(userproject.ByUsername(sql.OrderAsc())).
		All(ctx)
}

// SelectProjectIdsByUsername은 특정 사용자명에 해당하는 프로젝트 ID 목록을 조회합니다.
func (dao *UserProjectDAO) SelectProjectIdsByUsername(ctx context.Context, username string) ([]*ent.UserProject, error) {
	return dao.dbms.UserProject.Query().
		Where(userproject.Username(username)).
		Select(userproject.FieldProjectID).
		Order(userproject.ByProjectID(sql.OrderDesc())).
		All(ctx)
}

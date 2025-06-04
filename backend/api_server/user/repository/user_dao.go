package repository

import (
	"api_server/ent"
	"api_server/ent/user"
	"api_server/logger"
	"api_server/utils"
	"context"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent/dialect/sql"
)

// IUserDAO는 사용자(User) 엔티티에 대한 데이터 접근 메서드를 정의한 인터페이스입니다.
//
// 이 인터페이스는 사용자 생성, 조회, 수정, 삭제 기능을 포함합니다.
type IUserDAO interface {

	// InsertOne은 새로운 사용자 레코드를 생성하여 데이터베이스에 저장합니다.
	//
	// 입력:
	//   - ctx: 요청 범위를 제어하기 위한 context.Context
	//   - req: 사용자 생성에 필요한 정보를 담은 UserDTO 구조체
	//
	// 처리:
	//   - 필수 필드(name, username, password, group, is_use)가 비어있으면 오류 반환
	//   - 해당 정보를 기반으로 ent.User.Create() 호출
	//
	// 반환:
	//   - 성공 시 생성된 *ent.User 객체
	//   - 실패 시 오류 정보 (필수값 누락, DB 저장 오류 등)
	InsertOne(ctx context.Context, req UserDTO) (*ent.User, error)

	// SelectAll은 모든 사용자 레코드를 조회하며, ID 기준 내림차순으로 정렬하여 반환합니다.
	//
	// 입력:
	//   - ctx: 요청 범위를 제어하기 위한 context.Context
	//
	// 반환:
	//   - []*ent.User 형태의 전체 사용자 목록
	//   - 쿼리 중 오류 발생 시 에러 반환
	SelectAll(ctx context.Context) ([]*ent.User, error)

	// SelectOne은 지정된 ID 값을 가진 단일 사용자 레코드를 조회합니다.
	//
	// 입력:
	//   - ctx: 요청 범위를 제어하기 위한 context.Context
	//   - id: 조회 대상 사용자의 고유 ID
	//
	// 처리:
	//   - dbms.User.Get(ctx, id)를 통해 직접 조회
	//
	// 반환:
	//   - 존재하는 경우 *ent.User 객체 반환
	//   - 존재하지 않거나 오류가 발생하면 에러 반환
	SelectOne(ctx context.Context, id int) (*ent.User, error)

	// SelectOneByUsername은 주어진 사용자명(username)에 해당하는 사용자 레코드를 조회합니다.
	//
	// 입력:
	//   - ctx: 요청 범위를 제어하기 위한 context.Context
	//   - username: 검색 대상 사용자명
	//
	// 처리:
	//   - ent.User.Query().Where(user.Username(username)).Only() 호출
	//   - username은 고유해야 하며, 여러 개 존재하면 오류 발생
	//
	// 반환:
	//   - 성공 시 해당 사용자 *ent.User 객체
	//   - 실패 시 오류 (예: 존재하지 않음, 여러 개 존재함)
	SelectOneByUsername(ctx context.Context, username string) (*ent.User, error)

	// SelectManyByGroupGT는 그룹 값이 지정한 maxGroup보다 큰 사용자들을 조회합니다.
	//
	// 입력:
	//   - ctx: 요청 범위를 제어하기 위한 context.Context
	//   - maxGroup: 상위 그룹 필터 기준값
	//
	// 처리:
	//   - user.GroupGT(maxGroup) 조건으로 필터링
	//
	// 반환:
	//   - 그룹 값이 큰 사용자 목록 []*ent.User
	//   - 쿼리 실패 시 에러 반환
	SelectManyByGroupGT(ctx context.Context, maxGroup int) ([]*ent.User, error)

	// UpdateOne은 지정된 ID를 가진 사용자의 필드를 선택적으로 수정합니다.
	//
	// 입력:
	//   - ctx: 요청 범위를 제어하기 위한 context.Context
	//   - req: 업데이트할 사용자 정보가 포함된 UserDTO (ID는 필수)
	//
	// 처리:
	//   - 해당 ID가 존재하는지 먼저 확인
	//   - 존재하지 않으면 에러 반환
	//   - req에 포함된 필드들만 업데이트 (nil이 아닌 필드만 적용)
	//   - UpdatedAt은 현재 시각으로 갱신
	//
	// 반환:
	//   - 업데이트된 사용자 *ent.User 객체
	//   - 실패 시 에러 반환 (예: ID 미존재, DB 오류 등)
	UpdateOne(ctx context.Context, req UserDTO) (*ent.User, error)

	// DeleteOne은 지정된 ID를 가진 사용자 레코드를 삭제합니다.
	//
	// 입력:
	//   - ctx: 요청 범위를 제어하기 위한 context.Context
	//   - id: 삭제 대상 사용자 ID
	//
	// 처리:
	//   - 해당 ID가 존재하지 않으면 내부적으로 에러 발생
	//
	// 반환:
	//   - 삭제 성공 시 nil
	//   - 실패 시 에러 반환 (예: 존재하지 않는 ID, DB 오류 등)
	DeleteOne(ctx context.Context, id int) error
}

// UserDAO는 IUserDAO 인터페이스를 구현하는 구조체로,
// 내부적으로 ent 클라이언트를 사용하여 데이터베이스에 접근합니다.
type UserDAO struct {
	dbms *ent.Client
}

var userOnce sync.Once
var userInstance *UserDAO

func NewUserDAO() *UserDAO {
	userOnce.Do(func() {
		logger.Debug("User DAO instance")
		userInstance = &UserDAO{
			dbms: utils.GetEntClient(),
		}
	})
	return userInstance
}

// InsertOne은 새로운 사용자를 데이터베이스에 추가합니다.
//
// 필수 필드(name, username, password, group, is_use)가 누락된 경우 에러를 반환합니다.
func (dao *UserDAO) InsertOne(ctx context.Context, req UserDTO) (*ent.User, error) {
	logger.Debug(fmt.Sprintf("%+v", req))

	if req.Name == nil || req.Username == nil || req.Password == nil || req.Group == nil || req.IsUse == nil {
		return nil, fmt.Errorf("필수 입력값 누락: name, username, password, group, is_use는 필수입니다")
	}

	return dao.dbms.User.Create().
		SetName(*req.Name).
		SetUsername(*req.Username).
		SetPassword(*req.Password).
		SetGroup(*req.Group).
		SetIsUse(*req.IsUse).
		Save(ctx)
}

// SelectOneByUsername은 주어진 사용자명(username)에 해당하는 사용자 정보를 조회합니다.
//
// 사용자가 존재하지 않을 경우 에러를 반환합니다.
func (dao *UserDAO) SelectOneByUsername(ctx context.Context, username string) (*ent.User, error) {
	logger.Debug(fmt.Sprintf("Username: %s", username))
	return dao.dbms.User.Query().
		Where(user.Username(username)).
		Only(ctx)
}

// SelectAll은 모든 사용자를 ID 기준 내림차순으로 조회합니다.
func (dao *UserDAO) SelectAll(ctx context.Context) ([]*ent.User, error) {
	return dao.dbms.User.Query().
		Order(user.ByID(sql.OrderDesc())).
		All(ctx)
}

// SelectManyByGroupGT는 주어진 그룹 값(maxGroup)보다 높은 그룹을 가진 사용자들을 조회합니다.
//
// 결과는 ID 기준 내림차순으로 정렬됩니다.
func (dao *UserDAO) SelectManyByGroupGT(ctx context.Context, maxGroup int) ([]*ent.User, error) {
	return dao.dbms.User.Query().
		Where(user.GroupGT(maxGroup)).
		Order(user.ByID(sql.OrderDesc())).
		All(ctx)
}

// SelectOne은 ID로 특정 사용자를 조회합니다.
//
// 사용자가 존재하지 않을 경우 에러를 반환합니다.
func (dao *UserDAO) SelectOne(ctx context.Context, id int) (*ent.User, error) {
	logger.Debug(fmt.Sprintf(`{"user_id": %d}`, id))
	return dao.dbms.User.Get(ctx, id)
}

// UpdateOne은 주어진 ID에 해당하는 사용자의 정보를 업데이트합니다.
//
// 존재하지 않는 사용자일 경우 에러를 반환하며, req에 설정된 필드만 선택적으로 갱신합니다.
func (dao *UserDAO) UpdateOne(ctx context.Context, req UserDTO) (*ent.User, error) {
	logger.Debug(fmt.Sprintf("%+v", req))

	exists, err := dao.dbms.User.Query().Where(user.ID(req.ID)).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("데이터베이스 조회 실패: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("ID %d에 해당하는 유저를 찾을 수 없음", req.ID)
	}

	update := dao.dbms.User.UpdateOneID(req.ID)

	if req.Name != nil {
		update.SetName(*req.Name)
	}
	if req.Username != nil {
		update.SetUsername(*req.Username)
	}
	if req.Password != nil {
		update.SetPassword(*req.Password)
	}
	if req.Group != nil {
		update.SetGroup(*req.Group)
	}
	if req.IsUse != nil {
		update.SetIsUse(*req.IsUse)
	}

	update.SetUpdatedAt(time.Now())

	return update.Save(ctx)
}

// DeleteOne은 ID에 해당하는 사용자를 데이터베이스에서 삭제합니다.
//
// 존재하지 않는 경우에도 ent는 자체적으로 에러를 반환합니다.
func (dao *UserDAO) DeleteOne(ctx context.Context, id int) error {
	logger.Debug(fmt.Sprintf(`{"id": %d}`, id))
	return dao.dbms.User.
		DeleteOneID(id).
		Exec(ctx)
}

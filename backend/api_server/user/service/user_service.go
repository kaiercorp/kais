package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"api_server/logger"
	repo "api_server/user/repository"
	"api_server/utils"
)

type IUserService interface {
	// Create는 신규 사용자를 생성하는 서비스 함수입니다.
	//
	// 입력:
	//   - req: repo.UserDTO (생성할 사용자 정보)
	//
	// 동작:
	//   - 라이선스 사용자 수 체크
	//   - 비밀번호 암호화
	//   - 사용자 활성화 상태 기본값 설정
	//   - DB에 사용자 저장
	//   - 활성 사용자 수 증가 (라이선스 목업)
	//
	// 반환:
	//   - 성공: 생성된 사용자 정보 DTO
	//   - 실패: 라이선스 초과 또는 DB 삽입 오류에 대한 Report
	Create(req repo.UserDTO) (*repo.UserDTO, *logger.Report)

	// Edit는 사용자명을 기반으로 기존 사용자를 조회하고,
	// 해당 사용자의 정보를 수정하는 서비스 함수입니다.
	//
	// 입력:
	//   - req: repo.UserDTO (수정할 정보 포함, username 필수)
	//
	// 동작:
	//   - username을 기반으로 사용자 조회
	//   - 조회된 사용자의 ID를 설정한 후 update 실행
	//
	// 반환:
	//   - 성공: 수정된 사용자 DTO
	//   - 실패: 사용자 없음 또는 DB 업데이트 실패
	Edit(req repo.UserDTO) (*repo.UserDTO, *logger.Report)

	// Delete는 지정된 사용자 ID의 계정을 삭제합니다.
	//
	// 입력:
	//   - userID: 삭제할 사용자 ID
	//
	// 동작:
	//   - 사용자 DB 삭제
	//   - 라이선스 사용자 수 감소 (목업)
	//
	// 반환:
	//   - 성공: nil
	//   - 실패: DB 삭제 오류 Report
	Delete(userID int) *logger.Report

	// ReadByUsername는 사용자명을 기반으로 사용자를 조회합니다.
	//
	// 입력:
	//   - username: 조회할 사용자 이름
	//
	// 반환:
	//   - 성공: 사용자 DTO
	//   - 실패: 사용자 없음 또는 DB 오류
	ReadByUsername(username string) (*repo.UserDTO, *logger.Report)

	// ReadManyByGroupGT는 지정된 그룹 번호보다 큰 그룹에 속한 사용자들을 조회합니다.
	//
	// 입력:
	//   - group: 기준이 되는 그룹 번호
	//
	// 동작:
	//   - group > 입력값 조건을 만족하는 사용자 목록 조회
	//
	// 반환:
	//   - 성공: 사용자 DTO 슬라이스
	//   - 실패: DB 조회 오류
	ReadManyByGroupGT(group int) ([]*repo.UserDTO, *logger.Report)

	// ResetPassword는 지정된 사용자 ID의 비밀번호를 기본값으로 초기화합니다.
	//
	// 입력:
	//   - userID: 비밀번호를 초기화할 사용자 ID
	//
	// 동작:
	//   - 사용자 조회
	//   - 비밀번호를 "password1234"로 암호화 후 설정
	//   - 사용자 정보 업데이트
	//
	// 반환:
	//   - 성공: 비밀번호 초기화된 사용자 DTO
	//   - 실패: 사용자 조회 또는 업데이트 오류
	ResetPassword(userID int) (*repo.UserDTO, *logger.Report)

	// CheckPassword는 사용자 DTO의 저장된 비밀번호와 주어진 평문 비밀번호가 일치하는지 확인합니다.
	//
	// 입력:
	//   - req: repo.UserDTO 포인터 (비교 대상 사용자)
	//   - password: 확인할 평문 비밀번호
	//
	// 반환:
	//   - true: 비밀번호 일치
	//   - false: 비밀번호 불일치
	//
	// 주의:
	//   - 비밀번호는 내부적으로 암호화 후 비교합니다.
	CheckPassword(req *repo.UserDTO, password string) bool

	// Activate는 사용자의 활성화 상태를 변경하는 서비스 함수입니다.
	//
	// 입력:
	//   - userID: 대상 사용자 ID
	//   - active: 활성화 여부 (true: 활성화, false: 비활성화)
	//
	// 동작:
	//   - DB에서 사용자 조회
	//   - 현재 상태와 동일하면 무시하고 반환
	//   - 상태 변경 시, 라이선스 사용자 수 조정 (목업)
	//   - 변경된 활성화 상태로 사용자 업데이트
	//
	// 반환:
	//   - 성공: 상태 변경된 사용자 DTO
	//   - 실패: 사용자 조회 실패 또는 업데이트 실패
	Activate(userID int, active bool) (*repo.UserDTO, *logger.Report)

	// ChangePassword는 사용자의 현재 비밀번호를 검증하고, 새 비밀번호로 변경합니다.
	//
	// 입력:
	//   - username: 비밀번호를 변경할 사용자 이름
	//   - oldPassword: 현재 비밀번호 (검증용)
	//   - newPassword: 새로 설정할 비밀번호
	//
	// 동작:
	//   - 사용자 조회
	//   - 현재 비밀번호 검증
	//   - 비밀번호 암호화 후 저장
	//
	// 반환:
	//   - 성공: 비밀번호가 변경된 사용자 DTO
	//   - 실패: 사용자 없음, 비밀번호 불일치, DB 오류 등
	ChangePassword(username string, oldPassword string, newPassword string) (*repo.UserDTO, *logger.Report)
}

type UserService struct {
	ctx context.Context
	dao repo.IUserDAO
}

var userInstance *UserService
var userOnce sync.Once

func NewUserService(dao repo.IUserDAO) *UserService {
	userOnce.Do(func() {
		logger.Debug("UserService instance created")
		userInstance = &UserService{
			ctx: context.Background(),
			dao: dao,
		}
	})
	return userInstance
}

// Create는 신규 사용자를 생성하는 서비스 함수입니다.
func (svc *UserService) Create(req repo.UserDTO) (*repo.UserDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("CreateUser: %+v", req))

	// 비밀번호 암호화
	s := utils.CreateSecurity()
	password, _ := s.Encrypt("password1234")
	isUse := true

	user := repo.UserDTO{
		Username:  req.Username,
		Name:      req.Name,
		Password:  &password,
		Group:     req.Group,
		IsUse:     &isUse, // 기본적으로 활성화 상태
		CreatedAt: time.Now(),
	}

	ent, err := svc.dao.InsertOne(svc.ctx, user)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_INSERT, err)
	}

	return repo.ConvertUserEntToDTO(ent), nil
}

// update는 내부적으로 사용자 정보를 DB에 반영하는 유틸성 함수입니다.
//
// 입력:
//   - req: repo.UserDTO (갱신할 정보 포함)
//
// 반환:
//   - 성공: 업데이트된 사용자 DTO
//   - 실패: DB 업데이트 실패 Report
//
// 주의:
//   - 이 함수는 외부에서 직접 호출되지 않으며, Edit/Activate 등에서 내부적으로 사용됩니다.
func (svc *UserService) update(req repo.UserDTO) (*repo.UserDTO, *logger.Report) {
	if req.Password != nil {
		s := utils.CreateSecurity()
		encryptedPassword, _ := s.Encrypt(*req.Password)
		req.Password = &encryptedPassword
	}

	ent, err := svc.dao.UpdateOne(svc.ctx, req)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	return repo.ConvertUserEntToDTO(ent), nil
}

// Edit는 사용자명을 기반으로 기존 사용자를 조회하고,
// 해당 사용자의 정보를 수정하는 서비스 함수입니다.
func (svc *UserService) Edit(req repo.UserDTO) (*repo.UserDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("EditUser: %+v", req))
	user, r := svc.ReadByUsername(*req.Username)
	if r != nil {
		return nil, r
	}
	req.ID = user.ID

	return svc.update(req)
}

// Activate는 사용자의 활성화 상태를 변경하는 서비스 함수입니다.
func (svc *UserService) Activate(userID int, active bool) (*repo.UserDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("ActivateUser: ID=%d", userID))
	//current activation check
	user, err := svc.dao.SelectOne(svc.ctx, userID)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	if user.IsUse == active {
		return repo.ConvertUserEntToDTO(user), nil
	}

	userDto := repo.UserDTO{
		ID:    userID,
		IsUse: &active,
	}
	return svc.update(userDto)
}

// ResetPassword는 지정된 사용자 ID의 비밀번호를 기본값으로 초기화합니다.
func (svc *UserService) ResetPassword(userID int) (*repo.UserDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("ResetPassword: ID=%d", userID))
	user, err := svc.dao.SelectOne(svc.ctx, userID)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}
	user.Password = "password1234"

	return svc.update(*repo.ConvertUserEntToDTO(user))
}

// Delete는 지정된 사용자 ID의 계정을 삭제합니다.
func (svc *UserService) Delete(userID int) *logger.Report {
	logger.Debug(fmt.Sprintf("DeleteUser: ID=%d", userID))

	err := svc.dao.DeleteOne(svc.ctx, userID)
	if err != nil {
		return logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}

	return nil
}

// CheckPassword는 사용자 DTO의 저장된 비밀번호와 주어진 평문 비밀번호가 일치하는지 확인합니다.
func (svc *UserService) CheckPassword(req *repo.UserDTO, password string) bool {
	s := utils.CreateSecurity()
	p, _ := s.Encrypt(password)
	matched := (*req.Password == p)

	return matched
}

// ChangePassword는 사용자의 현재 비밀번호를 검증하고, 새 비밀번호로 변경합니다.
func (svc *UserService) ChangePassword(username string, oldPassword string, newPassword string) (*repo.UserDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("ChangePassword: ID=%s", username))
	user, r := svc.ReadByUsername(username)
	if r != nil {
		return nil, r
	}
	if !svc.CheckPassword(user, oldPassword) {
		return nil, logger.CreateReport(&logger.CODE_LOGIN_PARAMS, fmt.Errorf("incorrected old password"))
	}
	userDto := repo.UserDTO{
		ID:       user.ID,
		Password: &newPassword,
	}
	return svc.update(userDto)
}

// ReadByUsername는 사용자명을 기반으로 사용자를 조회합니다.
func (svc *UserService) ReadByUsername(username string) (*repo.UserDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("Username: %s", username))
	user, err := svc.dao.SelectOneByUsername(svc.ctx, username)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}
	return repo.ConvertUserEntToDTO(user), nil
}

// ReadManyByGroupGT는 지정된 그룹 번호보다 큰 그룹에 속한 사용자들을 조회합니다.
func (svc *UserService) ReadManyByGroupGT(group int) ([]*repo.UserDTO, *logger.Report) {
	logger.Debug(fmt.Sprintf("Group: %d", group))
	users, err := svc.dao.SelectManyByGroupGT(svc.ctx, group)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}
	userList := make([]*repo.UserDTO, len(users))

	for i, ent := range users {
		userList[i] = repo.ConvertUserEntToDTO(ent)
	}

	return userList, nil
}

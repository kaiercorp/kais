package csvformat

import (
	"api_server/ent"
	repo "api_server/project/repository"
	"api_server/utils"
	"strconv"
)

// ProjectCSVFormat 구조체는 Project 객체를 CSV 형식으로 변환하거나 CSV 데이터를 Project 객체로 변환하는 데 사용됩니다.
type ProjectCSVFormat struct{}

// GetHeader는 Project 객체의 CSV 헤더를 반환하는 함수입니다.
//
// 반환 값:
//   - []string: CSV 헤더(컬럼 이름) 목록
func (f *ProjectCSVFormat) GetHeader() []string {
	return []string{
		"id", "title", "description", "is_use", "created_at", "updated_at",
	}
}

// ConvertToRecord는 주어진 ent.Project 객체를 CSV 레코드 형식으로 변환하는 함수입니다.
//
// 매개변수:
//   - project: 변환할 ent.Project 객체
//
// 반환 값:
//   - []string: ent.Project 객체를 CSV 레코드 형식으로 변환한 결과
func (f *ProjectCSVFormat) ConvertToRecord(project *ent.Project) []string {
	return []string{
		strconv.Itoa(project.ID),             // ID
		project.Title,                        // Title
		project.Description,                  // Description
		strconv.FormatBool(project.IsUse),    // IsUse (bool 값을 문자열로 변환)
		utils.FormatTime(&project.CreatedAt), // CreatedAt (시간 형식 변환)
		utils.FormatTime(&project.UpdatedAt), // UpdatedAt (시간 형식 변환)
	}
}

// ParseRecord는 CSV 레코드를 ProjectDTO 객체로 변환하는 함수입니다.
//
// 매개변수:
//   - record: 변환할 CSV 레코드
//
// 반환 값:
//   - *repo.ProjectDTO: 변환된 ProjectDTO 객체
//   - error: 변환 중 발생할 수 있는 오류
func (f *ProjectCSVFormat) ParseRecord(record []string) (*repo.ProjectDTO, error) {
	// CSV에서 ID 값을 정수로 변환
	oldID, _ := strconv.Atoi(record[0])

	// CSV에서 Title, Description 값을 추출하여 ProjectDTO 객체 생성
	project := &repo.ProjectDTO{
		ID:          oldID,
		Title:       record[1],
		Description: record[2],
		//IsUse:       utils.ParseBool(record[3]), // Assumes ParseBool is a utility function that converts the string to a bool
	}

	// ProjectDTO 객체 반환
	return project, nil
}

package csvformat

import (
	"api_server/ent"
	repo "api_server/task/repository"
	"api_server/utils"
	"strconv"
)

// TaskCSVFormat 구조체는 Task 객체를 CSV 형식으로 변환하거나 CSV 데이터를 Task 객체로 변환하는 데 사용됩니다.
type TaskCSVFormat struct{}

// GetHeader는 Task 객체의 CSV 헤더를 반환하는 함수입니다.
//
// 반환 값:
//   - []string: CSV 헤더(컬럼 이름) 목록
func (f *TaskCSVFormat) GetHeader() []string {
	return []string{
		"id", "dataset_id", "title", "description", "engine_type",
		"target_metric", "params", "created_at", "updated_at", "project_id",
	}
}

// ConvertToRecord는 주어진 ent.Task 객체를 CSV 레코드 형식으로 변환하는 함수입니다.
//
// 매개변수:
//   - task: 변환할 ent.Task 객체
//
// 반환 값:
//   - []string: ent.Task 객체를 CSV 레코드 형식으로 변환한 결과
func (f *TaskCSVFormat) ConvertToRecord(task *ent.Task) []string {
	return []string{
		strconv.Itoa(task.ID),             // ID
		strconv.Itoa(task.DatasetID),      // DatasetID
		task.Title,                        // Title
		task.Description,                  // Description
		task.EngineType,                   // EngineType
		task.TargetMetric,                 // TargetMetric
		utils.JSONString(task.Params),     // Params (JSON 형식으로 변환)
		utils.FormatTime(&task.CreatedAt), // CreatedAt (시간 형식 변환)
		utils.FormatTime(&task.UpdatedAt), // UpdatedAt (시간 형식 변환)
		strconv.Itoa(task.ProjectID),      // ProjectID
	}
}

// ParseRecord는 CSV 레코드를 TaskDTO 객체로 변환하는 함수입니다.
//
// 매개변수:
//   - record: 변환할 CSV 레코드
//
// 반환 값:
//   - *repo.TaskDTO: 변환된 TaskDTO 객체
//   - error: 변환 중 발생할 수 있는 오류
func (f *TaskCSVFormat) ParseRecord(record []string) (*repo.TaskDTO, error) {
	// CSV에서 ID 값을 정수로 변환
	oldID, _ := strconv.Atoi(record[0])
	// DatasetID는 CSV에 명시되지 않으므로 기본값 0으로 설정
	datasetID := 0
	// Params는 JSON 형식이므로 이를 파싱
	params, _ := utils.ParseJSONB(record[6])
	// ProjectID 값을 정수로 변환
	oldProjectID, _ := strconv.Atoi(record[9])

	// 변환된 값을 바탕으로 TaskDTO 객체 생성
	task := &repo.TaskDTO{
		ID:           oldID,
		ProjectID:    &oldProjectID,
		DatasetID:    &datasetID,
		Title:        record[2],
		Description:  record[3],
		EngineType:   record[4],
		TargetMetric: record[5],
		Params:       params,
	}

	// TaskDTO 객체 반환
	return task, nil
}

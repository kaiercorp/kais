package csvformat

import (
	"api_server/ent"
	repo "api_server/task/repository"
	"api_server/utils"
	"strconv"
)

// ModelingCSVFormat 구조체는 Modeling 객체를 CSV 형식으로 변환하거나 CSV 데이터를 Modeling 객체로 변환하는 데 사용됩니다.
type ModelingCSVFormat struct{}

// GetHeader는 Modeling 객체의 CSV 헤더를 반환하는 함수입니다.
//
// 반환 값:
//   - []string: CSV 헤더(컬럼 이름) 목록
func (f *ModelingCSVFormat) GetHeader() []string {
	return []string{
		"id", "local_id", "parent_id", "parent_local_id", "dataset_id", "params",
		"dataset_stat", "modeling_type", "modeling_step", "performance", "progress",
		"created_at", "updated_at", "task_id",
	}
}

// ConvertToRecord는 주어진 ent.Modeling 객체를 CSV 레코드 형식으로 변환하는 함수입니다.
//
// 매개변수:
//   - modeling: 변환할 ent.Modeling 객체
//
// 반환 값:
//   - []string: ent.Modeling 객체를 CSV 레코드 형식으로 변환한 결과
func (f *ModelingCSVFormat) ConvertToRecord(modeling *ent.Modeling) []string {
	return []string{
		strconv.Itoa(modeling.ID),                          // ID
		strconv.Itoa(modeling.LocalID),                     // LocalID
		strconv.Itoa(modeling.ParentID),                    // ParentID
		strconv.Itoa(modeling.ParentLocalID),               // ParentLocalID
		strconv.Itoa(modeling.DatasetID),                   // DatasetID
		utils.JSONString(modeling.Params),                  // Params (JSON 형식으로 변환)
		utils.JSONString(modeling.DatasetStat),             // DatasetStat (JSON 형식으로 변환)
		modeling.ModelingType,                              // ModelingType
		modeling.ModelingStep,                              // ModelingStep
		utils.JSONString(modeling.Performance),             // Performance (JSON 형식으로 변환)
		strconv.FormatFloat(modeling.Progress, 'f', 2, 64), // Progress (소수점 2자리로 포맷)
		utils.FormatTime(&modeling.CreatedAt),              // CreatedAt (시간 형식 변환)
		utils.FormatTime(&modeling.UpdatedAt),              // UpdatedAt (시간 형식 변환)
		strconv.Itoa(modeling.TaskID),                      // TaskID
	}
}

// ParseRecord는 CSV 레코드를 ModelingDTO 객체로 변환하는 함수입니다.
//
// 매개변수:
//   - record: 변환할 CSV 레코드
//
// 반환 값:
//   - *repo.ModelingDTO: 변환된 ModelingDTO 객체
//   - error: 변환 중 발생할 수 있는 오류
func (f *ModelingCSVFormat) ParseRecord(record []string) (*repo.ModelingDTO, error) {
	// CSV에서 ID, TaskID, ParentID, ParentLocalID, DatasetID 값을 정수로 변환
	oldID, _ := strconv.Atoi(record[0])
	oldTaskID, _ := strconv.Atoi(record[13])
	parentID, _ := strconv.Atoi(record[2])
	parentLocalID, _ := strconv.Atoi(record[3])
	datasetID, _ := strconv.Atoi(record[4])

	// CSV에서 Params, DatasetStat, Performance 값을 JSON으로 파싱
	params, _ := utils.ParseJSONB(record[5])
	datasetStat, _ := utils.ParseJSONB(record[6])
	performance, _ := utils.ParseJSONB(record[9])

	// CSV에서 Progress 값을 실수로 변환
	progress, _ := strconv.ParseFloat(record[10], 64)

	// 변환된 값을 바탕으로 ModelingDTO 객체 생성
	modeling := &repo.ModelingDTO{
		ID:            oldID,
		ParentID:      parentID,
		ParentLocalID: parentLocalID,
		DatasetID:     datasetID,
		Params:        params,
		DatasetStat:   datasetStat,
		ModelingType:  record[7],
		ModelingStep:  record[8],
		Performance:   performance,
		Progress:      progress,
		TaskID:        oldTaskID,
	}

	// ModelingDTO 객체 반환
	return modeling, nil
}

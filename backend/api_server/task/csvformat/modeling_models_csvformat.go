package csvformat

import (
	"api_server/ent"
	repo "api_server/task/repository"
	"api_server/utils"
	"strconv"
)

// ModelingModelCSVFormat 구조체는 ModelingModel 객체를 CSV 형식으로 변환하거나 CSV 데이터를 ModelingModel 객체로 변환하는 데 사용됩니다.
type ModelingModelCSVFormat struct{}

// GetHeader는 ModelingModel 객체의 CSV 헤더를 반환하는 함수입니다.
//
// 반환 값:
//   - []string: CSV 헤더(컬럼 이름) 목록
func (f *ModelingModelCSVFormat) GetHeader() []string {
	return []string{
		"id", "data_type", "data", "created_at", "modeling_id",
	}
}

// ConvertToRecord는 주어진 ent.ModelingModels 객체를 CSV 레코드 형식으로 변환하는 함수입니다.
//
// 매개변수:
//   - model: 변환할 ent.ModelingModels 객체
//
// 반환 값:
//   - []string: ent.ModelingModels 객체를 CSV 레코드 형식으로 변환한 결과
func (f *ModelingModelCSVFormat) ConvertToRecord(model *ent.ModelingModels) []string {
	return []string{
		strconv.Itoa(model.ID),             // ID
		model.DataType,                     // DataType
		model.Data,                         // Data
		utils.FormatTime(&model.CreatedAt), // CreatedAt (시간 형식 변환)
		strconv.Itoa(model.ModelingID),     // ModelingID
	}
}

// ParseRecord는 CSV 레코드를 ModelingModels 객체로 변환하는 함수입니다.
//
// 매개변수:
//   - record: 변환할 CSV 레코드
//
// 반환 값:
//   - *repo.ModelingModels: 변환된 ModelingModels 객체
//   - error: 변환 중 발생할 수 있는 오류
func (f *ModelingModelCSVFormat) ParseRecord(record []string) (*repo.ModelingModels, error) {
	// CSV에서 ID 및 ModelingID 값을 정수로 변환
	oldID, _ := strconv.Atoi(record[0])
	oldModelingID, _ := strconv.Atoi(record[4])

	// 변환된 값을 바탕으로 ModelingModelDTO 객체 생성
	modelingModel := &repo.ModelingModels{
		ID:         oldID,
		ModelingID: oldModelingID,
		DataType:   record[1],
		Data:       record[2],
	}

	// ModelingModelDTO 객체 반환
	return modelingModel, nil
}

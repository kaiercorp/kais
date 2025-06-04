package csvformat

import (
	"api_server/ent"
	repo "api_server/task/repository"
	"api_server/utils"
	"strconv"
)

// ModelingDetailCSVFormat 구조체는 ModelingDetail 객체를 CSV 형식으로 변환하거나 CSV 데이터를 ModelingDetail 객체로 변환하는 데 사용됩니다.
type ModelingDetailCSVFormat struct{}

// GetHeader는 ModelingDetail 객체의 CSV 헤더를 반환하는 함수입니다.
//
// 반환 값:
//   - []string: CSV 헤더(컬럼 이름) 목록
func (f *ModelingDetailCSVFormat) GetHeader() []string {
	return []string{
		"id", "model", "data_type", "data", "created_at", "modeling_id",
	}
}

// ConvertToRecord는 주어진 ent.ModelingDetails 객체를 CSV 레코드 형식으로 변환하는 함수입니다.
//
// 매개변수:
//   - detail: 변환할 ent.ModelingDetails 객체
//
// 반환 값:
//   - []string: ent.ModelingDetails 객체를 CSV 레코드 형식으로 변환한 결과
func (f *ModelingDetailCSVFormat) ConvertToRecord(detail *ent.ModelingDetails) []string {
	return []string{
		strconv.Itoa(detail.ID),             // ID
		detail.Model,                        // Model
		detail.DataType,                     // DataType
		utils.JSONString(detail.Data),       // Data (JSON 형식으로 변환)
		utils.FormatTime(&detail.CreatedAt), // CreatedAt (시간 형식 변환)
		strconv.Itoa(detail.ModelingID),     // ModelingID
	}
}

// ParseRecord는 CSV 레코드를 ModelingDetailDTO 객체로 변환하는 함수입니다.
//
// 매개변수:
//   - record: 변환할 CSV 레코드
//
// 반환 값:
//   - *repo.ModelingDetailDTO: 변환된 ModelingDetailDTO 객체
//   - error: 변환 중 발생할 수 있는 오류
func (f *ModelingDetailCSVFormat) ParseRecord(record []string) (*repo.ModelingDetailDTO, error) {
	// CSV에서 ID 및 ModelingID 값을 정수로 변환
	oldID, _ := strconv.Atoi(record[0])
	oldModelingID, _ := strconv.Atoi(record[5])

	// CSV에서 Data 값을 JSON으로 파싱
	data, _ := utils.ParseJSONB(record[3])

	// 변환된 값을 바탕으로 ModelingDetailDTO 객체 생성
	modelingDetail := &repo.ModelingDetailDTO{
		ID:         oldID,
		Model:      record[1],
		DataType:   record[2],
		Data:       data,
		ModelingID: oldModelingID,
	}

	// ModelingDetailDTO 객체 반환
	return modelingDetail, nil
}

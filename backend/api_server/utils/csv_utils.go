package utils

import (
	"api_server/logger"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ExportToCSV는 주어진 데이터와 헤더를 CSV 파일로 내보내는 함수입니다.
// 파일을 생성하고, CSV 형식으로 데이터를 기록합니다.
//
// 매개변수:
//   - filename: 저장할 파일의 이름
//   - header: CSV 파일의 첫 번째 행으로, 컬럼 헤더를 나타냅니다.
//   - records: CSV 파일에 기록할 데이터 배열
//
// 반환 값:
//   - filename: 생성된 파일 이름
//   - 오류 보고서. 오류가 있을 경우 logger.Report가 반환됩니다.
func ExportToCSV(filename string, header []string, records [][]string) (string, *logger.Report) {
	file, err := os.Create(filename)
	if err != nil {
		return "", logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 헤더 작성
	if err := writer.Write(header); err != nil {
		return "", logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	// 데이터 기록
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return "", logger.CreateReport(&logger.CODE_DB_UPDATE, err)
		}
	}

	return filename, nil
}

// FormatTime은 주어진 시간을 RFC3339 형식의 문자열로 변환하는 함수입니다.
// 만약 입력값이 nil이면 빈 문자열을 반환합니다.
//
// 매개변수:
//   - t: 변환할 시간 (time.Time 타입)
//
// 반환 값:
//   - RFC3339 형식으로 변환된 시간 문자열
func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

// JSONString은 주어진 값을 JSON 문자열로 변환하는 헬퍼 함수입니다.
//
// 매개변수:
//   - v: 변환할 값
//
// 반환 값:
//   - JSON 문자열
func JSONString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// ParseJSON은 JSON 형식의 문자열을 맵(map[string]interface{})로 변환하는 헬퍼 함수입니다.
// 만약 값이 빈 문자열이면 nil을 반환합니다.
//
// 매개변수:
//   - value: 변환할 JSON 문자열
//
// 반환 값:
//   - 변환된 맵 (변환 실패 시 nil 반환)
func ParseJSON(value string) map[string]interface{} {
	if value == "" {
		return nil
	}

	var userParams map[string]interface{}
	err := json.Unmarshal([]byte(value), &userParams)
	if err != nil {
		return nil
	}
	return userParams
}

// ParseTime은 주어진 RFC3339 형식의 문자열을 time.Time 객체로 변환하는 함수입니다.
// 만약 값이 빈 문자열이면 nil을 반환합니다.
//
// 매개변수:
//   - s: 변환할 시간 문자열
//
// 반환 값:
//   - 변환된 time.Time 객체 또는 오류
func ParseTime(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil // 빈 문자열일 경우 nil을 반환
	}

	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err // RFC3339 형식이 아닐 경우 오류 반환
	}

	return &parsedTime, nil
}

// ParseJSONB는 postgresql JSONB 형식의 문자열을 []string으로 변환하는 헬퍼 함수입니다.
// 만약 값이 빈 문자열이면 nil을 반환합니다.
//
// 매개변수:
//   - value: 변환할 JSON 문자열
//
// 반환 값:
//   - 변환된 문자열 슬라이스 또는 오류
func ParseJSONB(value string) ([]string, error) {
	if value == "" {
		return nil, nil // 빈 문자열일 경우 nil을 반환
	}

	// JSON 문자열을 []string으로 파싱
	var params []string
	err := json.Unmarshal([]byte(value), &params)
	if err != nil {
		return nil, fmt.Errorf("params 필드 파싱 오류: %v", err)
	}

	return params, nil
}

// CSVImportFunc는 CSV 레코드를 해당 타입으로 변환하는 함수 타입입니다.
//
// 매개변수:
// - record: CSV에서 읽은 한 줄의 데이터
//
// 반환 값:
// - 변환된 객체와 기존 ID, 오류
type CSVImportFunc[T any] func(record []string) (*T, int, error)

// ImportFromCSV는 CSV 파일을 읽고, 각 레코드를 지정된 importFunc로 변환한 뒤,
// insertFunc를 호출하여 데이터베이스에 삽입하는 공통 함수입니다.
// 실패한 레코드는 별도로 기록합니다. 각 service 들의 importCSV 부분 참조.
//
// 매개변수:
//   - filename: 읽을 CSV 파일의 경로
//   - importFunc: 각 레코드를 객체로 변환하는 함수
//   - insertFunc: 변환된 객체를 데이터베이스에 삽입하는 함수
//
// 반환 값:
//   - 성공적인 ID 변환 맵과 오류 보고서
//     오류가 있을 경우 logger.Report가 반환됩니다.
func ImportFromCSV[T any](
	filename string,
	importFunc CSVImportFunc[T],
	insertFunc func(*T) (int, error),
) (map[int]int, *logger.Report) {

	logger.Debug(fmt.Sprintf("Importing data from %s", filename))

	// CSV 파일 열기
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_INSERT, err)
	}
	defer csvFile.Close()

	// CSV 리더 설정
	csvReader := csv.NewReader(csvFile)

	// 첫 번째 행(헤더) 건너뛰기
	_, err = csvReader.Read()
	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_INSERT, err)
	}

	var failedRecords []string
	newIdMap := make(map[int]int)

	// CSV 데이터 한 줄씩 읽기
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err.Error() != "EOF" {
				return nil, logger.CreateReport(&logger.CODE_DB_INSERT, err)
			}
			break
		}

		// 데이터 변환
		entity, oldID, err := importFunc(record)
		if err != nil {
			failedRecords = append(failedRecords, fmt.Sprintf("Failed to parse record: %v", record))
			continue
		}

		// 데이터베이스 삽입
		newID, err := insertFunc(entity)
		if err != nil {
			failedRecords = append(failedRecords, fmt.Sprintf("Failed to insert: %v", record))
			continue
		}

		newIdMap[oldID] = newID
	}

	// 실패한 레코드 처리
	if len(failedRecords) > 0 {
		return nil, logger.CreateReport(&logger.CODE_DB_INSERT, fmt.Errorf("일부 데이터 import에 실패했습니다"))
	}

	return newIdMap, nil
}

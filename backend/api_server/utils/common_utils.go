package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"time"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

type ObjectStruct struct {
	Index string
	Value float64
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func GenerateRandomString(n int) string {
	if n <= 0 {
		n = 16
	}

	// Set String Bytes Size
	b := make([]byte, n)

	// Seed Init
	rand.NewSource(time.Now().UnixNano())

	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits
		remain--
	}

	return string(b[:])
}

func AnsiToUniString(src string) string {
	return string(AnsiToUniBytes([]byte(src)))
}

func AnsiToUniBytes(src []byte) []byte {
	got, _, _ := transform.String(korean.EUCKR.NewDecoder(), string(src))
	return []byte(got)
}

func GetLicenseType(trialType string) string {
	licenseType := ""
	switch trialType {
	case JOB_TYPE_VISION_CLS_SL:
		licenseType = VIS_CLS_SL_TGUI
	case JOB_TYPE_VISION_CLS_ML:
		licenseType = VIS_CLS_ML_TGUI
	case JOB_TYPE_VISION_SEG:
		licenseType = VIS_SEG_TGUI
	case JOB_TYPE_VISION_OD:
		licenseType = VIS_OD_TGUI
	case JOB_TYPE_VISION_OCR:
		licenseType = VIS_OCR_TGUI
	case JOB_TYPE_VISION_AD:
		licenseType = VIS_AD_TGUI
	case JOB_TYPE_TABLE_CLS:
		licenseType = TAB_CLS_TGUI
	case JOB_TYPE_TABLE_REG:
		licenseType = TAB_REG_TGUI
	case JOB_TYPE_TS_AD:
		licenseType = TS_AD_TGUI
	case JOB_TYPE_TS_DF:
		licenseType = TS_DF_TGUI
	default:
		licenseType = ""
	}

	return licenseType
}

func GetMetricList(engineType string) []string {
	switch engineType {
	case JOB_TYPE_VISION_CLS_SL:
		return []string{"wa", "uwa", "f1", "recall", "precision"}
	case JOB_TYPE_VISION_CLS_ML:
		return []string{"image_accuracy", "image_precision", "image_recall", "image_f1_score", "label_accuracy", "label_precision", "label_recall", "label_f1_score"}
	case JOB_TYPE_VISION_AD:
		return []string{"wa", "uwa", "f1", "recall", "precision", "auroc", "prauc"}
	case JOB_TYPE_TABLE_CLS:
		return []string{"wa", "uwa", "f1", "recall", "precision", "aucroc"}
	case JOB_TYPE_TABLE_REG:
		return []string{"mse", "rmse", "mae", "r2", "xvar"}
	case JOB_TYPE_TS_AD:
		return []string{"mse", "rmse", "mae"}
	}

	return []string{}
}

func GetMaxWeightMetric(metrics map[string]interface{}) string {
	var targetMetric string
	var maxWeight float64

	for metric, weight := range metrics {
		if weight.(float64) > maxWeight {
			targetMetric = metric
			maxWeight = weight.(float64)
		}
	}

	return targetMetric
}

func SetTargetMetricMap(engineType string, targetMetric string) map[string]interface{} {
	switch engineType {
	case JOB_TYPE_VISION_CLS_ML:
		return genTargetMetricVCLSML(targetMetric)
	}

	return nil
}

func genTargetMetricVCLSML(targetMetric string) map[string]interface{} {
	targetMetricMap := make(map[string]interface{})
	targetMetricMap["image_accuracy"] = 0
	targetMetricMap["image_recall"] = 0
	targetMetricMap["image_precision"] = 0
	targetMetricMap["image_f1_score"] = 0
	targetMetricMap["label_accuracy"] = 0
	targetMetricMap["label_recall"] = 0
	targetMetricMap["label_precision"] = 0
	targetMetricMap["label_f1_score"] = 0
	targetMetricMap[targetMetric] = 1

	return targetMetricMap
}

func StringToMap(source string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if err := json.Unmarshal([]byte(source), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func StringToMapSlice(source string) (map[string][]interface{}, error) {
	result := make(map[string][]interface{})
	if err := json.Unmarshal([]byte(source), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func StringToMapMap(source string) (map[string]map[string]interface{}, error) {
	result := make(map[string]map[string]interface{})
	if err := json.Unmarshal([]byte(source), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func StringToReversMap(source string) (map[string]string, error) {
	result := make(map[string]string)

	origin := make(map[string]interface{})
	if err := json.Unmarshal([]byte(source), &origin); err != nil {
		return nil, err
	}

	for k, v := range origin {
		_v := fmt.Sprintf("%v", v)
		result[_v] = k
	}

	return result, nil
}

func SortObjectToMap(source map[string]interface{}) ([]interface{}, error) {
	sorted_keys := make([]float64, 0, len(source))
	for k := 0; k < len(source); k++ {
		strKey := fmt.Sprintf("%d", k)

		if v, ok := source[strKey].(float64); ok {
			sorted_keys = append(sorted_keys, v)
		} else {
			return nil, errors.New("value is not a float64")
		}
	}
	var values []interface{}
	for _, val := range sorted_keys {
		values = append(values, val)
	}

	return values, nil
}

func SortByValue(source map[string]interface{}) ([]ObjectStruct, error) {
	var sortedValues []ObjectStruct

	for key, value := range source {
		if v, ok := value.(float64); ok {
			sortedValues = append(sortedValues, ObjectStruct{key, v})
		}
	}

	sort.Slice(sortedValues, func(i, j int) bool {
		return sortedValues[i].Value < sortedValues[j].Value
	})

	return sortedValues, nil
}

func ImagePathToURL(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return string(' '), errors.New("failed to open image path")
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return string(' '), errors.New("failed to read image file")
	}

	base64Image := base64.StdEncoding.EncodeToString(imageBytes)
	fileType := http.DetectContentType(imageBytes)
	dataURL := fmt.Sprintf("data:%s;base64,%s", fileType, base64Image)

	return dataURL, nil
}

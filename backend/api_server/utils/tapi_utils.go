package utils

func ConvertVCLSResult(result map[string]interface{}) map[string]interface{} {
	rawLabels := result["prediction"].([]interface{})
	probsRaw := result["probabilities"].([]interface{})
	inferenceTime := result["inference_time"].(float64)

	var labels []string
	for _, v := range rawLabels {
		if labelStr, ok := v.(string); ok {
			labels = append(labels, labelStr)
		}
	}

	var probFloats []float64
	if len(probsRaw) > 0 {
		innerList := probsRaw[0].([]interface{})
		for _, pairRaw := range innerList {
			pair := pairRaw.([]interface{})
			if len(pair) >= 2 {
				if prob, ok := pair[1].(float64); ok {
					probFloats = append(probFloats, prob)
				}
			}
		}
	}

	data := map[string]interface{}{
		"labels":         labels,
		"probabilities":  probFloats,
		"inference_time": inferenceTime,
	}

	return data
}

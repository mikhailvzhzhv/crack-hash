package util

import (
	"encoding/json"

	"github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
)

func StructToJSON(obj any) []byte {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		FailOnError(err, "Failed to marshal struct to JSON")
		return nil
	}

	return jsonData
}

func JSONToTask(data []byte) *models.Task {
	var task *models.Task

	err := json.Unmarshal(data, task)
	if err != nil {
		FailOnError(err, "Failed to unmarshal JSON to Task")
		return nil
	}

	return task
}

func JSONToResult(data []byte) *models.Result {
	var result *models.Result

	err := json.Unmarshal(data, result)
	if err != nil {
		FailOnError(err, "Failed to unmarshal JSON to Result")
		return nil
	}

	return result
}

package models

import "github.com/google/uuid"

type Result struct {
	RequestID        uuid.UUID `json:"request_id"`
	TaskID           uuid.UUID `json:"task_id"`
	PartNumber       int       `json:"part_number"`
	PartCount        int       `json:"part_count"`
	Words            []string  `json:"words"`
	AvgExecutionTime int       `json:"avg_execution_time"`
}

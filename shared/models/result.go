package models

import "github.com/google/uuid"

type Result struct {
	RequestID  uuid.UUID `json:"request_id"`
	PartNumber int       `json:"part_number"`
	PartCount  int       `json:"part_count"`
	Words      []string  `json:"words"`
}

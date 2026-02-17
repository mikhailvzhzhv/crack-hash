package models

import "github.com/google/uuid"

type Result struct {
	RequestID  uuid.UUID
	PartNumber int
	PartCount  int
	Words      []string
}

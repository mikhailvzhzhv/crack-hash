package models

import "github.com/google/uuid"

type Task struct {
	TaskID     uuid.UUID
	PartNumber int
	PartCount  int
	Alphabet   string
	WordLen    int
	TargetHash string
}

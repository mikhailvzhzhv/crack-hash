package models

import "github.com/google/uuid"

type Algorithm string

const (
	MD5 = Algorithm("MD5")
)

type Task struct {
	RequestID  uuid.UUID `json:"request_id"`
	TaskID     uuid.UUID `json:"task_id"`
	PartNumber int       `json:"part_number"`
	PartCount  int       `json:"part_count"`
	Alphabet   string    `json:"alphabet"`
	MaxWordLen int       `json:"max_word_len"`
	TargetHash [16]byte  `json:"target_hash"`
	Algorithm  Algorithm `json:"algorithm"`
}

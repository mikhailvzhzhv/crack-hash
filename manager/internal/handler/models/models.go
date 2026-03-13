package models

import "github.com/google/uuid"

type Status string

const (
	StatusReady      = Status("READY")
	StatusInProgress = Status("IN_PROGRESS")
	StatusError      = Status("ERROR")
	StatusCancelled  = Status("CANCELLED")
)

type CrackHashRequest struct {
	Hash      string `json:"hash"`
	MaxLength int    `json:"maxLength"`
	Algorithm string `json:"algorithm"`
	Alphabet  string `json:"alphabet"`
}

type CrackHashResponse struct {
	RequestId             uuid.UUID `json:"requestId"`
	EstimatedCombinations int       `json:"estimated_combinations"`
}

type ResultResponse struct {
	Status Status   `json:"status"`
	Data   []string `json:"data"`
	Error  string   `json:"error,omitempty"`
}

type HashRequest struct {
	Word string `form:"word" binding:"required"`
}

type HashResponse struct {
	Hash string `json:"hash"`
}

type Statistic struct {
	TotalTasks       int `json:"total_tasks"`
	ActiveTasks      int `json:"active_tasks"`
	CompletedTasks   int `json:"completed_tasks"`
	AvgExecutionTime int `json:"avg_execution_time"`
	TasksParts       int `json:"-"`
}

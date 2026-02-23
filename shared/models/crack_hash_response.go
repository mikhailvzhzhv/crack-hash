package models

type Status string

const (
	StatusReady      = Status("READY")
	StatusInProgress = Status("IN_PROGRESS")
	StatusError      = Status("ERROR")
)

type CrackHashResponse struct {
	Status Status   `json:"status"`
	Data   []string `json:"data"`
}

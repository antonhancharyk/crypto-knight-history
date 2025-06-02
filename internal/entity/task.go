package entity

import "time"

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID      string               `json:"id"`
	Params  GetKlinesQueryParams `json:"params"`
	Status  TaskStatus           `json:"status"`
	StartAt time.Time            `json:"start_at,omitempty"`
	EndAt   time.Time            `json:"end_at,omitempty"`
	Result  []History            `json:"result,omitempty"`
	Error   string               `json:"error,omitempty"`
}

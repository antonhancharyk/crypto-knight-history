package entity

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID     string               `json:"id"`
	Params GetKlinesQueryParams `json:"params"`
	Status TaskStatus           `json:"status"`
	Result []History            `json:"result,omitempty"`
	Error  string               `json:"error,omitempty"`
}

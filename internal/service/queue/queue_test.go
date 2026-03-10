package queue

import (
	"context"
	"sync"
	"testing"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
)

type mockProcessor struct {
	mu       sync.Mutex
	called   bool
	lastParams entity.GetKlinesQueryParams
	result   []entity.History
	err     error
}

func (m *mockProcessor) ProcessHistory(ctx context.Context, params entity.GetKlinesQueryParams) ([]entity.History, error) {
	m.mu.Lock()
	m.called = true
	m.lastParams = params
	res, err := m.result, m.err
	m.mu.Unlock()
	return res, err
}

func TestCreateTask_GetTask(t *testing.T) {
	mock := &mockProcessor{result: []entity.History{{Symbol: "BTCUSDT"}}}
	q := New(mock)

	params := entity.GetKlinesQueryParams{From: "2024-01-01", To: "2024-01-02", Interval: "1h"}
	task := q.CreateTask(params)

	if task.ID == "" {
		t.Error("task ID is empty")
	}
	if task.Params.From != params.From || task.Params.Interval != params.Interval {
		t.Errorf("task params = %+v", task.Params)
	}
	if task.Status != entity.StatusPending {
		t.Errorf("status = %v", task.Status)
	}

	got, ok := q.GetTask(task.ID)
	if !ok {
		t.Fatal("GetTask: not found")
	}
	if got.ID != task.ID {
		t.Errorf("GetTask id = %q", got.ID)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	q := New(&mockProcessor{})
	_, ok := q.GetTask("nonexistent")
	if ok {
		t.Error("GetTask expected false for missing id")
	}
}

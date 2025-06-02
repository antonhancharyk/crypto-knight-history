package queue

import (
	"context"
	"sync"
	"time"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/service/kline"
	"github.com/google/uuid"
)

type HandlerFunc func(task *entity.Task) error

type TaskQueue struct {
	tasks    sync.Map
	jobs     chan string
	klineSvc *kline.Kline
}

func New(klineSvc *kline.Kline) *TaskQueue {
	q := &TaskQueue{
		jobs:     make(chan string, 100),
		klineSvc: klineSvc,
	}

	go q.worker()

	return q
}

func (q *TaskQueue) CreateTask(params entity.GetKlinesQueryParams) *entity.Task {
	id := uuid.New().String()
	task := &entity.Task{
		ID:     id,
		Params: params,
		Status: entity.StatusPending,
	}
	q.tasks.Store(id, task)
	q.jobs <- id

	return task
}

func (q *TaskQueue) GetTask(id string) (*entity.Task, bool) {
	val, ok := q.tasks.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*entity.Task), true
}

func (q *TaskQueue) worker() {
	for id := range q.jobs {
		val, ok := q.tasks.Load(id)
		if !ok {
			continue
		}
		task := val.(*entity.Task)
		task.Status = entity.StatusRunning
		task.StartAt = time.Now().UTC()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		result, err := q.klineSvc.ProcessHistory(ctx, task.Params)
		cancel()

		task.EndAt = time.Now().UTC()

		if err != nil {
			task.Status = entity.StatusFailed
			task.Error = err.Error()
		} else {
			task.Status = entity.StatusCompleted
			task.Result = result
		}

		go func() {
			time.Sleep(1 * time.Hour)
			q.tasks.Delete(task.ID)
		}()
	}
}

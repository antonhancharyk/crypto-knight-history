package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	apperr "github.com/antonhancharyk/crypto-knight-history/internal/errors"
)

type fakeBackend struct {
	klines    []entity.Kline
	klinesErr error
	tasks     map[string]*entity.Task
	nextTask  *entity.Task
}

func (f *fakeBackend) GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error) {
	if f.klinesErr != nil {
		return nil, f.klinesErr
	}
	return f.klines, nil
}

func (f *fakeBackend) CreateTask(params entity.GetKlinesQueryParams) *entity.Task {
	t := f.nextTask
	if t == nil {
		t = &entity.Task{ID: "test-id", Params: params, Status: entity.StatusPending}
	}
	if f.tasks != nil {
		f.tasks[t.ID] = t
	}
	return t
}

func (f *fakeBackend) GetTask(id string) (*entity.Task, bool) {
	if f.tasks != nil {
		t, ok := f.tasks[id]
		return t, ok
	}
	if f.nextTask != nil && f.nextTask.ID == id {
		return f.nextTask, true
	}
	return nil, false
}

func TestHandleHealth(t *testing.T) {
	srv := New(&fakeBackend{}, struct{ Port string }{Port: "0"})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	srv.handleHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if body := rec.Body.String(); body != "ok" {
		t.Errorf("body = %q, want \"ok\"", body)
	}
}

func TestHandleCreateTask(t *testing.T) {
	backend := &fakeBackend{
		tasks:    make(map[string]*entity.Task),
		nextTask: &entity.Task{ID: "created-123", Params: entity.GetKlinesQueryParams{Interval: "1h"}, Status: entity.StatusPending},
	}
	srv := New(backend, struct{ Port string }{Port: "0"})
	req := httptest.NewRequest(http.MethodPost, "/task?from=2024-01-01&to=2024-01-02&interval=1h", nil)
	rec := httptest.NewRecorder()

	srv.handleCreateTask(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", rec.Code)
	}
	var out map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out["task_id"] != "created-123" {
		t.Errorf("task_id = %q, want \"created-123\"", out["task_id"])
	}
}

func TestHandleTaskStatus_MissingID(t *testing.T) {
	srv := New(&fakeBackend{}, struct{ Port string }{Port: "0"})
	req := httptest.NewRequest(http.MethodGet, "/task/status", nil)
	rec := httptest.NewRecorder()

	srv.handleTaskStatus(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestHandleTaskStatus_NotFound(t *testing.T) {
	backend := &fakeBackend{tasks: make(map[string]*entity.Task)}
	srv := New(backend, struct{ Port string }{Port: "0"})
	req := httptest.NewRequest(http.MethodGet, "/task/status?id=missing", nil)
	rec := httptest.NewRecorder()

	srv.handleTaskStatus(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestHandleTaskStatus_Ok(t *testing.T) {
	task := &entity.Task{ID: "t1", Params: entity.GetKlinesQueryParams{}, Status: entity.StatusCompleted}
	backend := &fakeBackend{tasks: map[string]*entity.Task{"t1": task}}
	srv := New(backend, struct{ Port string }{Port: "0"})
	req := httptest.NewRequest(http.MethodGet, "/task/status?id=t1", nil)
	rec := httptest.NewRecorder()

	srv.handleTaskStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	var out entity.Task
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.ID != "t1" {
		t.Errorf("task id = %q", out.ID)
	}
}

func TestHandleGetKlines_Ok(t *testing.T) {
	klines := []entity.Kline{{Symbol: "BTCUSDT", Interval: "1h"}}
	backend := &fakeBackend{klines: klines}
	srv := New(backend, struct{ Port string }{Port: "0"})
	req := httptest.NewRequest(http.MethodGet, "/klines?"+url.Values{"from": {"2024-01-01"}, "to": {"2024-01-02"}, "interval": {"1h"}}.Encode(), nil)
	rec := httptest.NewRecorder()

	srv.handleGetKlines(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	var out []entity.Kline
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out) != 1 || out[0].Symbol != "BTCUSDT" {
		t.Errorf("klines = %v", out)
	}
}

func TestHandleGetKlines_BadRequest(t *testing.T) {
	backend := &fakeBackend{klinesErr: apperr.ErrBadRequest}
	srv := New(backend, struct{ Port string }{Port: "0"})
	req := httptest.NewRequest(http.MethodGet, "/klines", nil)
	rec := httptest.NewRecorder()

	srv.handleGetKlines(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestHandleGetKlines_InternalError(t *testing.T) {
	backend := &fakeBackend{klinesErr: errors.New("db error")}
	srv := New(backend, struct{ Port string }{Port: "0"})
	req := httptest.NewRequest(http.MethodGet, "/klines", nil)
	rec := httptest.NewRecorder()

	srv.handleGetKlines(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

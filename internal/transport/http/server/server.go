package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/antonhancharyk/crypto-knight-history/internal/config"
	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	apperr "github.com/antonhancharyk/crypto-knight-history/internal/errors"
)

// Backend is the interface required by HTTP handlers (allows testing with a fake).
type Backend interface {
	GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error)
	CreateTask(params entity.GetKlinesQueryParams) *entity.Task
	GetTask(id string) (*entity.Task, bool)
}

type HTTPServer struct {
	server *http.Server
	svc    Backend
}

func New(svc Backend, cfg config.Server) *HTTPServer {
	s := &HTTPServer{svc: svc}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/task", s.handleCreateTask)
	mux.HandleFunc("/task/status", s.handleTaskStatus)
	mux.HandleFunc("/klines", s.handleGetKlines)

	s.server = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	return s
}

func (s *HTTPServer) Start() error {
	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	return s.server.Serve(ln)
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *HTTPServer) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	params := entity.GetKlinesQueryParams{
		From:     r.URL.Query().Get("from"),
		To:       r.URL.Query().Get("to"),
		Symbol:   r.URL.Query().Get("symbol"),
		Interval: r.URL.Query().Get("interval"),
	}

	task := s.svc.CreateTask(params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"task_id": task.ID,
	})
}

func (s *HTTPServer) handleTaskStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, apperr.ErrBadRequest)
		return
	}

	task, ok := s.svc.GetTask(id)
	if !ok {
		writeError(w, apperr.ErrNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (s *HTTPServer) handleGetKlines(w http.ResponseWriter, r *http.Request) {
	params := entity.GetKlinesQueryParams{
		From:     r.URL.Query().Get("from"),
		To:       r.URL.Query().Get("to"),
		Symbol:   r.URL.Query().Get("symbol"),
		Interval: r.URL.Query().Get("interval"),
	}

	klines, err := s.svc.GetKlines(params)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(klines)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	msg := err.Error()

	switch {
	case errors.Is(err, apperr.ErrNotFound):
		status = http.StatusNotFound
		msg = "not found"
	case errors.Is(err, apperr.ErrBadRequest):
		status = http.StatusBadRequest
		msg = "bad request"
	}

	if status >= 500 {
		log.Printf("server error: %v", err)
	}
	http.Error(w, msg, status)
}

package server

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/service"
)

type HTTPServer struct {
	server *http.Server
	svc    *service.Service
	// queue  *queue.TaskQueue
}

func New(svc *service.Service) *HTTPServer {
	s := &HTTPServer{svc: svc}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/task", s.handleCreateTask)
	mux.HandleFunc("/task/status", s.handleTaskStatus)

	s.server = &http.Server{
		Addr:    ":" + os.Getenv("APP_SERVER_PORT"),
		Handler: mux,
	}

	return s
}

func (s *HTTPServer) Start() error {
	log.Printf("HTTP server is starting on %s...", s.server.Addr)

	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	return s.server.Serve(ln)
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	log.Println("shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *HTTPServer) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	params := entity.GetKlinesQueryParams{
		From:   r.URL.Query().Get("from"),
		To:     r.URL.Query().Get("to"),
		Symbol: r.URL.Query().Get("symbol"),
	}

	task := s.svc.Queue.CreateTask(params)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"task_id": task.ID,
	})
}

func (s *HTTPServer) handleTaskStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	task, ok := s.svc.Queue.GetTask(id)
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

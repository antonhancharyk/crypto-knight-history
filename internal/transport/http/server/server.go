package server

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/service"
)

type HTTPServer struct {
	server *http.Server
	svc    *service.Service
}

func New(addr string, svc *service.Service) *HTTPServer {
	s := &HTTPServer{svc: svc}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/history", s.handleHistory)

	s.server = &http.Server{
		Addr:    addr,
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

func (s *HTTPServer) handleHistory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	res, err := s.svc.Kline.ProcessHistory(ctx, entity.GetKlinesQueryParams{From: r.URL.Query().Get("from"),
		To:     r.URL.Query().Get("to"),
		Symbol: r.URL.Query().Get("symbol")})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

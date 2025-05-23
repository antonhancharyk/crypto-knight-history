package server

import (
	"context"
	"log"
	"net"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

func New(addr string) *HTTPServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/history", handleHistory)

	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &HTTPServer{server: s}
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

func handleHistory(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	symbol := r.URL.Query().Get("symbol")

	if from == "" || to == "" || symbol == "" {
		http.Error(w, "Missing query parameters: from, to, symbol", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","from":"` + from + `","to":"` + to + `","symbol":"` + symbol + `"}`))
}

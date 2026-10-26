package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Iposhka54/pulse/internal/collector"
)

type HttpServer struct {
	cfg       Config
	collector *collector.Collector
}

func NewHttpServer(cfg Config, col *collector.Collector) *HttpServer {
	return &HttpServer{cfg: cfg, collector: col}
}

func (s *HttpServer) Run(ctx context.Context) error {
	if s.collector == nil {
		return errors.New("http server: collector is nil")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /livez", s.handleLivez)
	mux.HandleFunc("GET /readyz", s.handleReadyz)
	mux.HandleFunc("GET /api/v1/status", s.handleStatus)

	srv := &http.Server{
		Addr:              s.cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", s.cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}

	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	log.Println("server stopped")
	return nil
}

func (s *HttpServer) handleLivez(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *HttpServer) handleReadyz(w http.ResponseWriter, _ *http.Request) {
	last, ok := s.collector.LastCollected()
	if !ok {
		http.Error(w, "no metrics collected yet", http.StatusServiceUnavailable)
		return
	}

	threshold := 3 * s.collector.Interval()
	if threshold < 15*time.Second {
		threshold = 15 * time.Second
	}

	if time.Since(last) > threshold {
		http.Error(w, fmt.Sprintf("metrics are stale: last collected %s ago", time.Since(last).Truncate(time.Second)), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ready"))
}

func (s *HttpServer) handleStatus(w http.ResponseWriter, _ *http.Request) {
	snapshot := s.collector.Latest()
	if snapshot == nil {
		http.Error(w, "no metrics collected yet", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(snapshot); err != nil {
		log.Printf("encode status: %v", err)
	}
}

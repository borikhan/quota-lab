package main

import (
	"encoding/json"
	"fmt"
	"github.com/borikhan/quota-lab/quota"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

func main() {
	capacity := envInt("QUOTA_CAPACITY", 100)
	limiter := quota.New(capacity, float64(envInt("REFILL_PER_SECOND", 10)))
	var allowed, rejected atomic.Uint64
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { write(w, 200, map[string]string{"status": "ok"}) })
	mux.HandleFunc("POST /v1/check", func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Tenant-ID")
		if key == "" {
			write(w, 400, map[string]string{"error": "X-Tenant-ID required"})
			return
		}
		d := limiter.Allow(key, 1)
		if !d.Allowed {
			rejected.Add(1)
			w.Header().Set("Retry-After", fmt.Sprint(max(1, int(d.RetryAfter.Seconds()))))
			write(w, 429, d)
			return
		}
		allowed.Add(1)
		write(w, 200, d)
	})
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("content-type", "text/plain; version=0.0.4")
		fmt.Fprintf(w, "quota_requests_total{decision=\"allowed\"} %d\nquota_requests_total{decision=\"rejected\"} %d\n", allowed.Load(), rejected.Load())
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := http.Server{Addr: ":" + port, Handler: mux, ReadHeaderTimeout: 3 * time.Second}
	slog.Info("quota api listening", "port", port)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("stopped", "error", err)
	}
}
func envInt(k string, f int) int {
	if v, err := strconv.Atoi(os.Getenv(k)); err == nil {
		return v
	}
	return f
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func write(w http.ResponseWriter, s int, v any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(s)
	_ = json.NewEncoder(w).Encode(v)
}

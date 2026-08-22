package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type pinger interface {
	Ping(context.Context) error
}

type app struct {
	db pinger
}

func newApp(db pinger) *app {
	return &app{db: db}
}

func (a *app) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", a.handleHealth)
	return mux
}

func (a *app) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	database := "not_configured"
	status := "ok"
	code := http.StatusOK
	if a.db != nil {
		if err := a.db.Ping(ctx); err != nil {
			database = "error"
			status = "degraded"
			code = http.StatusServiceUnavailable
		} else {
			database = "ok"
		}
	}

	writeJSON(w, code, map[string]string{
		"service":  "upskillrh-backend",
		"status":   status,
		"database": database,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

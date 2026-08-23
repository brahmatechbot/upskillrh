package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	platformauth "github.com/brahmatechbot/upskillrh/backend/internal/platform/auth"
)

type pinger interface {
	Ping(context.Context) error
}

type app struct {
	db          pinger
	authHandler *platformauth.Handler
}

func newApp(db pinger, authHandler ...*platformauth.Handler) *app {
	a := &app{db: db}
	if len(authHandler) > 0 {
		a.authHandler = authHandler[0]
	}
	return a
}

func (a *app) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", a.handleHealth)
	if a.authHandler != nil {
		mux.HandleFunc("/login", a.authHandler.LoginPage)
		mux.HandleFunc("/cadastro", a.authHandler.RegisterPage)
		mux.HandleFunc("/candidate", a.authHandler.CandidatePage)
		mux.HandleFunc("/app", a.authHandler.CompanyAppPage)
		mux.HandleFunc("/api/v1/auth/login", a.authHandler.LoginAPI)
		mux.HandleFunc("/api/v1/auth/register", a.authHandler.RegisterAPI)
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	}
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

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	platformauth "github.com/brahmatechbot/upskillrh/backend/internal/platform/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := configFromEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("create database pool: %v", err)
	}
	defer pool.Close()

	authRepo := platformauth.NewPostgresRepository(pool)
	authService := platformauth.NewLoginService(authRepo, platformauth.NewMemoryRateLimiter(5, 5*time.Minute))
	authHandler, err := platformauth.NewHandler(authService, "web/templates/auth/login.html")
	if err != nil {
		log.Fatalf("load login template: %v", err)
	}

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           newApp(pool, authHandler).routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("upskillrh backend listening on %s", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}

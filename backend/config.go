package main

import "os"

type config struct {
	HTTPAddr    string
	DatabaseURL string
}

func configFromEnv() config {
	return config{
		HTTPAddr:    envOrDefault("UPSKILLRH_HTTP_ADDR", "127.0.0.1:8092"),
		DatabaseURL: envOrDefault("UPSKILLRH_DATABASE_URL", "postgres://upskillrh:upskillrh@localhost:5432/upskillrh_dev?sslmode=disable"),
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

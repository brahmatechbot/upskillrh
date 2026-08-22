package main

import "testing"

func TestConfigFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("UPSKILLRH_HTTP_ADDR", "")
	t.Setenv("UPSKILLRH_DATABASE_URL", "")

	cfg := configFromEnv()

	if cfg.HTTPAddr != "127.0.0.1:8092" {
		t.Fatalf("expected default HTTP addr 127.0.0.1:8092, got %q", cfg.HTTPAddr)
	}
	if cfg.DatabaseURL == "" {
		t.Fatal("expected default database URL")
	}
}

func TestConfigFromEnvReadsOverrides(t *testing.T) {
	t.Setenv("UPSKILLRH_HTTP_ADDR", ":9090")
	t.Setenv("UPSKILLRH_DATABASE_URL", "postgres://custom")

	cfg := configFromEnv()

	if cfg.HTTPAddr != ":9090" {
		t.Fatalf("expected override HTTP addr, got %q", cfg.HTTPAddr)
	}
	if cfg.DatabaseURL != "postgres://custom" {
		t.Fatalf("expected override database URL, got %q", cfg.DatabaseURL)
	}
}

package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakePinger struct {
	err error
}

func (p fakePinger) Ping(_ context.Context) error {
	return p.err
}

func TestHealthEndpointReportsOKWhenDatabasePings(t *testing.T) {
	app := newApp(fakePinger{})
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	app.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"status":"ok"`) {
		t.Fatalf("expected status ok in body, got %s", body)
	}
	if !strings.Contains(body, `"database":"ok"`) {
		t.Fatalf("expected database ok in body, got %s", body)
	}
}

func TestHealthEndpointRejectsUnsupportedMethod(t *testing.T) {
	app := newApp(fakePinger{})
	req := httptest.NewRequest(http.MethodPost, "/api/health", nil)
	rec := httptest.NewRecorder()

	app.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

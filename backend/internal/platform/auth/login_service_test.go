package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeStore struct {
	user          User
	hash          string
	createSession bool
}

func (s *fakeStore) FindUserCredentials(_ context.Context, _ string) (User, string, error) {
	return s.user, s.hash, nil
}
func (s *fakeStore) CreateSession(_ context.Context, _ Session) error {
	s.createSession = true
	return nil
}
func (s *fakeStore) UpdateLastLogin(_ context.Context, _ string) error { return nil }

func TestPasswordHashAndVerify(t *testing.T) {
	hash, err := HashPassword("senha secreta")
	if err != nil {
		t.Fatal(err)
	}
	ok, err := VerifyPassword("senha secreta", hash)
	if err != nil || !ok {
		t.Fatalf("expected valid password, ok=%v err=%v", ok, err)
	}
	ok, err = VerifyPassword("outra", hash)
	if err != nil || ok {
		t.Fatalf("expected invalid password, ok=%v err=%v", ok, err)
	}
}

func TestLoginServiceCreatesOpaqueSession(t *testing.T) {
	hash, err := HashPassword("senha")
	if err != nil {
		t.Fatal(err)
	}
	store := &fakeStore{user: User{ID: "00000000-0000-0000-0000-000000000001", Email: "pessoa@empresa.com.br", DisplayName: "Pessoa", Status: "active", SessionVersion: 1}, hash: hash}
	service := NewLoginService(store, NewMemoryRateLimiter(5, 300000000000))
	result, err := service.Login(context.Background(), " pessoa@empresa.com.br ", "senha", false, "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if !store.createSession || result.SessionRaw == "" || result.NextURL != "/app" {
		t.Fatalf("expected session and next_url, result=%+v", result)
	}
}

func TestLoginHandlerRejectsInvalidCredentialsWithoutEnumeration(t *testing.T) {
	hash, err := HashPassword("senha-correta")
	if err != nil {
		t.Fatal(err)
	}
	store := &fakeStore{user: User{ID: "00000000-0000-0000-0000-000000000001", Email: "pessoa@empresa.com.br", DisplayName: "Pessoa", Status: "active", SessionVersion: 1}, hash: hash}
	handler := &Handler{service: NewLoginService(store, NewMemoryRateLimiter(5, 300000000000))}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"email":"pessoa@empresa.com.br","password":"errada","remember_me":false}`))
	req.Header.Set("X-CSRF-Token", "csrf")
	req.AddCookie(&http.Cookie{Name: "upskillrh_login_csrf", Value: "csrf"})
	rec := httptest.NewRecorder()
	handler.LoginAPI(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "E-mail ou senha inválidos") {
		t.Fatalf("expected generic message, got %s", rec.Body.String())
	}
}

func TestRateLimit(t *testing.T) {
	limiter := NewMemoryRateLimiter(1, 300000000000)
	if ok, _ := limiter.Allow("k"); !ok {
		t.Fatal("first attempt should pass")
	}
	if ok, _ := limiter.Allow("k"); ok {
		t.Fatal("second attempt should be blocked")
	}
}

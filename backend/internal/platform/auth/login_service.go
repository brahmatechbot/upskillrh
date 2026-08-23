package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccessUnavailable  = errors.New("access unavailable")
	ErrRateLimited        = errors.New("rate limited")
)

type User struct {
	ID             string
	Email          string
	DisplayName    string
	Status         string
	SessionVersion int
}

type Store interface {
	FindUserCredentials(ctx context.Context, email string) (User, string, error)
	CreateSession(ctx context.Context, session Session) error
	UpdateLastLogin(ctx context.Context, userID string) error
}

type LoginDestinationResolver interface {
	LoginDestination(ctx context.Context, userID string) (string, error)
}

type RateLimiter interface {
	Allow(key string) (bool, time.Duration)
}

type LoginService struct {
	store       Store
	rateLimiter RateLimiter
	now         func() time.Time
}

type Session struct {
	ID                string
	UserID            string
	TokenHash         []byte
	CSRFHash          []byte
	SessionVersion    int
	RememberMe        bool
	CreatedAt         time.Time
	LastSeenAt        time.Time
	IdleExpiresAt     time.Time
	AbsoluteExpiresAt time.Time
}

type LoginResult struct {
	User        User
	NextURL     string
	SessionID   string
	SessionRaw  string
	CSRFToken   string
	IdleExpiry  time.Time
	AbsoluteExp time.Time
}

func NewLoginService(store Store, limiter RateLimiter) *LoginService {
	return &LoginService{store: store, rateLimiter: limiter, now: time.Now}
}

func (s *LoginService) Login(ctx context.Context, email, password string, rememberMe bool, remoteAddr string) (LoginResult, error) {
	if s.store == nil {
		return LoginResult{}, ErrAccessUnavailable
	}
	key := strings.ToLower(strings.TrimSpace(email)) + "|" + remoteAddr
	if s.rateLimiter != nil {
		if ok, _ := s.rateLimiter.Allow(key); !ok {
			return LoginResult{}, ErrRateLimited
		}
	}
	user, hash, err := s.store.FindUserCredentials(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}
	if user.Status != "active" {
		return LoginResult{}, ErrAccessUnavailable
	}
	ok, err := VerifyPassword(password, hash)
	if err != nil || !ok {
		return LoginResult{}, ErrInvalidCredentials
	}
	rawToken, err := randomToken(32)
	if err != nil {
		return LoginResult{}, err
	}
	csrfToken, err := randomToken(32)
	if err != nil {
		return LoginResult{}, err
	}
	now := s.now().UTC()
	idle := now.Add(12 * time.Hour)
	absolute := now.Add(24 * time.Hour)
	if rememberMe {
		idle = now.Add(30 * 24 * time.Hour)
		absolute = now.Add(30 * 24 * time.Hour)
	}
	session := Session{
		ID:                "",
		UserID:            user.ID,
		TokenHash:         hashToken(rawToken),
		CSRFHash:          hashToken(csrfToken),
		SessionVersion:    user.SessionVersion,
		RememberMe:        rememberMe,
		CreatedAt:         now,
		LastSeenAt:        now,
		IdleExpiresAt:     idle,
		AbsoluteExpiresAt: absolute,
	}
	if err := s.store.CreateSession(ctx, session); err != nil {
		return LoginResult{}, err
	}
	_ = s.store.UpdateLastLogin(ctx, user.ID)
	nextURL := "/app"
	if resolver, ok := s.store.(LoginDestinationResolver); ok {
		if resolved, err := resolver.LoginDestination(ctx, user.ID); err == nil && resolved != "" {
			nextURL = resolved
		}
	}
	return LoginResult{User: user, NextURL: nextURL, SessionRaw: rawToken, CSRFToken: csrfToken, IdleExpiry: idle, AbsoluteExp: absolute}, nil
}

func randomToken(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}

func secureCookie(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}

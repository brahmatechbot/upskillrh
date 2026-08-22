package auth

import (
	"encoding/json"
	"errors"
	"html/template"
	"net"
	"net/http"
	"strings"
	"unicode/utf8"
)

type Handler struct {
	service *LoginService
	tmpl    *template.Template
}

func NewHandler(service *LoginService, templatePath string) (*Handler, error) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}
	return &Handler{service: service, tmpl: t}, nil
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	csrf, err := randomToken(32)
	if err != nil {
		http.Error(w, "Não foi possível entrar agora. Tente novamente.", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "upskillrh_login_csrf", Value: csrf, Path: "/login", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: secureCookie(r), MaxAge: 1800})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.tmpl.ExecuteTemplate(w, "login.html", map[string]string{"CSRFToken": csrf})
}

type loginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

func (h *Handler) LoginAPI(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID, _ = randomToken(12)
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed", requestID, nil)
		return
	}
	if !h.validCSRF(r) {
		writeAPIError(w, http.StatusForbidden, "CSRF_INVALID", "Atualize a página e tente novamente.", requestID, nil)
		return
	}
	defer r.Body.Close()
	var in loginRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 2048)).Decode(&in); err != nil {
		writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "JSON inválido.", requestID, nil)
		return
	}
	fields := validateLogin(in)
	if len(fields) > 0 {
		writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Revise os campos informados.", requestID, fields)
		return
	}
	result, err := h.service.Login(r.Context(), strings.TrimSpace(in.Email), in.Password, in.RememberMe, clientIP(r))
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			writeAPIError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "E-mail ou senha inválidos.", requestID, nil)
		case errors.Is(err, ErrAccessUnavailable):
			writeAPIError(w, http.StatusForbidden, "ACCESS_UNAVAILABLE", "Acesso indisponível. Fale com o suporte.", requestID, nil)
		case errors.Is(err, ErrRateLimited):
			w.Header().Set("Retry-After", "300")
			writeAPIError(w, http.StatusTooManyRequests, "RATE_LIMITED", "Muitas tentativas. Aguarde alguns minutos antes de tentar novamente.", requestID, nil)
		default:
			writeAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Não foi possível entrar agora. Tente novamente.", requestID, nil)
		}
		return
	}
	exp := result.AbsoluteExp
	http.SetCookie(w, &http.Cookie{Name: "upskillrh_session", Value: result.SessionRaw, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: secureCookie(r), Expires: exp})
	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]any{"user": map[string]string{"id": result.User.ID, "display_name": result.User.DisplayName}, "next_url": result.NextURL}, "meta": map[string]string{"request_id": requestID}})
}

func (h *Handler) validCSRF(r *http.Request) bool {
	cookie, err := r.Cookie("upskillrh_login_csrf")
	if err != nil || cookie.Value == "" {
		return false
	}
	return r.Header.Get("X-CSRF-Token") == cookie.Value
}

func validateLogin(in loginRequest) map[string]string {
	fields := map[string]string{}
	email := strings.TrimSpace(in.Email)
	if email == "" {
		fields["email"] = "Informe seu e-mail."
	} else if len(email) > 254 || !strings.Contains(email, "@") || strings.ContainsAny(email, " \t\n\r") {
		fields["email"] = "Informe um e-mail válido."
	}
	if in.Password == "" || len(in.Password) > 1024 || !utf8.ValidString(in.Password) {
		fields["password"] = "Informe sua senha."
	}
	return fields
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]); forwarded != "" {
		return forwarded
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeAPIError(w http.ResponseWriter, status int, code, message, requestID string, fields map[string]string) {
	body := map[string]any{"error": map[string]any{"code": code, "message": message}, "meta": map[string]string{"request_id": requestID}}
	if len(fields) > 0 {
		body["error"].(map[string]any)["fields"] = fields
	}
	writeJSON(w, status, body)
}

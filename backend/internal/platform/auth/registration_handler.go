package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	csrf, err := randomToken(32)
	if err != nil {
		http.Error(w, "Não foi possível abrir o cadastro agora.", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "upskillrh_register_csrf", Value: csrf, Path: "/cadastro", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: secureCookie(r), MaxAge: 1800})
	initialType := r.URL.Query().Get("tipo")
	if initialType != "empresa" && initialType != "candidato" {
		initialType = ""
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.tmpl.ExecuteTemplate(w, "register.html", map[string]string{"CSRFToken": csrf, "InitialType": initialType})
}

func (h *Handler) RegisterAPI(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID, _ = randomToken(12)
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed", requestID, nil)
		return
	}
	if h.registrationService == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "REGISTRATION_UNAVAILABLE", "Cadastro indisponível no momento.", requestID, nil)
		return
	}
	if !h.validRegisterCSRF(r) {
		writeAPIError(w, http.StatusForbidden, "CSRF_INVALID", "Atualize a página e tente novamente.", requestID, nil)
		return
	}
	defer r.Body.Close()
	var in RegisterInput
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192)).Decode(&in); err != nil {
		writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "JSON inválido.", requestID, nil)
		return
	}
	result, fields, err := h.registrationService.Register(r.Context(), in, requestID)
	if len(fields) > 0 {
		writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Revise os campos informados.", requestID, fields)
		return
	}
	if err != nil {
		if errors.Is(err, ErrAccessUnavailable) {
			writeAPIError(w, http.StatusServiceUnavailable, "ACCESS_UNAVAILABLE", "Cadastro indisponível no momento.", requestID, nil)
			return
		}
		writeAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Não foi possível criar a conta agora.", requestID, nil)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "upskillrh_session", Value: result.SessionRaw, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: secureCookie(r), Expires: result.SessionExpiresAt})
	writeJSON(w, http.StatusCreated, map[string]any{"data": result, "meta": map[string]string{"request_id": requestID}})
}

func (h *Handler) CandidatePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx, ok := h.requireSession(w, r, "/candidate")
	if !ok {
		return
	}
	if ctx.CandidateProfileID == "" {
		http.Error(w, "Perfil de candidato não encontrado.", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.tmpl.ExecuteTemplate(w, "candidate.html", ctx)
}

func (h *Handler) CompanyAppPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx, ok := h.requireSession(w, r, "/app")
	if !ok {
		return
	}
	if ctx.OrganizationID == "" {
		http.Error(w, "Vínculo empresarial não encontrado.", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.tmpl.ExecuteTemplate(w, "app.html", ctx)
}

func (h *Handler) requireSession(w http.ResponseWriter, r *http.Request, back string) (SessionContext, bool) {
	if h.registrationService == nil || h.registrationService.store == nil {
		http.Error(w, "Sessão indisponível.", http.StatusServiceUnavailable)
		return SessionContext{}, false
	}
	raw, ok := sessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/login?next="+back, http.StatusFound)
		return SessionContext{}, false
	}
	ctx, err := h.registrationService.store.SessionContext(r.Context(), raw)
	if err != nil {
		http.Redirect(w, r, "/login?next="+back, http.StatusFound)
		return SessionContext{}, false
	}
	return ctx, true
}

func (h *Handler) validRegisterCSRF(r *http.Request) bool {
	cookie, err := r.Cookie("upskillrh_register_csrf")
	if err != nil || cookie.Value == "" {
		return false
	}
	return strings.TrimSpace(r.Header.Get("X-CSRF-Token")) == cookie.Value
}

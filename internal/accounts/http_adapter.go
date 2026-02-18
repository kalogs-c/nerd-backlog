package accounts

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/kalogs-c/nerd-backlog/pkg/auth"
	"github.com/kalogs-c/nerd-backlog/pkg/httpjson"
)

type HTTPAdapter struct {
	service domain.AccountService
	logger  *slog.Logger
}

func NewHTTPAdapter(s domain.AccountService, logger *slog.Logger) *HTTPAdapter {
	return &HTTPAdapter{s, logger}
}

func (h *HTTPAdapter) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := httpjson.DecodeValid[*LoginPayload](r)
	if err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusBadRequest, "invalid payload", err)
		return
	}

	account, session, err := h.service.Login(ctx, payload.Email, payload.Password)
	if err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusInternalServerError, "failed to login", err)
		return
	}

	setSessionCookie(w, r, session)

	response := MountAccountResponse(account)
	if err := httpjson.Encode(w, r, http.StatusCreated, response); err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusInternalServerError, "failed to encode account", err)
	}
}

func (h *HTTPAdapter) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := httpjson.DecodeValid[*RegisterPayload](r)
	if err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusBadRequest, "invalid payload", err)
		return
	}

	account, session, err := h.service.Register(ctx, payload.Nickname, payload.Email, payload.Password)
	if err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusInternalServerError, "failed to register", err)
		return
	}

	setSessionCookie(w, r, session)

	response := MountAccountResponse(account)
	if err := httpjson.Encode(w, r, http.StatusCreated, response); err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusInternalServerError, "failed to encode account", err)
	}
}

func (h *HTTPAdapter) Logout(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := auth.SessionTokenFromContext(r.Context())
	if !ok {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusUnauthorized, "missing session", auth.ErrInvalidToken)
		return
	}

	if err := h.service.LogoutSession(r.Context(), sessionToken); err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusInternalServerError, "failed to logout", err)
		return
	}

	clearSessionCookie(w, r)
	w.WriteHeader(http.StatusNoContent)
}

func setSessionCookie(w http.ResponseWriter, r *http.Request, session domain.Session) {
	age := int(time.Until(session.ExpiresAt).Seconds())
	maxAge := max(age, 0)

	cookie := &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

func clearSessionCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

func isSecureRequest(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}

	forwardedProto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto"))
	if forwardedProto != "" {
		return strings.EqualFold(forwardedProto, "https")
	}

	return forwardedHeaderIsHTTPS(r.Header.Get("Forwarded"))
}

func forwardedHeaderIsHTTPS(value string) bool {
	if value == "" {
		return false
	}

	for entry := range strings.SplitSeq(value, ",") {
		for part := range strings.SplitSeq(entry, ";") {
			part = strings.TrimSpace(part)
			if !strings.HasPrefix(strings.ToLower(part), "proto=") {
				continue
			}
			proto := strings.TrimSpace(strings.TrimPrefix(part, "proto="))
			proto = strings.Trim(proto, `"`)
			return strings.EqualFold(proto, "https")
		}
	}

	return false
}

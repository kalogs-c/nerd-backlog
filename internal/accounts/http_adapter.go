package accounts

import (
	"log/slog"
	"net/http"

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

	account, tokenPair, err := h.service.Login(ctx, payload.Email, payload.Password)
	if err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusInternalServerError, "failed to login", err)
		return
	}

	response := LoginResponse{
		Account:   MountAccountResponse(account),
		TokenPair: MountTokenPairResponse(tokenPair),
	}

	if err := httpjson.Encode(w, r, http.StatusCreated, response); err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusInternalServerError, "failed to encode account", err)
	}
}

func (h *HTTPAdapter) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := httpjson.DecodeValid[*SignupPayload](r)
	if err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusBadRequest, "invalid payload", err)
		return
	}

	account, tokenPair, err := h.service.Signup(ctx, payload.Nickname, payload.Email, payload.Password)
	if err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusInternalServerError, "failed to signup", err)
		return
	}

	response := LoginResponse{
		Account:   MountAccountResponse(account),
		TokenPair: MountTokenPairResponse(tokenPair),
	}

	if err := httpjson.Encode(w, r, http.StatusCreated, response); err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusInternalServerError, "failed to encode account", err)
	}
}

func (h *HTTPAdapter) Logout(w http.ResponseWriter, r *http.Request) {
	accountID, ok := auth.AccountIDFromContext(r.Context())
	if !ok {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusUnauthorized, "missing account", auth.ErrInvalidToken)
		return
	}

	if err := h.service.Logout(r.Context(), accountID); err != nil {
		httpjson.NotifyHTTPError(w, r, h.logger, http.StatusInternalServerError, "failed to logout", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

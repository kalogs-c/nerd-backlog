package games

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/kalogs-c/nerd-backlog/pkg/httpjson"
	"github.com/kalogs-c/nerd-backlog/pkg/validator"
)

type HTTPAdapter struct {
	service domain.GameService
	logger  *slog.Logger
}

func NewHTTPAdapter(s domain.GameService, logger *slog.Logger) *HTTPAdapter {
	return &HTTPAdapter{s, logger}
}

func (h *HTTPAdapter) error(w http.ResponseWriter, r *http.Request, code int, title string, err error) {
	httpjson.NotifyError(
		r.Context(),
		w,
		r,
		h.logger,
		code,
		title,
		err.Error(),
		err,
	)
}

func (h *HTTPAdapter) CreateGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := httpjson.DecodeValid[*CreateGamePayload](r)
	if err != nil {
		switch e := err.(type) {
		case validator.ValidationError:
			httpjson.EncodeValidationErrors(w, r, e.Problems)
		default:
			h.error(w, r, http.StatusBadRequest, "invalid payload", err)
		}
		return
	}

	room, err := h.service.CreateGame(ctx, payload.Title)
	if err != nil {
		h.error(w, r, http.StatusInternalServerError, "failed to create room", err)
		return
	}

	if err := httpjson.Encode(w, r, http.StatusCreated, room); err != nil {
		h.error(w, r, http.StatusInternalServerError, "failed to encode room", err)
	}
}

func (h *HTTPAdapter) GetGameByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idString := r.PathValue("id")

	id, err := uuid.Parse(idString)
	if err != nil {
		h.error(w, r, http.StatusUnprocessableEntity, "failed to parse id", err)
		return
	}

	room, err := h.service.GetGameByID(ctx, id)
	if err != nil {
		h.error(w, r, http.StatusNotFound, "failed to retrieve room", err)
		return
	}

	if err := httpjson.Encode(w, r, http.StatusOK, room); err != nil {
		h.error(w, r, http.StatusInternalServerError, "failed to encode room", err)
	}
}

func (h *HTTPAdapter) ListGames(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.service.ListGames(r.Context())
	if err != nil {
		h.error(w, r, http.StatusInternalServerError, "failed to list rooms", err)
		return
	}

	if err := httpjson.Encode(w, r, http.StatusOK, rooms); err != nil {
		h.error(w, r, http.StatusInternalServerError, "failed to encode rooms", err)
	}
}

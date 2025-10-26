package games

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHTTPAdapter_CreateGame(t *testing.T) {
	mockSvc := new(MockGameService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	want := domain.Game{ID: uuid.New(), Title: "Zelda"}
	mockSvc.On("CreateGame", mock.Anything, "Zelda").Return(want, nil)

	body := bytes.NewBufferString(`{"title":"Zelda"}`)
	req := httptest.NewRequest(http.MethodPost, "/games", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateGame(w, req)
	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusCreated, res.StatusCode)

	var got domain.Game
	require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.Title, got.Title)

	mockSvc.AssertExpectations(t)
}

func TestHTTPAdapter_CreateGame_BadPayload(t *testing.T) {
	mockSvc := new(MockGameService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	body := bytes.NewBufferString(`{"title":123}`)
	req := httptest.NewRequest(http.MethodPost, "/games", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateGame(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHTTPAdapter_GetGameByID(t *testing.T) {
	mockSvc := new(MockGameService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	id := uuid.New()
	want := domain.Game{ID: id, Title: "Celeste"}
	mockSvc.On("GetGameByID", mock.Anything, id).Return(want, nil)

	req := httptest.NewRequest(http.MethodGet, "/games/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	handler.GetGameByID(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var got domain.Game
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	require.Equal(t, want.Title, got.Title)
	mockSvc.AssertExpectations(t)
}

func TestHTTPAdapter_GetGameByID_InvalidUUID(t *testing.T) {
	mockSvc := new(MockGameService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	req := httptest.NewRequest(http.MethodGet, "/games/invalid-uuid", nil)
	req.SetPathValue("id", "invalid-uuid")
	w := httptest.NewRecorder()

	handler.GetGameByID(w, req)
	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestHTTPAdapter_ListGames(t *testing.T) {
	mockSvc := new(MockGameService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	want := []domain.Game{
		{ID: uuid.New(), Title: "Hades"},
		{ID: uuid.New(), Title: "Celeste"},
	}
	mockSvc.On("ListGames", mock.Anything).Return(want, nil)

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()

	handler.ListGames(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var got []domain.Game
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	require.Len(t, got, 2)
	require.Equal(t, "Hades", got[0].Title)
	mockSvc.AssertExpectations(t)
}

func TestHTTPAdapter_DeleteGameByID(t *testing.T) {
	mockSvc := new(MockGameService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	id := uuid.New()
	mockSvc.On("DeleteGameByID", mock.Anything, id).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/games/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	handler.DeleteGameByID(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestHTTPAdapter_DeleteGameByID_BadUUID(t *testing.T) {
	mockSvc := new(MockGameService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	req := httptest.NewRequest(http.MethodDelete, "/games/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.DeleteGameByID(w, req)
	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestHTTPAdapter_DeleteGameByID_ServiceError(t *testing.T) {
	mockSvc := new(MockGameService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	id := uuid.New()
	mockSvc.On("DeleteGameByID", mock.Anything, id).Return(errors.New("delete failed"))

	req := httptest.NewRequest(http.MethodDelete, "/games/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	handler.DeleteGameByID(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertExpectations(t)
}

package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/kalogs-c/nerd-backlog/pkg/auth"
)

func TestWithAuth_MissingToken(t *testing.T) {
	jwtManager := auth.NewJWTManager([]byte("secret"), time.Minute, time.Hour)
	middleware := WithAuth(jwtManager, nil)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWithAuth_InvalidToken(t *testing.T) {
	jwtManager := auth.NewJWTManager([]byte("secret"), time.Minute, time.Hour)
	middleware := WithAuth(jwtManager, nil)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	req.Header.Set("Authorization", "Bearer not-a-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWithAuth_ValidToken(t *testing.T) {
	jwtManager := auth.NewJWTManager([]byte("secret"), time.Minute, time.Hour)
	accountID := uuid.New()
	accessToken, err := jwtManager.GenerateAccessToken(accountID)
	require.NoError(t, err)

	middleware := WithAuth(jwtManager, nil)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID, ok := AccountIDFromContext(r.Context())
		if !ok || gotID != accountID {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/kalogs-c/nerd-backlog/pkg/auth"
)

type stubSessionStore struct {
	accountID uuid.UUID
	err       error
}

func (s stubSessionStore) GetSessionAccountID(ctx context.Context, token string) (uuid.UUID, error) {
	if s.err != nil {
		return uuid.Nil, s.err
	}
	return s.accountID, nil
}

func TestWithAuth_MissingToken(t *testing.T) {
	middleware := WithAuth(nil, nil)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWithAuth_MissingSessionCookie(t *testing.T) {
	sessionStore := stubSessionStore{accountID: uuid.New()}
	middleware := WithAuth(sessionStore, nil)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWithAuth_InvalidSession(t *testing.T) {
	sessionStore := stubSessionStore{err: auth.ErrInvalidToken}
	middleware := WithAuth(sessionStore, nil)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "session-token"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWithAuth_SessionCookie(t *testing.T) {
	accountID := uuid.New()
	sessionStore := stubSessionStore{accountID: accountID}

	middleware := WithAuth(sessionStore, nil)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID, ok := auth.AccountIDFromContext(r.Context())
		if !ok || gotID != accountID {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		token, ok := auth.SessionTokenFromContext(r.Context())
		if !ok || token != "session-token" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "session-token"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

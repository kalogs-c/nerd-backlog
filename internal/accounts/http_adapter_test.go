package accounts

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/kalogs-c/nerd-backlog/pkg/auth"
)

func TestHTTPAdapter_Login(t *testing.T) {
	mockSvc := new(MockAccountService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	account := domain.Account{
		ID:       uuid.New(),
		Nickname: "nerd",
		Email:    "nerd@example.com",
	}
	session := domain.Session{
		Token:     "session-token",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockSvc.On("Login", mock.Anything, account.Email, "password").Return(account, session, nil)

	body := bytes.NewBufferString(`{"email":"nerd@example.com","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var got AccountResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	require.Equal(t, account.ID, got.ID)
	require.Equal(t, account.Nickname, got.Nickname)
	require.Equal(t, account.Email, got.Email)

	res := w.Result()
	var cookie *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == auth.SessionCookieName {
			cookie = c
			break
		}
	}
	require.NotNil(t, cookie)
	require.Equal(t, session.Token, cookie.Value)
	require.NoError(t, res.Body.Close())

	mockSvc.AssertExpectations(t)
}

func TestHTTPAdapter_Login_BadPayload(t *testing.T) {
	mockSvc := new(MockAccountService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	body := bytes.NewBufferString(`{"email":123}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHTTPAdapter_Login_ServiceError(t *testing.T) {
	mockSvc := new(MockAccountService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	mockSvc.On("Login", mock.Anything, "nerd@example.com", "password").Return(domain.Account{}, domain.Session{}, errors.New("login failed"))

	body := bytes.NewBufferString(`{"email":"nerd@example.com","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	mockSvc.AssertExpectations(t)
}

func TestHTTPAdapter_Register(t *testing.T) {
	mockSvc := new(MockAccountService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	account := domain.Account{
		ID:       uuid.New(),
		Nickname: "nerd",
		Email:    "nerd@example.com",
	}
	session := domain.Session{
		Token:     "session-token",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockSvc.On("Register", mock.Anything, account.Nickname, account.Email, "password").Return(account, session, nil)

	body := bytes.NewBufferString(`{"nickname":"nerd","email":"nerd@example.com","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var got AccountResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	require.Equal(t, account.ID, got.ID)
	require.Equal(t, account.Nickname, got.Nickname)
	require.Equal(t, account.Email, got.Email)

	res := w.Result()
	var cookie *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == auth.SessionCookieName {
			cookie = c
			break
		}
	}
	require.NotNil(t, cookie)
	require.Equal(t, session.Token, cookie.Value)
	require.NoError(t, res.Body.Close())

	mockSvc.AssertExpectations(t)
}

func TestHTTPAdapter_Register_BadPayload(t *testing.T) {
	mockSvc := new(MockAccountService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	body := bytes.NewBufferString(`{"nickname":123}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHTTPAdapter_Register_ServiceError(t *testing.T) {
	mockSvc := new(MockAccountService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	mockSvc.On("Register", mock.Anything, "nerd", "nerd@example.com", "password").Return(domain.Account{}, domain.Session{}, errors.New("register failed"))

	body := bytes.NewBufferString(`{"nickname":"nerd","email":"nerd@example.com","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	mockSvc.AssertExpectations(t)
}

func TestHTTPAdapter_Logout(t *testing.T) {
	mockSvc := new(MockAccountService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	mockSvc.On("LogoutSession", mock.Anything, "session-token").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req = req.WithContext(auth.WithSessionToken(req.Context(), "session-token"))
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	res := w.Result()
	var cookie *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == auth.SessionCookieName {
			cookie = c
			break
		}
	}
	require.NotNil(t, cookie)
	require.Equal(t, -1, cookie.MaxAge)
	require.NoError(t, res.Body.Close())

	mockSvc.AssertExpectations(t)
}

func TestHTTPAdapter_Logout_MissingSession(t *testing.T) {
	mockSvc := new(MockAccountService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHTTPAdapter_Logout_ServiceError(t *testing.T) {
	mockSvc := new(MockAccountService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	mockSvc.On("LogoutSession", mock.Anything, "session-token").Return(errors.New("logout failed"))

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req = req.WithContext(auth.WithSessionToken(req.Context(), "session-token"))
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestIsSecureRequest(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*http.Request)
		expected bool
	}{
		{
			name: "tls",
			setup: func(r *http.Request) {
				r.TLS = &tls.ConnectionState{}
			},
			expected: true,
		},
		{
			name: "x-forwarded-proto https",
			setup: func(r *http.Request) {
				r.Header.Set("X-Forwarded-Proto", "https")
			},
			expected: true,
		},
		{
			name: "x-forwarded-proto http",
			setup: func(r *http.Request) {
				r.Header.Set("X-Forwarded-Proto", "http")
			},
			expected: false,
		},
		{
			name: "forwarded proto https",
			setup: func(r *http.Request) {
				r.Header.Set("Forwarded", "for=1.1.1.1;proto=https;host=example.com")
			},
			expected: true,
		},
		{
			name: "forwarded proto quoted https",
			setup: func(r *http.Request) {
				r.Header.Set("Forwarded", "proto=\"https\"")
			},
			expected: true,
		},
		{
			name: "forwarded proto http",
			setup: func(r *http.Request) {
				r.Header.Set("Forwarded", "for=1.1.1.1;proto=http")
			},
			expected: false,
		},
		{
			name:     "no tls or forwarded",
			setup:    func(r *http.Request) {},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			test.setup(req)
			require.Equal(t, test.expected, isSecureRequest(req))
		})
	}
}

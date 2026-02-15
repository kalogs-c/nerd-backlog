package accounts

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

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
	tokenPair := domain.TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	mockSvc.On("Login", mock.Anything, account.Email, "password").Return(account, tokenPair, nil)

	body := bytes.NewBufferString(`{"email":"nerd@example.com","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var got LoginResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	require.Equal(t, account.ID, got.Account.ID)
	require.Equal(t, account.Nickname, got.Account.Nickname)
	require.Equal(t, account.Email, got.Account.Email)
	require.Equal(t, tokenPair.AccessToken, got.TokenPair.AccessToken)
	require.Equal(t, tokenPair.RefreshToken, got.TokenPair.RefreshToken)

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

	mockSvc.On("Login", mock.Anything, "nerd@example.com", "password").Return(domain.Account{}, domain.TokenPair{}, errors.New("login failed"))

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
	tokenPair := domain.TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	mockSvc.On("Register", mock.Anything, account.Nickname, account.Email, "password").Return(account, tokenPair, nil)

	body := bytes.NewBufferString(`{"nickname":"nerd","email":"nerd@example.com","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var got LoginResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	require.Equal(t, account.ID, got.Account.ID)
	require.Equal(t, account.Nickname, got.Account.Nickname)
	require.Equal(t, account.Email, got.Account.Email)
	require.Equal(t, tokenPair.AccessToken, got.TokenPair.AccessToken)
	require.Equal(t, tokenPair.RefreshToken, got.TokenPair.RefreshToken)

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

	mockSvc.On("Register", mock.Anything, "nerd", "nerd@example.com", "password").Return(domain.Account{}, domain.TokenPair{}, errors.New("register failed"))

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

	accountID := uuid.New()
	mockSvc.On("Logout", mock.Anything, accountID).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req = req.WithContext(auth.WithAccountID(req.Context(), accountID))
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestHTTPAdapter_Logout_MissingAccount(t *testing.T) {
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

	accountID := uuid.New()
	mockSvc.On("Logout", mock.Anything, accountID).Return(errors.New("logout failed"))

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req = req.WithContext(auth.WithAccountID(req.Context(), accountID))
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

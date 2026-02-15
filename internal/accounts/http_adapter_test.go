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

func TestHTTPAdapter_Signup(t *testing.T) {
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

	mockSvc.On("Signup", mock.Anything, account.Nickname, account.Email, "password").Return(account, tokenPair, nil)

	body := bytes.NewBufferString(`{"nickname":"nerd","email":"nerd@example.com","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/signup", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Signup(w, req)

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

func TestHTTPAdapter_Signup_BadPayload(t *testing.T) {
	mockSvc := new(MockAccountService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	body := bytes.NewBufferString(`{"nickname":123}`)
	req := httptest.NewRequest(http.MethodPost, "/signup", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHTTPAdapter_Signup_ServiceError(t *testing.T) {
	mockSvc := new(MockAccountService)
	logger := slog.Default()
	handler := NewHTTPAdapter(mockSvc, logger)

	mockSvc.On("Signup", mock.Anything, "nerd", "nerd@example.com", "password").Return(domain.Account{}, domain.TokenPair{}, errors.New("signup failed"))

	body := bytes.NewBufferString(`{"nickname":"nerd","email":"nerd@example.com","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/signup", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	mockSvc.AssertExpectations(t)
}

package accounts

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/kalogs-c/nerd-backlog/pkg/auth"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Login_Success(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	jwtManager := auth.NewJWTManager([]byte("nerd_secret"), time.Hour, time.Hour)
	svc := NewService(mockRepo, jwtManager)

	ctx := context.Background()
	email := "service_test@example.com"
	password := "password$123"
	hashedPassword, _ := auth.HashPassword(password)

	user := domain.Account{
		ID:             uuid.New(),
		Nickname:       "testuser",
		Email:          email,
		HashedPassword: hashedPassword,
		TimeStamps: domain.TimeStamps{
			InsertedAt: time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	accessToken, _ := jwtManager.GenerateAccessToken(user.ID)
	refreshToken, _, _ := jwtManager.GenerateRefreshToken(user.ID)

	mockRepo.On("GetAccountByEmail", ctx, email).Return(user, nil)
	mockRepo.On("StoreRefreshToken", ctx, user.ID, refreshToken, mock.AnythingOfType(reflect.TypeFor[time.Time]().Name())).Return(nil)

	account, tokenPair, err := svc.Login(ctx, email, password)
	require.NoError(t, err)
	require.Equal(t, user.Email, account.Email)
	require.Equal(t, accessToken, tokenPair.AccessToken)
	require.Equal(t, refreshToken, tokenPair.RefreshToken)
	mockRepo.AssertExpectations(t)
}

func TestService_Login_EmailNotFound(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	svc := NewService(mockRepo, auth.NewJWTManager([]byte("nerd_secret"), time.Hour, time.Hour))

	ctx := context.Background()
	email := "not_found@example.com"
	password := "password123"

	mockRepo.On("GetAccountByEmail", ctx, email).Return(domain.Account{}, domain.ErrAccountNotFound)

	_, _, err := svc.Login(ctx, email, password)
	require.Error(t, err)
	require.EqualError(t, err, domain.ErrAccountNotFound.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Login_PasswordIncorrect(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	svc := NewService(mockRepo, auth.NewJWTManager([]byte("nerd_secret"), time.Hour, time.Hour))

	ctx := context.Background()
	email := "service_test@example.com"
	password := "wrong_password"

	user := domain.Account{
		ID:             uuid.New(),
		Nickname:       "testuser",
		Email:          email,
		HashedPassword: "password$123",
		TimeStamps: domain.TimeStamps{
			InsertedAt: time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	mockRepo.On("GetAccountByEmail", ctx, email).Return(user, nil)

	_, _, err := svc.Login(ctx, email, password)
	require.Error(t, err)
	require.EqualError(t, err, domain.ErrAccountNotFound.Error())
	mockRepo.AssertExpectations(t)
}

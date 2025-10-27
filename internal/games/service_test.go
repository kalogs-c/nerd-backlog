package games

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestService_CreateGame(t *testing.T) {
	mockRepo := new(MockGameRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	gameIn := domain.Game{Title: "Hollow Knight"}
	gameOut := domain.Game{ID: uuid.New(), Title: "Hollow Knight"}

	mockRepo.On("CreateGame", ctx, gameIn).Return(gameOut, nil)

	got, err := svc.CreateGame(ctx, "Hollow Knight")
	require.NoError(t, err)
	require.Equal(t, gameOut.ID, got.ID)
	require.Equal(t, "Hollow Knight", got.Title)

	mockRepo.AssertExpectations(t)
}

func TestService_CreateGame_Error(t *testing.T) {
	mockRepo := new(MockGameRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	gameIn := domain.Game{Title: "Error Game"}
	mockRepo.On("CreateGame", ctx, gameIn).Return(domain.Game{}, errors.New("db error"))

	_, err := svc.CreateGame(ctx, "Error Game")
	require.Error(t, err)
	require.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

func TestService_GetGameByID(t *testing.T) {
	mockRepo := new(MockGameRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()
	id := uuid.New()

	want := domain.Game{ID: id, Title: "Zelda"}
	mockRepo.On("GetGameByID", ctx, id).Return(want, nil)

	got, err := svc.GetGameByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, want, got)
	mockRepo.AssertExpectations(t)
}

func TestService_GetGameByID_Error(t *testing.T) {
	mockRepo := new(MockGameRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()
	id := uuid.New()

	mockRepo.On("GetGameByID", ctx, id).Return(domain.Game{}, errors.New("not found"))

	_, err := svc.GetGameByID(ctx, id)
	require.Error(t, err)
	require.EqualError(t, err, "not found")
	mockRepo.AssertExpectations(t)
}

func TestService_ListGames(t *testing.T) {
	mockRepo := new(MockGameRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	want := []domain.Game{
		{ID: uuid.New(), Title: "Hades"},
		{ID: uuid.New(), Title: "Celeste"},
	}

	mockRepo.On("ListGames", ctx).Return(want, nil)

	got, err := svc.ListGames(ctx)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, "Hades", got[0].Title)
	mockRepo.AssertExpectations(t)
}

func TestService_ListGames_Error(t *testing.T) {
	mockRepo := new(MockGameRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	mockRepo.On("ListGames", ctx).Return([]domain.Game{}, errors.New("db error"))

	_, err := svc.ListGames(ctx)
	require.Error(t, err)
	require.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteGameByID_Success(t *testing.T) {
	mockRepo := new(MockGameRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()
	id := uuid.New()

	game := domain.Game{ID: id, Title: "To Delete"}
	mockRepo.On("GetGameByID", ctx, id).Return(game, nil)
	mockRepo.On("DeleteGameByID", ctx, id).Return(nil)

	err := svc.DeleteGameByID(ctx, id)
	require.NoError(t, err)

	mockRepo.AssertCalled(t, "GetGameByID", ctx, id)
	mockRepo.AssertCalled(t, "DeleteGameByID", ctx, id)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteGameByID_NotFound(t *testing.T) {
	mockRepo := new(MockGameRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()
	id := uuid.New()

	mockRepo.On("GetGameByID", ctx, id).Return(domain.Game{}, errors.New("not found"))

	err := svc.DeleteGameByID(ctx, id)
	require.Error(t, err)
	require.EqualError(t, err, "not found")

	mockRepo.AssertCalled(t, "GetGameByID", ctx, id)
	mockRepo.AssertNotCalled(t, "DeleteGameByID", ctx, id)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteGameByID_DeleteFails(t *testing.T) {
	mockRepo := new(MockGameRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()
	id := uuid.New()

	game := domain.Game{ID: id, Title: "Fail Delete"}
	mockRepo.On("GetGameByID", ctx, id).Return(game, nil)
	mockRepo.On("DeleteGameByID", ctx, id).Return(errors.New("delete error"))

	err := svc.DeleteGameByID(ctx, id)
	require.Error(t, err)
	require.EqualError(t, err, "delete error")
	mockRepo.AssertExpectations(t)
}

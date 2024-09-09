package auth_test

import (
	"context"
	"providerHub/internal/domain/dto"
	"providerHub/internal/infra/config"
	"providerHub/internal/service/auth"
	"providerHub/internal/service/auth/mocks"
	mock "providerHub/pkg/logger/mock"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	ctx := context.Background()
	tokenDTO := dto.Token{
		Email:    "test",
		Password: "test",
		JWT: config.JWT{
			AccessTTL:  time.Second * 5,
			RefreshTTL: time.Second * 15,
			Secret:     "test",
		},
	}

	tokenMock := mocks.NewInterface(t)
	loggerMock := mock.NewMockLogger()

	service := auth.New(loggerMock, tokenMock)

	_, err := service.Token(ctx, tokenDTO)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

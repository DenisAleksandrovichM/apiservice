package api

import (
	"context"
	mockPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

type userFixture struct {
	ctx      context.Context
	userRepo *mockPkg.MockInterface
	service  *Implementation
}

func userSetUp(t *testing.T) userFixture {
	t.Parallel()

	f := userFixture{}
	f.userRepo = mockPkg.NewMockInterface(gomock.NewController(t))
	f.service = New(f.userRepo)
	f.ctx = context.Background()
	return f
}

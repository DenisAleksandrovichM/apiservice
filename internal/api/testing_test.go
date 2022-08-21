package api

import (
	"context"
	"github.com/golang/mock/gomock"
	mockPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/mocks"
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

package user

import (
	"context"
	mockPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

type userFixture struct {
	ctx      context.Context
	userRepo *mockPkg.MockInterface
	service  *implementation
}

func userSetUp(t *testing.T) userFixture {
	t.Parallel()

	f := userFixture{}
	f.userRepo = mockPkg.NewMockInterface(gomock.NewController(t))
	f.service = New()
	return f
}

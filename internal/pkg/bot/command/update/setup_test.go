package add

import (
	"context"
	commandPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/command"
	mockPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

type userFixture struct {
	ctx      context.Context
	userRepo *mockPkg.MockInterface
	service  commandPkg.Interface
}

func userSetUp(t *testing.T) userFixture {
	t.Parallel()

	f := userFixture{}
	f.userRepo = mockPkg.NewMockInterface(gomock.NewController(t))
	f.service = New(f.userRepo)
	return f
}

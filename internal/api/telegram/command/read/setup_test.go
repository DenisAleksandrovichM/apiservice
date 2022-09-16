package read

import (
	"context"
	mockPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/core/user/mocks"
	commandPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/telegram/command"
	"github.com/golang/mock/gomock"
	"testing"
)

type userFixture struct {
	ctx      context.Context
	userRepo *mockPkg.MockInterface
	service  commandPkg.Command
}

func userSetUp(t *testing.T) userFixture {
	t.Parallel()

	f := userFixture{}
	f.userRepo = mockPkg.NewMockInterface(gomock.NewController(t))
	f.service = New(f.userRepo)
	return f
}

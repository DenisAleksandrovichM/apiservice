package grpc

import (
	"context"
	mockPkg "github.com/DenisAleksandrovichM/apiservice/internal/database/core/user/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

type userFixture struct {
	ctx      context.Context
	userRepo *mockPkg.MockUser
	service  *implementation
}

func userSetUp(t *testing.T) userFixture {
	t.Parallel()

	f := userFixture{}
	f.ctx = context.Background()
	f.userRepo = mockPkg.NewMockUser(gomock.NewController(t))
	f.service = New(f.userRepo)
	return f
}

package grpc

import (
	"context"
	"github.com/DenisAleksandrovichM/apiservice/pkg/api"
	modelsPkg "github.com/DenisAleksandrovichM/apiservice/pkg/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"testing"
)

var (
	errSome = errors.New("some error")
	user1   = modelsPkg.User{
		Login:     "test_login",
		FirstName: "test_fn",
		LastName:  "test_ln",
		Weight:    80,
		Height:    180,
		Age:       60,
	}
	user2 = modelsPkg.User{
		Login:     "test_login1",
		FirstName: "test_fn1",
		LastName:  "test_ln1",
		Weight:    81,
		Height:    181,
		Age:       61,
	}
)

func TestUserCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Create(f.ctx, user1).Return(user1, nil)
		req := &api.UserCreateRequest{
			Login:     user1.Login,
			FirstName: user1.FirstName,
			LastName:  user1.LastName,
			Weight:    float64(user1.Weight),
			Height:    uint32(user1.Height),
			Age:       uint32(user1.Age),
		}
		// act
		resp, err := f.service.UserCreate(f.ctx, req)
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, &api.UserCreateResponse{
			Login:     user1.Login,
			FirstName: user1.FirstName,
			LastName:  user1.LastName,
			Weight:    float64(user1.Weight),
			Height:    uint32(user1.Height),
			Age:       uint32(user1.Age),
		})
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Create(f.ctx, user1).Return(user1, errSome)
		// act
		_, err := f.service.UserCreate(f.ctx,
			&api.UserCreateRequest{
				Login:     user1.Login,
				FirstName: user1.FirstName,
				LastName:  user1.LastName,
				Weight:    float64(user1.Weight),
				Height:    uint32(user1.Height),
				Age:       uint32(user1.Age),
			})
		// assert
		assert.EqualError(t, err, status.Error(codes.Internal, errSome.Error()).Error())
	})
}

func TestUserRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Read(f.ctx, user1.Login).Return(user1, nil)
		// act
		resp, err := f.service.UserRead(f.ctx, &api.UserReadRequest{Login: user1.Login})
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, &api.UserReadResponse{
			Login:     user1.Login,
			FirstName: user1.FirstName,
			LastName:  user1.LastName,
			Weight:    float64(user1.Weight),
			Height:    uint32(user1.Height),
			Age:       uint32(user1.Age),
		})
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Read(f.ctx, user1.Login).Return(modelsPkg.User{}, errSome)
		// act
		_, err := f.service.UserRead(f.ctx, &api.UserReadRequest{Login: user1.Login})
		// assert
		assert.EqualError(t, err, status.Error(codes.Internal, errSome.Error()).Error())
	})
}

func TestUserUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Update(f.ctx, user1).Return(user1, nil)
		req := &api.UserUpdateRequest{
			Login:     user1.Login,
			FirstName: user1.FirstName,
			LastName:  user1.LastName,
			Weight:    float64(user1.Weight),
			Height:    uint32(user1.Height),
			Age:       uint32(user1.Age),
		}
		// act
		resp, err := f.service.UserUpdate(f.ctx, req)
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, &api.UserUpdateResponse{
			Login:     user1.Login,
			FirstName: user1.FirstName,
			LastName:  user1.LastName,
			Weight:    float64(user1.Weight),
			Height:    uint32(user1.Height),
			Age:       uint32(user1.Age),
		})
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Update(f.ctx, user1).Return(user1, errSome)
		req := &api.UserUpdateRequest{
			Login:     user1.Login,
			FirstName: user1.FirstName,
			LastName:  user1.LastName,
			Weight:    float64(user1.Weight),
			Height:    uint32(user1.Height),
			Age:       uint32(user1.Age),
		}
		// act
		_, err := f.service.UserUpdate(f.ctx, req)
		// assert
		assert.EqualError(t, err, status.Error(codes.Internal, errSome.Error()).Error())
	})
}

func TestUserDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Delete(f.ctx, user1.Login).Return(user1, nil)
		// act
		resp, err := f.service.UserDelete(f.ctx, &api.UserDeleteRequest{Login: user1.Login})
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, &api.UserDeleteResponse{})
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Delete(f.ctx, user1.Login).Return(user1, errSome)
		// act
		_, err := f.service.UserDelete(f.ctx, &api.UserDeleteRequest{Login: user1.Login})
		// assert
		assert.EqualError(t, err, status.Error(codes.Internal, errSome.Error()).Error())
	})
}

func TestUserList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		users := []modelsPkg.User{user1, user2}
		f.userRepo.EXPECT().List(f.ctx, map[string]interface{}{}).Return(users, nil)
		// act
		resp, err := f.service.UserList(f.ctx, &api.UserListRequest{})
		// assert
		respUsers := make([]*api.UserListResponse_User, 2, 2)
		respUsers[0] = &api.UserListResponse_User{
			Login:     user1.Login,
			FirstName: user1.FirstName,
			LastName:  user1.LastName,
			Weight:    float64(user1.Weight),
			Height:    uint32(user1.Height),
			Age:       uint32(user1.Age),
		}
		respUsers[1] = &api.UserListResponse_User{
			Login:     user2.Login,
			FirstName: user2.FirstName,
			LastName:  user2.LastName,
			Weight:    float64(user2.Weight),
			Height:    uint32(user2.Height),
			Age:       uint32(user2.Age),
		}
		require.NoError(t, err)
		assert.Equal(t, resp, &api.UserListResponse{Users: respUsers})
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().List(f.ctx, map[string]interface{}{}).Return(nil, errSome)
		// act
		_, err := f.service.UserList(context.Background(), &api.UserListRequest{})
		// assert
		assert.Equal(t, err.Error(), status.Error(codes.Internal, errSome.Error()).Error())
	})
}

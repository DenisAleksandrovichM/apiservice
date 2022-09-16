package api

import (
	"context"
	"github.com/DenisAleksandrovichM/apiservice/pkg/api"
	modelsPkg "github.com/DenisAleksandrovichM/apiservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	user1 = modelsPkg.User{
		Login:     "test_login",
		FirstName: "test_fn",
		LastName:  "test_ln",
		Weight:    80,
		Height:    180,
		Age:       60,
	}
)

func TestUserCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		f := userSetUp(t)
		resp, err := f.userClient.UserCreate(context.Background(), &api.UserCreateRequest{
			Login:     user1.Login,
			FirstName: user1.FirstName,
			LastName:  user1.LastName,
			Weight:    float64(user1.Weight),
			Height:    uint32(user1.Height),
			Age:       uint32(user1.Age),
		})
		assert.Nil(t, err)
		assert.Equal(t, resp.Login, user1.Login)
		assert.Equal(t, resp.FirstName, user1.FirstName)
		assert.Equal(t, resp.LastName, user1.LastName)
		assert.Equal(t, float32(resp.Weight), user1.Weight)
		assert.Equal(t, uint(resp.Height), user1.Height)
		assert.Equal(t, uint(resp.Age), user1.Age)

	})
}

func TestUserRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		f := userSetUp(t)
		resp, err := f.userClient.UserRead(context.Background(), &api.UserReadRequest{
			Login: user1.Login,
		})
		assert.Nil(t, err)
		assert.Equal(t, resp.Login, user1.Login)
		assert.Equal(t, resp.FirstName, user1.FirstName)
		assert.Equal(t, resp.LastName, user1.LastName)
		assert.Equal(t, float32(resp.Weight), user1.Weight)
		assert.Equal(t, uint(resp.Height), user1.Height)
		assert.Equal(t, uint(resp.Age), user1.Age)

	})
}

func TestUserUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		f := userSetUp(t)
		resp, err := f.userClient.UserUpdate(context.Background(), &api.UserUpdateRequest{
			Login:     user1.Login,
			FirstName: user1.FirstName,
			LastName:  user1.LastName,
			Weight:    float64(user1.Weight),
			Height:    uint32(user1.Height),
			Age:       uint32(user1.Age),
		})
		assert.Nil(t, err)
		assert.Equal(t, resp.Login, user1.Login)
		assert.Equal(t, resp.FirstName, user1.FirstName)
		assert.Equal(t, resp.LastName, user1.LastName)
		assert.Equal(t, float32(resp.Weight), user1.Weight)
		assert.Equal(t, uint(resp.Height), user1.Height)
		assert.Equal(t, uint(resp.Age), user1.Age)

	})
}

func TestUserDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		f := userSetUp(t)
		_, err := f.userClient.UserDelete(context.Background(), &api.UserDeleteRequest{
			Login: user1.Login,
		})
		assert.Nil(t, err)
	})
}

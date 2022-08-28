package user

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	modelsPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/models"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/validate"
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
	user2 = modelsPkg.User{
		Login:     "test_login1",
		FirstName: "test_fn1",
		LastName:  "test_ln1",
		Weight:    81,
		Height:    181,
		Age:       61,
	}
	errSome = errors.New("some error")
)

func TestUserCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Add(f.ctx, user1).Return(user1, nil)
		// act
		resp, err := f.service.Create(f.ctx, user1)
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, user1)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Run("validation login error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Create(f.ctx, modelsPkg.User{Login: ""})
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
		t.Run("validation first name error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Create(f.ctx, modelsPkg.User{
				Login:     user1.Login,
				FirstName: ""},
			)
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
		t.Run("validation last name error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Create(f.ctx, modelsPkg.User{
				Login:     user1.Login,
				FirstName: user1.FirstName,
				LastName:  ""},
			)
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
		t.Run("validation weight error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Create(f.ctx, modelsPkg.User{
				Login:     user1.Login,
				FirstName: user1.FirstName,
				LastName:  user1.LastName,
				Weight:    0},
			)
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
		t.Run("validation height error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Create(f.ctx, modelsPkg.User{
				Login:     user1.Login,
				FirstName: user1.FirstName,
				LastName:  user1.LastName,
				Weight:    user1.Weight,
				Height:    0},
			)
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
		t.Run("validation age error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Create(f.ctx, modelsPkg.User{
				Login:     user1.Login,
				FirstName: user1.FirstName,
				LastName:  user1.LastName,
				Weight:    user1.Weight,
				Height:    user1.Height,
				Age:       200},
			)
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
	})
}

func TestUserRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Read(f.ctx, user1.Login).Return(user1, nil)
		// act
		resp, err := f.service.Read(f.ctx, user1.Login)
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, user1)
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		// act
		_, err := f.service.Read(f.ctx, "")
		// assert
		assert.ErrorIs(t, err, validate.ErrValidation)
	})
}

func TestUserUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Update(f.ctx, user1).Return(user1, nil)
		// act
		resp, err := f.service.Update(f.ctx, user1)
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, user1)
	})

	t.Run("error", func(t *testing.T) {
		t.Run("validation login error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Update(f.ctx, modelsPkg.User{Login: ""})
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
		t.Run("validation first name error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Update(f.ctx, modelsPkg.User{
				Login:     user1.Login,
				FirstName: ""},
			)
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
		t.Run("validation last name error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Update(f.ctx, modelsPkg.User{
				Login:     user1.Login,
				FirstName: user1.FirstName,
				LastName:  ""},
			)
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
		t.Run("validation weight error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Update(f.ctx, modelsPkg.User{
				Login:     user1.Login,
				FirstName: user1.FirstName,
				LastName:  user1.LastName,
				Weight:    0},
			)
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
		t.Run("validation height error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Update(f.ctx, modelsPkg.User{
				Login:     user1.Login,
				FirstName: user1.FirstName,
				LastName:  user1.LastName,
				Weight:    user1.Weight,
				Height:    0},
			)
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
		t.Run("validation age error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Update(f.ctx, modelsPkg.User{
				Login:     user1.Login,
				FirstName: user1.FirstName,
				LastName:  user1.LastName,
				Weight:    user1.Weight,
				Height:    user1.Height,
				Age:       200},
			)
			// assert
			assert.ErrorIs(t, err, validate.ErrValidation)
		})
	})
}

func TestUserDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().Delete(f.ctx, user1.Login).Return(user1, nil)
		// act
		err := f.service.Delete(f.ctx, user1.Login)
		// assert
		require.NoError(t, err)
		//assert.Equal(t, resp, user1)
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		// act
		err := f.service.Delete(f.ctx, "")
		// assert
		assert.ErrorIs(t, err, validate.ErrValidation)
	})
}

func TestUserList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		users := []modelsPkg.User{user1, user2}
		f.userRepo.EXPECT().List(f.ctx, map[string]interface{}{}).Return(users, nil)
		// act
		resp, err := f.service.List(f.ctx, map[string]interface{}{})
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, users)
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		//users := []modelsPkg.User{user1, user2}
		f.userRepo.EXPECT().List(f.ctx, map[string]interface{}{}).Return([]modelsPkg.User{}, errSome)
		// act
		_, err := f.service.List(f.ctx, map[string]interface{}{})
		// assert
		assert.Equal(t, err, errSome)
	})
}

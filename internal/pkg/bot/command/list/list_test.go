package list

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	modelsPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/models"
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

func TestProcess(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		userStr := fmt.Sprintf("Login: %s, first name: %s, last name: %s,\nweight: %.2f, height: %d, age: %d",
			user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		f.userRepo.EXPECT().List(f.ctx, map[string]interface{}{}).Return([]modelsPkg.User{user1}, nil)
		f.userRepo.EXPECT().String(user1).Return(userStr)
		// act
		resp, err := f.service.Process(f.ctx, "test_login")
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, userStr)
	})

	t.Run("error", func(t *testing.T) {
		errSome := errors.New("some error")
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().List(f.ctx, map[string]interface{}{}).Return([]modelsPkg.User{}, errSome)
		// act
		_, err := f.service.Process(f.ctx, "")
		// assert
		assert.ErrorIs(t, err, errSome)
	})
}

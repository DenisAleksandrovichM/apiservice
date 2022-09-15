package delete

import (
	"fmt"
	modelsPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		f.userRepo.EXPECT().Delete(f.ctx, user1.Login).Return(user1, nil)
		userStr := fmt.Sprintf("Login: %s, first name: %s, last name: %s,\nweight: %.2f, height: %d, age: %d",
			user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		f.userRepo.EXPECT().String(user1).Return(userStr)
		// act
		resp, err := f.service.Process(f.ctx, "test_login")
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, userStr)
	})

	t.Run("error", func(t *testing.T) {
		t.Run("error validate args", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			// act
			_, err := f.service.Process(f.ctx, "test login")
			// assert
			assert.ErrorIs(t, err, errDelete)
		})

		t.Run("db create error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			f.userRepo.EXPECT().Delete(f.ctx, user1.Login).Return(modelsPkg.User{}, errDelete)
			// act
			_, err := f.service.Process(f.ctx, user1.Login)
			// assert
			assert.ErrorIs(t, err, errDelete)
		})
	})
}

package add

import (
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
		f.userRepo.EXPECT().Update(f.ctx, user1).Return(user1, nil)
		userStr := fmt.Sprintf("Login: %s, first name: %s, last name: %s,\nweight: %.2f, height: %d, age: %d",
			user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		f.userRepo.EXPECT().String(user1).Return(userStr)
		args := "test_login test_fn test_ln 80 180 60"
		// act
		resp, err := f.service.Process(f.ctx, args)
		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, userStr)
	})

	t.Run("error", func(t *testing.T) {
		t.Run("db update error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			f.userRepo.EXPECT().Update(f.ctx, user1).Return(modelsPkg.User{}, errUpdate)
			args := "test_login test_fn test_ln 80 180 60"
			// act
			_, err := f.service.Process(f.ctx, args)
			// assert
			assert.ErrorIs(t, err, errUpdate)
		})

		t.Run("error validate args", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			args := "test_login test_fn test_ln 80 180"
			// act
			_, err := f.service.Process(f.ctx, args)
			// assert
			assert.ErrorIs(t, err, errUpdate)
		})
	})

}

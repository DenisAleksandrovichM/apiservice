//go:generate mockgen -source ./user.go -destination=./mocks/user.go -package=mock_user
package user

import (
	"context"
	"fmt"
	cachePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/cache"
	localCachePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/cache/local"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/models"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/validate"
)

type Interface interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	Read(ctx context.Context, login string) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
	Delete(ctx context.Context, login string) error
	List(ctx context.Context, queryParams map[string]interface{}) ([]models.User, error)
	String(user models.User) string
}

type core struct {
	cache cachePkg.Interface
}

func New() *core {
	return &core{
		cache: localCachePkg.New(),
	}
}

func (c *core) Create(ctx context.Context, user models.User) (models.User, error) {
	if err := validate.ValidateUser(user); err != nil {
		return models.User{}, err
	}
	return c.cache.Add(ctx, user)
}

func (c *core) Read(ctx context.Context, login string) (models.User, error) {
	if err := validate.ValidateLogin(login); err != nil {
		return models.User{}, err
	}
	return c.cache.Read(ctx, login)
}

func (c *core) Update(ctx context.Context, user models.User) (models.User, error) {
	if err := validate.ValidateUser(user); err != nil {
		return models.User{}, err
	}
	return c.cache.Update(ctx, user)
}

func (c *core) Delete(ctx context.Context, login string) error {
	if err := validate.ValidateLogin(login); err != nil {
		return err
	}
	return c.cache.Delete(ctx, login)
}

func (c *core) List(ctx context.Context, queryParams map[string]interface{}) ([]models.User, error) {
	return c.cache.List(ctx, queryParams)
}

func (c *core) String(user models.User) string {
	return fmt.Sprintf("Login: %s, first name: %s, last name: %s,\nweight: %.2f, height: %d, age: %d",
		user.Login, user.FirstName, user.LastName, user.Weight, user.Height, user.Age)
}

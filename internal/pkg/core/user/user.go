package user

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	cachePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user/cache"
	localCachePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user/cache/local"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user/models"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user/validate"
)

type Interface interface {
	Create(ctx context.Context, user models.User) error
	Read(ctx context.Context, login string) (models.User, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, login string) error
	List(ctx context.Context, queryParams map[string]interface{}) ([]models.User, error)
	String(user models.User) string
}

func New(pool *pgxpool.Pool) Interface {
	return &core{
		cache: localCachePkg.New(pool),
	}
}

type core struct {
	cache cachePkg.Interface
}

func (c *core) Create(ctx context.Context, user models.User) error {
	if err := validate.ValidateUser(user); err != nil {
		return err
	}
	return c.cache.Add(ctx, user)
}

func (c *core) Read(ctx context.Context, login string) (models.User, error) {
	if err := validate.ValidateLogin(login); err != nil {
		return models.User{}, err
	}

	user, err := c.cache.Read(ctx, login)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (c *core) Update(ctx context.Context, user models.User) error {
	if err := validate.ValidateUser(user); err != nil {
		return err
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

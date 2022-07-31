package user

import (
	"fmt"
	cachePkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user/cache"
	localCachePkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user/cache/local"
	"gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user/models"
	"gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user/validate"
)

type Interface interface {
	Create(user models.User) error
	Read(login string) (models.User, error)
	Update(user models.User) error
	Delete(login string) error
	List() []models.User
	String(user models.User) string
}

func New() Interface {
	return &core{
		cache: localCachePkg.New(),
	}
}

type core struct {
	cache cachePkg.Interface
}

func (c *core) Create(user models.User) error {
	if err := validate.ValidateUser(user); err != nil {
		return err
	}
	return c.cache.Add(user)
}

func (c *core) Read(login string) (models.User, error) {
	if err := validate.ValidateLogin(login); err != nil {
		return models.User{}, err
	}

	user, err := c.cache.Read(login)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (c *core) Update(user models.User) error {
	if err := validate.ValidateUser(user); err != nil {
		return err
	}
	return c.cache.Update(user)
}

func (c *core) Delete(login string) error {
	if err := validate.ValidateLogin(login); err != nil {
		return err
	}
	return c.cache.Delete(login)
}

func (c *core) List() []models.User {
	return c.cache.List()
}

func (c *core) String(user models.User) string {
	return fmt.Sprintf("Login: %s, first name: %s, last name: %s,\nweight: %.2f, height: %d, age: %d",
		user.Login, user.FirstName, user.LastName, user.Weight, user.Height, user.Age)
}

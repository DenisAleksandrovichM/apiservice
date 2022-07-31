package local

import (
	"sync"

	"github.com/pkg/errors"
	cachePkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user/cache"
	"gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user/models"
)

const poolSize = 10

var (
	ErrUserNotExists = errors.New("user does not exist")
	ErrUserExists    = errors.New("user exists")
)

func New() cachePkg.Interface {
	return &cache{
		mu:     sync.RWMutex{},
		data:   map[string]models.User{},
		poolCh: make(chan struct{}, poolSize),
	}
}

type cache struct {
	mu     sync.RWMutex
	data   map[string]models.User
	poolCh chan struct{}
}

func (c *cache) List() []models.User {
	c.poolCh <- struct{}{}
	c.mu.RLock()
	defer func() {
		c.mu.RUnlock()
		<-c.poolCh
	}()

	result := make([]models.User, 0, len(c.data))
	for _, value := range c.data {
		result = append(result, value)
	}
	return result
}

func (c *cache) Add(user models.User) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	if _, ok := c.data[user.Login]; ok {
		return errors.Wrapf(ErrUserExists, "user-login: [%s]", user.Login)
	}
	c.data[user.Login] = user
	return nil
}

func (c *cache) Read(login string) (models.User, error) {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	if user, ok := c.data[login]; ok {
		return user, nil
	}

	return models.User{}, errors.Wrapf(ErrUserNotExists, "user-login: [%s]", login)
}

func (c *cache) Update(user models.User) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	if _, ok := c.data[user.Login]; !ok {
		return errors.Wrapf(ErrUserNotExists, "user-Login: [%s]", user.Login)
	}
	c.data[user.Login] = user
	return nil
}

func (c *cache) Delete(login string) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	if _, ok := c.data[login]; ok {
		delete(c.data, login)
		return nil
	}
	return errors.Wrapf(ErrUserNotExists, "user-login: [%s]", login)
}

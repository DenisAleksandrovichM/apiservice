package help

import (
	"context"
	"github.com/pkg/errors"
	commandPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user"
)

var errRead = errors.New("read process error")

func New(user userPkg.Interface) commandPkg.Interface {
	return &command{
		user: user,
	}
}

type command struct {
	user userPkg.Interface
}

func (c *command) Name() string {
	return "read"
}

func (c *command) Description() string {
	return "read user by id"
}

func (c *command) Process(ctx context.Context, args string) (string, error) {
	params, err := commandPkg.ValidateParams(args, 1)
	if err != nil {
		return "", errors.Wrap(errRead, err.Error())
	}
	login := params[0]
	user, err := c.user.Read(ctx, login)
	if err != nil {
		return "", errors.Wrap(errRead, err.Error())
	}
	return c.user.String(user), nil
}

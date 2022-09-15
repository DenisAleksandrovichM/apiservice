package read

import (
	"context"
	commandPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/command"
	userPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
	"github.com/pkg/errors"
)

var errRead = errors.New("read process error")

func New(user userPkg.User) commandPkg.Interface {
	return &command{
		user: user,
	}
}

type command struct {
	user userPkg.User
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

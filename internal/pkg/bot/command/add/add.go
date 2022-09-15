package add

import (
	"context"
	commandPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/command"
	userPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
	"github.com/pkg/errors"
)

var errAdd = errors.New("add process error")

func New(user userPkg.User) commandPkg.Interface {
	return &command{
		user: user,
	}
}

type command struct {
	user userPkg.User
}

func (c *command) Name() string {
	return "add"
}

func (c *command) Description() string {
	return "create user"
}

func (c *command) Process(ctx context.Context, args string) (string, error) {
	user, err := commandPkg.ProcessAddOrUpdate(ctx, args, c.user.Create)
	if err != nil {
		return "", errors.Wrap(errAdd, err.Error())
	}
	return c.user.String(user), nil
}

package add

import (
	"context"
	commandPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/command"
	userPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
	"github.com/pkg/errors"
)

var errUpdate = errors.New("update process error")

func New(user userPkg.User) commandPkg.Interface {
	return &command{
		user: user,
	}
}

type command struct {
	user userPkg.User
}

func (c *command) Name() string {
	return "update"
}

func (c *command) Description() string {
	return "update user"
}

func (c *command) Process(ctx context.Context, args string) (string, error) {
	user, err := commandPkg.ProcessAddOrUpdate(ctx, args, c.user.Update)
	if err != nil {
		return "", errors.Wrap(errUpdate, err.Error())
	}
	return c.user.String(user), nil
}

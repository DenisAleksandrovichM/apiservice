package add

import (
	"context"
	"github.com/pkg/errors"
	commandPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
)

var errUpdate = errors.New("update process error")

func New(user userPkg.Interface) commandPkg.Interface {
	return &command{
		user: user,
	}
}

type command struct {
	user userPkg.Interface
}

func (c *command) Name() string {
	return "update"
}

func (c *command) Description() string {
	return "update user"
}

func (c *command) Process(ctx context.Context, args string) (string, error) {
	msg, err := commandPkg.ProcessAddOrUpdate(ctx, args, c.user.Update)
	if err != nil {
		return "", errors.Wrap(errUpdate, err.Error())
	}
	return msg, nil
}

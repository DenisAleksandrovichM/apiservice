package add

import (
	commandPkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user"
)

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

func (c *command) Process(args string) (string, error) {
	msg, err := commandPkg.ProcessAddOrUpdate(args, c.user.Update)
	if err != nil {
		return "", err
	}
	return msg, nil
}

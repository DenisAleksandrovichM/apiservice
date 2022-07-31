package help

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
	return "read"
}

func (c *command) Description() string {
	return "read user by id"
}

func (c *command) Process(args string) (string, error) {
	params, err := commandPkg.ValidateParams(args, 1)
	if err != nil {
		return "", err
	}
	login := params[0]
	user, err := c.user.Read(login)
	if err != nil {
		return "", err
	}
	return c.user.String(user), nil
}

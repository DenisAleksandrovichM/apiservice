package add

import (
	"github.com/pkg/errors"
	commandPkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user"
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

func (c *command) Process(args string) (string, error) {
	msg, err := commandPkg.ProcessAddOrUpdate(args, c.user.Update)
	if err != nil {
		return "", errors.Wrap(errUpdate, err.Error())
	}
	return msg, nil
}

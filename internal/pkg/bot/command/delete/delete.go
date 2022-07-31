package add

import (
	"github.com/pkg/errors"
	commandPkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user"
	validatePkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user/validate"
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
	return "delete"
}

func (c *command) Description() string {
	return "delete user by id"
}

func (c *command) Process(args string) (string, error) {
	params, err := commandPkg.ValidateParams(args, 1)
	if err != nil {
		return "", err
	}
	login := params[0]
	if err := c.user.Delete(login); err != nil {
		if errors.Is(err, validatePkg.ErrValidation) {
			return "", errors.Wrap(commandPkg.BadArgument, err.Error())
		}
		return "", errors.Wrap(commandPkg.BadArgument, "internal error")
	}

	return "success", nil
}

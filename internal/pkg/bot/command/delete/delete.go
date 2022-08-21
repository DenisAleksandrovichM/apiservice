package delete

import (
	"context"
	"github.com/pkg/errors"
	commandPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
	validatePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/validate"
)

var errDelete = errors.New("delete process error")

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

func (c *command) Process(ctx context.Context, args string) (string, error) {
	params, err := commandPkg.ValidateParams(args, 1)
	if err != nil {
		return "", errors.Wrap(errDelete, err.Error())
	}
	login := params[0]
	user, err := c.user.Delete(ctx, login)
	if err != nil {
		if errors.Is(err, validatePkg.ErrValidation) {
			return "", errors.Wrap(errDelete, err.Error())
		}
		return "", errors.Wrap(errDelete, "internal error")
	}
	return c.user.String(user), nil
}

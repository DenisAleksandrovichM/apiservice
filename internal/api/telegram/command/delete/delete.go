package delete

import (
	"context"
	userPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/core/user"
	validatePkg "github.com/DenisAleksandrovichM/apiservice/internal/api/core/user/validate"
	commandPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/telegram/command"
	"github.com/pkg/errors"
)

var errDelete = errors.New("delete process error")

const (
	emptyResult   = ""
	correctResult = "request has been sent"
)

func New(user userPkg.User) commandPkg.Command {
	return &command{
		user: user,
	}
}

type command struct {
	user userPkg.User
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
		return emptyResult, errors.Wrap(errDelete, err.Error())
	}
	login := params[0]
	err = c.user.Delete(ctx, login)
	if err != nil {
		if errors.Is(err, validatePkg.ErrValidation) {
			return emptyResult, errors.Wrap(errDelete, err.Error())
		}
		return emptyResult, errors.Wrap(errDelete, "internal error")
	}
	return correctResult, nil
}

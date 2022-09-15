package list

import (
	"context"
	commandPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/command"
	userPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
	"strings"
)

func New(user userPkg.User) commandPkg.Interface {
	return &command{
		user: user,
	}
}

type command struct {
	user userPkg.User
}

func (c *command) Name() string {
	return "list"
}

func (c *command) Description() string {
	return "list of users"
}

func (c *command) Process(ctx context.Context, _ string) (string, error) {
	usersList, err := c.user.List(ctx, map[string]interface{}{})
	if err != nil {
		return "", err
	}
	if len(usersList) == 0 {
		return "no users", nil
	}
	result := make([]string, 0, len(usersList))
	for _, user := range usersList {
		result = append(result, c.user.String(user))
	}
	return strings.Join(result, "\n"), nil
}

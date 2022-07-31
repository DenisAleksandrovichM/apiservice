package help

import (
	commandPkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user"
	"strings"
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
	return "list"
}

func (c *command) Description() string {
	return "list of users"
}

func (c *command) Process(_ string) (string, error) {
	usersList := c.user.List()
	if len(usersList) == 0 {
		return "no users", nil
	}
	result := make([]string, 0, len(usersList))
	for _, user := range usersList {
		result = append(result, c.user.String(user))
	}
	return strings.Join(result, "\n"), nil
}

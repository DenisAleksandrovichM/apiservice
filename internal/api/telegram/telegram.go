package telegram

import (
	"context"
	"fmt"
	"github.com/DenisAleksandrovichM/apiservice/internal/api/config"
	commandPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/telegram/command"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Commander interface {
	Run(ctx context.Context) error
	RegisterHandler(cmd commandPkg.Command)
}

func MustNew() Commander {
	bot, err := tgbotapi.NewBotAPI(config.ApiKey)
	if err != nil {
		log.Fatal(errors.Wrap(err, "on init tgbot"))
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &implementation{
		bot:   bot,
		route: make(map[string]commandPkg.Command),
	}
}

type implementation struct {
	bot   *tgbotapi.BotAPI
	route map[string]commandPkg.Command
}

func (c *implementation) RegisterHandler(cmd commandPkg.Command) {
	c.route[cmd.Name()] = cmd
}

func (c *implementation) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if cmdName := update.Message.Command(); cmdName != "" {
			if cmd, ok := c.route[cmdName]; ok {
				text, err := cmd.Process(ctx, update.Message.CommandArguments())
				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = text
				}
			} else {
				msg.Text = "Unknown command"
			}
		} else {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg.Text = fmt.Sprintf("you send <%v>", update.Message.Text)
		}
		_, err := c.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "send tg message")
		}
	}
	return nil
}

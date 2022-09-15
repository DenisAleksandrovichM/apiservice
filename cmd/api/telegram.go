package main

import (
	"context"
	botPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot"
	cmdAddPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/add"
	cmdDeletePkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/delete"
	cmdHelpPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/help"
	cmdListPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/list"
	cmdReadPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/read"
	cmdUpdatePkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/update"
	userPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
	"github.com/pkg/errors"
)

func runBot(user userPkg.User, errSignals chan error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var bot botPkg.Interface
	{
		bot = botPkg.MustNew()

		commandAdd := cmdAddPkg.New(user)
		bot.RegisterHandler(commandAdd)

		commandDelete := cmdDeletePkg.New(user)
		bot.RegisterHandler(commandDelete)

		commandUpdate := cmdUpdatePkg.New(user)
		bot.RegisterHandler(commandUpdate)

		commandList := cmdListPkg.New(user)
		bot.RegisterHandler(commandList)

		commandRead := cmdReadPkg.New(user)
		bot.RegisterHandler(commandRead)

		commandHelp := cmdHelpPkg.New(map[string]string{
			commandAdd.Name():    commandAdd.Description(),
			commandDelete.Name(): commandDelete.Description(),
			commandUpdate.Name(): commandUpdate.Description(),
			commandRead.Name():   commandRead.Description(),
			commandList.Name():   commandList.Description(),
		})
		bot.RegisterHandler(commandHelp)
	}

	errSignals <- errors.Wrapf(bot.Run(ctx), "telegram client failure")
}

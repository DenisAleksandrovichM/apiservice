package main

import (
	"context"
	userPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/core/user"
	telegramPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/telegram"
	cmdAddPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/telegram/command/add"
	cmdDeletePkg "github.com/DenisAleksandrovichM/apiservice/internal/api/telegram/command/delete"
	cmdHelpPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/telegram/command/help"
	cmdListPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/telegram/command/list"
	cmdReadPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/telegram/command/read"
	cmdUpdatePkg "github.com/DenisAleksandrovichM/apiservice/internal/api/telegram/command/update"
	"github.com/pkg/errors"
)

func runTelegramBot(user userPkg.User, errSignals chan error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var commander telegramPkg.Commander
	{
		commander = telegramPkg.MustNew()

		commandAdd := cmdAddPkg.New(user)
		commander.RegisterHandler(commandAdd)

		commandDelete := cmdDeletePkg.New(user)
		commander.RegisterHandler(commandDelete)

		commandUpdate := cmdUpdatePkg.New(user)
		commander.RegisterHandler(commandUpdate)

		commandList := cmdListPkg.New(user)
		commander.RegisterHandler(commandList)

		commandRead := cmdReadPkg.New(user)
		commander.RegisterHandler(commandRead)

		commandHelp := cmdHelpPkg.New(map[string]string{
			commandAdd.Name():    commandAdd.Description(),
			commandDelete.Name(): commandDelete.Description(),
			commandUpdate.Name(): commandUpdate.Description(),
			commandRead.Name():   commandRead.Description(),
			commandList.Name():   commandList.Description(),
		})
		commander.RegisterHandler(commandHelp)
	}

	errSignals <- errors.Wrapf(commander.Run(ctx), "telegram client failure")
}

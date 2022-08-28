package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	botPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/add"
	cmdDeletePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/delete"
	cmdHelpPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/help"
	cmdListPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/list"
	cmdReadPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/read"
	cmdUpdatePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/update"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
)

func main() {
	log.SetLevel(log.DebugLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := userPkg.New()

	go runBot(user)
	go runREST(ctx)
	runGRPCServer(user)
}

func runBot(user userPkg.Interface) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
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

	if err := bot.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

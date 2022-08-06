package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	configPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/config"
	botPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/add"
	cmdDeletePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/delete"
	cmdHelpPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/help"
	cmdListPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/list"
	cmdReadPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/read"
	cmdUpdatePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/command/update"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configPkg.Host, configPkg.Port, configPkg.User, configPkg.Password, configPkg.DBname)
	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("ping database error", err)
	}

	config := pool.Config()
	config.MaxConnIdleTime = configPkg.MaxConnIdleTime
	config.MaxConnLifetime = configPkg.MaxConnLifetime
	config.MinConns = configPkg.MinConns
	config.MaxConns = configPkg.MaxConns

	user := userPkg.New(pool)

	go runBot(user)
	go runREST()
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

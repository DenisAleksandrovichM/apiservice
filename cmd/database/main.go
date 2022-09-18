package main

import (
	"github.com/DenisAleksandrovichM/apiservice/internal/database/config"
	userPkg "github.com/DenisAleksandrovichM/apiservice/internal/database/core/user"
	"github.com/DenisAleksandrovichM/apiservice/pkg/server"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	_ "net/http/pprof"
)

func main() {
	log.SetLevel(log.InfoLevel)
	db, err := createRepository()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	user := userPkg.New(db)
	errSignals := make(chan error)
	go server.RunHTTPServer(config.HTTPEndpoint, config.HTTPAddress, errSignals)
	go server.RunGRPCServer(user, config.GRPCNetwork, config.GRPCAddress, errSignals)
	go runConsumer(user, errSignals)

	err = <-errSignals
	log.Fatalf("Stopping service. Cause: %s", err.Error())
}

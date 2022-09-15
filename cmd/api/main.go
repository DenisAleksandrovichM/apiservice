package main

import (
	userPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
	log "github.com/sirupsen/logrus"
)

func main() {

	user := userPkg.New()
	errSignals := make(chan error)

	go runBot(user, errSignals)
	go runREST(errSignals)
	go runGRPCServer(user, errSignals)

	err := <-errSignals

	log.Fatalf("Stopping service. Cause: %s", err.Error())
}

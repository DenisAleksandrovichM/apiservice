package main

import (
	userPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/core/user"
	log "github.com/sirupsen/logrus"
)

func main() {

	user := userPkg.New()
	errSignals := make(chan error)

	go runTelegramBot(user, errSignals)
	go runHttpServer(errSignals)
	go runGRPCServer(user, errSignals)

	err := <-errSignals
	log.Fatalf("Stopping service. Cause: %s", err.Error())
}

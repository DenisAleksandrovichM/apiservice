package main

import (
	"github.com/DenisAleksandrovichM/apiservice/internal/api/config"
	userPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/core/user"
	"github.com/DenisAleksandrovichM/apiservice/pkg/server"
	log "github.com/sirupsen/logrus"
)

func main() {

	user := userPkg.New()
	errSignals := make(chan error)

	go runTelegramBot(user, errSignals)
	go server.RunHTTPServer(config.HTTPEndpoint, config.HTTPAddress, errSignals)
	go server.RunGRPCServer(user, config.GRPCNetwork, config.GRPCAddress, errSignals)

	err := <-errSignals
	log.Fatalf("Stopping service. Cause: %s", err.Error())
}

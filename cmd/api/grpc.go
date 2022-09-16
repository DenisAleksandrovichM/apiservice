package main

import (
	"github.com/DenisAleksandrovichM/apiservice/internal/api/config"
	userPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/core/user"
	apiPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/grpc"
	pb "github.com/DenisAleksandrovichM/apiservice/pkg/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

func runGRPCServer(user userPkg.User, errSignals chan error) {
	listener, err := net.Listen(config.GRPCNetwork, config.GRPCAddress)
	if err != nil {
		errSignals <- errors.Wrapf(err, "failed to create a %s listener at %s", config.GRPCNetwork, config.GRPCAddress)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServer(grpcServer, apiPkg.New(user))
	errSignals <- errors.Wrapf(grpcServer.Serve(listener), "grpc server failure")
}

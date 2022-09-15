package main

import (
	apiPkg "github.com/DenisAleksandrovichM/homework-1/internal/api"
	userPkg "github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
	pb "github.com/DenisAleksandrovichM/homework-1/pkg/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

const (
	network = "tcp"
	address = ":9081"
)

func runGRPCServer(user userPkg.User, errSignals chan error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		errSignals <- errors.Wrapf(err, "failed to create a %s listener at %s", network, address)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServer(grpcServer, apiPkg.New(user))
	errSignals <- errors.Wrapf(grpcServer.Serve(listener), "grpc server failure")
}

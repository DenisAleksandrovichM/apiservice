package server

import (
	userPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/core/user"
	apiPkg "github.com/DenisAleksandrovichM/apiservice/internal/api/grpc"
	pb "github.com/DenisAleksandrovichM/apiservice/pkg/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

func RunGRPCServer(user userPkg.User, network, address string, errSignals chan error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		errSignals <- errors.Wrapf(err, "failed to create a %s listener at %s", network, address)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServer(grpcServer, apiPkg.New(user))
	errSignals <- errors.Wrapf(grpcServer.Serve(listener), "grpc server failure")
}

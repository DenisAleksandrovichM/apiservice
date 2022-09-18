package server

import (
	"context"
	pb "github.com/DenisAleksandrovichM/apiservice/pkg/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func RunHTTPServer(endpoint, address string, errSignals chan error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterAdminHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		errSignals <- errors.Wrapf(err, "failed to register HTTP handler")
		return
	}

	errSignals <- errors.Wrapf(http.ListenAndServe(address, mux), "http server failure")
}

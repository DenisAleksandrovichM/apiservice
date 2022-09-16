package main

import (
	"context"
	"github.com/DenisAleksandrovichM/apiservice/internal/api/config"
	pb "github.com/DenisAleksandrovichM/apiservice/pkg/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func runHttpServer(errSignals chan error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterAdminHandlerFromEndpoint(ctx, mux, config.HTTPEndpoint, opts); err != nil {
		errSignals <- errors.Wrapf(err, "failed to register HTTP handler")
		return
	}

	errSignals <- errors.Wrapf(http.ListenAndServe(config.HTTPAddress, mux), "http server failure")
}

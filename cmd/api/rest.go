package main

import (
	"context"
	pb "github.com/DenisAleksandrovichM/homework-1/pkg/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

const (
	port8081 = ":9081"
	port8080 = ":9080"
)

func runREST(errSignals chan error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterAdminHandlerFromEndpoint(ctx, mux, port8081, opts); err != nil {
		errSignals <- errors.Wrapf(err, "failed to register HTTP handler")
		return
	}

	errSignals <- errors.Wrapf(http.ListenAndServe(port8080, mux), "http server failure")
}

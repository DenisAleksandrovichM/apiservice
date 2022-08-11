package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	apiPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/api"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
	pb "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

const (
	port8081 = ":8081"
	port8080 = ":8080"
)

func runGRPCServer(user userPkg.Interface) {
	listener, err := net.Listen("tcp", port8081)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServer(grpcServer, apiPkg.New(user))

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func runREST(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcherREST),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterAdminHandlerFromEndpoint(ctx, mux, port8081, opts); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(port8080, mux); err != nil {
		log.Fatal(err)
	}
}

func headerMatcherREST(key string) (string, bool) {
	switch key {
	case "Custom":
		return key, true
	default:
		return key, false
	}
}

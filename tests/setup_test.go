package tests

import (
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

const port = ":9081"

type userFixture struct {
	userClient api.AdminClient
}

func userSetUp(t *testing.T) userFixture {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	f := userFixture{api.NewAdminClient(conn)}
	return f
}

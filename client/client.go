package main

import (
	"context"
	"log"

	pb "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conns, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewAdminClient(conns)
	ctx := context.Background()

	createUser(client, ctx, "login1", "fn1", "ln1", 55, 155, 55)
	readUser(client, ctx, "login1")
	createUser(client, ctx, "login2", "fn2", "ln1", 65, 165, 65)
	createUser(client, ctx, "login3", "fn3", "ln3", 66, 166, 66)
	listUser(client, ctx)
	createUser(client, ctx, "login3", "fn3", "ln3", 66, 166, 66)
	updateUser(client, ctx, "login1", "fn11", "ln11", 55, 155, 55)
	listUser(client, ctx)
	deleteUser(client, ctx, "login3")
	listUser(client, ctx)
	deleteUser(client, ctx, "login5")
	listUser(client, ctx)
	updateUser(client, ctx, "login55", "fn11", "ln11", 55, 155, 55)
	listUser(client, ctx)
	readUser(client, ctx, "login3")

}

func createUser(client pb.AdminClient, ctx context.Context, login, firstName, lastName string, weight float32, height, age uint) {
	response, err := client.UserCreate(ctx, &pb.UserCreateRequest{
		Login:     login,
		FirstName: firstName,
		LastName:  lastName,
		Weight:    float64(weight),
		Height:    uint32(height),
		Age:       uint32(age),
	})

	if err != nil {
		log.Printf("response create err: [%v]", err)
	} else {
		log.Printf("response create: [%v]", response)
	}

}

func updateUser(client pb.AdminClient, ctx context.Context, login, firstName, lastName string, weight float32, height, age uint) {
	response, err := client.UserUpdate(ctx, &pb.UserUpdateRequest{
		Login:     login,
		FirstName: firstName,
		LastName:  lastName,
		Weight:    float64(weight),
		Height:    uint32(height),
		Age:       uint32(age),
	})

	if err != nil {
		log.Printf("response update err: [%v]", err)
	} else {
		log.Printf("response update: [%v]", response)
	}
}

func deleteUser(client pb.AdminClient, ctx context.Context, login string) {
	response, err := client.UserDelete(ctx, &pb.UserDeleteRequest{
		Login: login,
	})

	if err != nil {
		log.Printf("response delete err: [%v]", err)
	} else {
		log.Printf("response delete: [%v]", response)
	}
}

func readUser(client pb.AdminClient, ctx context.Context, login string) {
	response, err := client.UserRead(ctx, &pb.UserReadRequest{
		Login: login,
	})

	if err != nil {
		log.Printf("response read err: [%v]", err)
	} else {
		log.Printf("response read: [%v]", response)
	}
}

func listUser(client pb.AdminClient, ctx context.Context) {
	response, err := client.UserList(ctx, &pb.UserListRequest{})
	if err != nil {
		log.Printf("response list err: [%v]", err)
	} else {
		log.Printf("response list: [%v]", response)
	}
}

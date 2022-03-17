package main

import (
	"context"
	"log"

	authv1 "github.com/cory-evans/gps-tracker-authentication/pkg/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	client := authv1.NewAuthServiceClient(conn)

	resp, err := client.GetUser(context.Background(), &authv1.GetUserRequest{})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(resp)
}

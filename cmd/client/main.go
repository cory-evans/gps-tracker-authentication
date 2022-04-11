package main

import (
	"context"
	"log"

	"github.com/cory-evans/gps-tracker-authentication/pkg/auth"
	"github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	client := auth.NewAuthServiceClient(conn)

	ctx := context.Background()

	signInResp, err := client.SignIn(ctx, &auth.SignInRequest{
		Email:    "cory@email.localhost",
		Password: "this is a strong password",
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(signInResp.Token)

	ctx = metadata.AppendToOutgoingContext(context.Background(), jwtauth.JWT_METADATA_KEY, signInResp.Token)

	deviceCreateResp, err := client.CreateDevice(ctx, &auth.CreateDeviceRequest{
		DeviceName: "My Device",
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(deviceCreateResp)
}

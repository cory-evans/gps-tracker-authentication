package main

import (
	"log"
	"net"

	serverv1 "github.com/cory-evans/gps-tracker-authentication/internal/server/v1"
	authv1 "github.com/cory-evans/gps-tracker-authentication/pkg/auth/v1"
	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()

	authv1.RegisterAuthServiceServer(grpcServer, &serverv1.AuthService{})

	grpcServer.Serve(listen)
}

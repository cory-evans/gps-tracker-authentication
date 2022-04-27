package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	database "github.com/cory-evans/gps-tracker-authentication/internal/database"
	service "github.com/cory-evans/gps-tracker-authentication/internal/service"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	auth "go.buf.build/grpc/go/corux/gps-tracker-auth/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func myAuthFunc(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("could not get metadata")
	}

	log.Println(md)

	return ctx, nil
}

func main() {
	// get the port from the environment variable
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "8080"
	}
	listen, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalln(err)
	}

	mongoCtx := context.Background()
	db, err := database.NewDatabaseClient(os.Getenv("MONGO_URI"), mongoCtx)
	if err != nil {
		log.Fatalln(err)
	}

	userDB := db.Database("auth")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(myAuthFunc)),
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(myAuthFunc)),
	)
	authService := &service.AuthService{
		DB: userDB,
	}

	auth.RegisterAuthServiceServer(grpcServer, authService)

	go func() {
		err := grpcServer.Serve(listen)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	log.Printf("Server started on port 8080 %v\n", listen)

	// TODO: make this work for docker shutdown signals
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c // block until signal is received

	log.Println("Shutting down server...")
	grpcServer.Stop()
	listen.Close()
	db.Disconnect(mongoCtx)
}

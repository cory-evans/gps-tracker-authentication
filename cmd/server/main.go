package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	databasev1 "github.com/cory-evans/gps-tracker-authentication/internal/database/v1"
	servicev1 "github.com/cory-evans/gps-tracker-authentication/internal/service/v1"
	authv1 "github.com/cory-evans/gps-tracker-authentication/pkg/auth/v1"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
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
	listen, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalln(err)
	}

	mongoCtx := context.Background()
	db, err := databasev1.NewDatabaseClient(os.Getenv("MONGO_URI"), mongoCtx)
	if err != nil {
		log.Fatalln(err)
	}

	userDB := db.Database("auth")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(myAuthFunc)),
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(myAuthFunc)),
	)
	authService := &servicev1.AuthService{
		DB: userDB,
	}

	authv1.RegisterAuthServiceServer(grpcServer, authService)

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

package service

import (
	"context"
	"log"

	"github.com/cory-evans/gps-tracker-authentication/pkg/auth"
	"github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	auth.UnimplementedAuthServiceServer

	DB *mongo.Database
}

func (s *AuthService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	switch fullMethodName {
	case "/pkg.auth.AuthService/CreateUser", "/pkg.auth.AuthService/SignIn":
		// don't require auth
		return ctx, nil
	}

	// all other requests require auth
	ctx, err := jwtauth.MapJWT(ctx)
	if err != nil {
		log.Println("error mapping JWT to context metadata:", err)
		return ctx, status.Errorf(codes.Unauthenticated, "Not authenticated.")
	}

	log.Println("INFO: Authenticated request to", fullMethodName)

	return ctx, nil
}

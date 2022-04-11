package service

import (
	"context"
	"log"

	"github.com/cory-evans/gps-tracker-authentication/pkg/auth"
	"github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	auth.UnimplementedAuthServiceServer

	DB *mongo.Database
}

func (s *AuthService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	// map incomming JWT to context metadata
	ctx, err := jwtauth.MapJWT(ctx)
	if err != nil {
		log.Println("error mapping JWT to context metadata:", err)
	}

	return ctx, nil
}

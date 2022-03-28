package servicev1

import (
	"context"
	"log"

	authv1 "github.com/cory-evans/gps-tracker-authentication/pkg/auth/v1"
	jwtauthv1 "github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth/v1"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	authv1.UnimplementedAuthServiceServer

	DB *mongo.Database
}

func (s *AuthService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	// map incomming JWT to context metadata
	ctx, err := jwtauthv1.MapJWT(ctx)
	if err != nil {
		log.Println("error mapping JWT to context metadata:", err)
	}

	return ctx, nil
}

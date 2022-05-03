package service

import (
	"context"
	"log"

	"github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth"
	auth "go.buf.build/grpc/go/corux/gps-tracker-auth/auth/v1"
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
	case "/auth.v1.AuthService/CreateUser", "/auth.v1.AuthService/SignIn":
		// don't require auth
		return ctx, nil
	}

	// all other requests require auth
	ctx, err := jwtauth.MapJWT(ctx)
	if err != nil {
		log.Println("error mapping JWT to context metadata:", err)
		return ctx, status.Errorf(codes.Unauthenticated, "Not authenticated.")
	}

	sessionId := jwtauth.GetSessionIdFromContext(ctx)

	// make sure session still exists
	resp, err := s.SessionIsValid(ctx, &auth.SessionIsValidRequest{SessionId: sessionId})
	if err != nil {
		return ctx, status.Errorf(codes.Internal, "Internal Server Error: %s", err.Error())
	}

	if !resp.IsValid {
		return ctx, status.Errorf(codes.Unauthenticated, "Session expired or no longer exists.")
	}

	log.Println("INFO: Authenticated request to", fullMethodName)

	return ctx, nil
}

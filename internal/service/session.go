package service

import (
	"context"
	"time"

	"github.com/cory-evans/gps-tracker-authentication/internal/models"
	"github.com/cory-evans/gps-tracker-authentication/pkg/auth"
	"github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthService) SignIn(ctx context.Context, req *auth.SignInRequest) (*auth.SignInResponse, error) {
	users := s.DB.Collection("user")

	var user models.User
	err := users.FindOne(ctx, bson.M{"Email": req.Email}).Decode(&user)

	if err != nil {
		return nil, status.Error(codes.NotFound, "No user found")
	}

	// create a new session
	sessionId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	refreshToken, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	token, err := jwtauth.CreateNewSession(
		user.UserId,
		sessionId.String(),
		time.Hour*3,
	)

	if err != nil {
		return nil, err
	}

	return &auth.SignInResponse{
		Token:        token,
		RefreshToken: refreshToken.String(),
	}, nil
}

func (s *AuthService) SignOut(ctx context.Context, req *auth.SignOutRequest) (*auth.SignOutResponse, error) {
	return &auth.SignOutResponse{}, nil
}

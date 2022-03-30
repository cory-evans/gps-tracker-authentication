package servicev1

import (
	"context"
	"time"

	modelsv1 "github.com/cory-evans/gps-tracker-authentication/internal/models/v1"
	authv1 "github.com/cory-evans/gps-tracker-authentication/pkg/auth/v1"
	jwtauthv1 "github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth/v1"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthService) SignIn(ctx context.Context, req *authv1.SignInRequest) (*authv1.SignInResponse, error) {
	users := s.DB.Collection("user")

	var user modelsv1.User
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

	token, err := jwtauthv1.CreateNewSession(
		user.UserId,
		sessionId.String(),
		time.Hour*3,
	)

	if err != nil {
		return nil, err
	}

	return &authv1.SignInResponse{
		Token:        token,
		RefreshToken: refreshToken.String(),
	}, nil
}

func (s *AuthService) SignOut(ctx context.Context, req *authv1.SignOutRequest) (*authv1.SignOutResponse, error) {
	return &authv1.SignOutResponse{}, nil
}

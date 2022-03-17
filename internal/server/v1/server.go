package serverv1

import (
	"context"

	authv1 "github.com/cory-evans/gps-tracker-authentication/pkg/auth/v1"
)

type AuthService struct {
	authv1.AuthServiceServer
}

func (s *AuthService) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	return &authv1.GetUserResponse{
		User: &authv1.User{
			Id: "1",
		},
	}, nil
}

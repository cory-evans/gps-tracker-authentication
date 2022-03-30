package servicev1

import (
	"context"
	"fmt"

	modelsv1 "github.com/cory-evans/gps-tracker-authentication/internal/models/v1"
	authv1 "github.com/cory-evans/gps-tracker-authentication/pkg/auth/v1"
	jwtauthv1 "github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth/v1"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/metadata"
)

func (s *AuthService) GetUserByMongoID(ctx context.Context, id interface{}) *modelsv1.User {
	users := s.DB.Collection(modelsv1.USER_COLLECTION)

	var user modelsv1.User
	users.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user
}

func (s *AuthService) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	users := s.DB.Collection(modelsv1.USER_COLLECTION)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("could not get metadata")
	}

	userId := jwtauthv1.GetUserIdFromMetadata(md)

	if userId != req.Id {
		return nil, fmt.Errorf("Signed in user does not match requested user")
	}

	var user modelsv1.User
	users.FindOne(ctx, bson.M{"UserId": req.Id}).Decode(&user)
	return &authv1.GetUserResponse{
		User: user.AsProtoBuf(),
	}, nil
}

func (s *AuthService) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	users := s.DB.Collection(modelsv1.USER_COLLECTION)

	userId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	result, err := users.InsertOne(ctx, bson.M{
		"UserId":       userId.String(),
		"DisplayName":  req.DisplayName,
		"Email":        req.Email,
		"PasswordHash": req.Password,
	})
	if err != nil {
		return nil, err
	}

	user := s.GetUserByMongoID(ctx, result.InsertedID)

	return &authv1.CreateUserResponse{
		User: user.AsProtoBuf(),
	}, nil
}

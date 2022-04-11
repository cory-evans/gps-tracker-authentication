package service

import (
	"context"
	"fmt"

	"github.com/cory-evans/gps-tracker-authentication/internal/models"
	"github.com/cory-evans/gps-tracker-authentication/pkg/auth"
	"github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/metadata"
)

func (s *AuthService) GetUserByMongoID(ctx context.Context, id interface{}) *models.User {
	users := s.DB.Collection(models.USER_COLLECTION)

	var user models.User
	users.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user
}

func (s *AuthService) GetUser(ctx context.Context, req *auth.GetUserRequest) (*auth.GetUserResponse, error) {
	users := s.DB.Collection(models.USER_COLLECTION)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("could not get metadata")
	}

	userId := jwtauth.GetUserIdFromMetadata(md)

	if userId != req.Id {
		return nil, fmt.Errorf("Signed in user does not match requested user")
	}

	var user models.User
	users.FindOne(ctx, bson.M{"UserId": req.Id}).Decode(&user)
	return &auth.GetUserResponse{
		User: user.AsProtoBuf(),
	}, nil
}

func (s *AuthService) CreateUser(ctx context.Context, req *auth.CreateUserRequest) (*auth.CreateUserResponse, error) {
	users := s.DB.Collection(models.USER_COLLECTION)

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

	return &auth.CreateUserResponse{
		User: user.AsProtoBuf(),
	}, nil
}

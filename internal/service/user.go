package service

import (
	"context"

	"github.com/cory-evans/gps-tracker-authentication/internal/models"
	"github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth"
	"github.com/google/uuid"
	auth "go.buf.build/grpc/go/corux/gps-tracker-auth/auth/v1"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthService) getUserByMongoID(ctx context.Context, id interface{}) *models.User {
	users := s.DB.Collection(models.USER_COLLECTION)

	var user models.User
	users.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user
}

func (s AuthService) getUserByUserID(ctx context.Context, userId string) *models.User {
	users := s.DB.Collection(models.USER_COLLECTION)

	var user models.User
	users.FindOne(ctx, bson.M{"UserId": userId}).Decode(&user)
	return &user
}

func (s *AuthService) GetUser(ctx context.Context, req *auth.GetUserRequest) (*auth.GetUserResponse, error) {
	users := s.DB.Collection(models.USER_COLLECTION)

	userId := jwtauth.GetUserIdFromContext(ctx)

	if userId != req.Id {
		return nil, status.Errorf(codes.Unauthenticated, "Signed in user does not match requested user")
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

	user := s.getUserByMongoID(ctx, result.InsertedID)

	return &auth.CreateUserResponse{
		User: user.AsProtoBuf(),
	}, nil
}

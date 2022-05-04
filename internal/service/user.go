package service

import (
	"context"

	"github.com/cory-evans/gps-tracker-authentication/internal/models"
	"github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth"
	"github.com/google/uuid"
	auth "go.buf.build/grpc/go/corux/gps-tracker-auth/auth/v1"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
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

func (s *AuthService) GetMe(ctx context.Context, req *auth.GetMeRequest) (*auth.GetMeResponse, error) {
	userId := jwtauth.GetUserIdFromContext(ctx)

	user := s.getUserByUserID(ctx, userId)

	return &auth.GetMeResponse{
		User: user.AsProtoBuf(),
	}, nil
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

	// check that the user doesn't already exist
	var userModel models.User
	err = users.FindOne(ctx, bson.M{"Email": req.Email}).Decode(&userModel)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "User already exists")
	}

	if (req.DisplayName == "") || (req.Email == "") || (req.Password == "") {
		return nil, status.Errorf(codes.InvalidArgument, "Missing required fields")
	}

	if (req.FirstName == "") || (req.LastName == "") {
		return nil, status.Errorf(codes.InvalidArgument, "Missing required fields")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error hashing password")
	}

	result, err := users.InsertOne(ctx, bson.M{
		"UserId":       userId.String(),
		"DisplayName":  req.DisplayName,
		"Email":        req.Email,
		"PasswordHash": passwordHash,
		"FirstName":    req.FirstName,
		"LastName":     req.LastName,
	})
	if err != nil {
		return nil, err
	}

	user := s.getUserByMongoID(ctx, result.InsertedID)

	return &auth.CreateUserResponse{
		User: user.AsProtoBuf(),
	}, nil
}

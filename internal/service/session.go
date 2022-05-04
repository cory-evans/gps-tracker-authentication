package service

import (
	"context"
	"time"

	"github.com/cory-evans/gps-tracker-authentication/internal/models"
	"github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth"
	"github.com/google/uuid"
	auth "go.buf.build/grpc/go/corux/gps-tracker-auth/auth/v1"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *AuthService) createNewUserSession(ctx context.Context, deviceId string) (string, *models.Session, error) {
	sessionId, err := uuid.NewUUID()
	if err != nil {
		return "", nil, err
	}
	refreshToken, err := uuid.NewUUID()
	if err != nil {
		return "", nil, err
	}

	tokenExpiresAt := time.Now().Add(time.Hour * 24 * 7).UTC()
	token, err := jwtauth.CreateJWTSession(
		sessionId.String(),
		deviceId,
		tokenExpiresAt,
	)

	if err != nil {
		return "", nil, err
	}

	sess := models.Session{
		ID:           sessionId.String(),
		Subject:      deviceId,
		RefreshToken: refreshToken.String(),
		ExpiresAt:    tokenExpiresAt,
		Type:         models.USER_TYPE,
	}

	s.DB.Collection(models.SESSION_COLLECTION).InsertOne(ctx, sess)

	return token, &sess, nil
}
func (s *AuthService) createNewDeviceSession(ctx context.Context, userId string) (string, *models.Session, error) {
	sessionId, err := uuid.NewUUID()
	if err != nil {
		return "", nil, err
	}
	refreshToken, err := uuid.NewUUID()
	if err != nil {
		return "", nil, err
	}

	tokenExpiresAt := time.Now().Add(time.Hour * 24 * 7 * 365).UTC()
	token, err := jwtauth.CreateJWTSession(
		sessionId.String(),
		userId,
		tokenExpiresAt,
	)

	if err != nil {
		return "", nil, err
	}

	sess := models.Session{
		ID:           sessionId.String(),
		Subject:      userId,
		RefreshToken: refreshToken.String(),
		ExpiresAt:    tokenExpiresAt,
		Type:         models.DEVICE_TYPE,
	}

	s.DB.Collection(models.SESSION_COLLECTION).InsertOne(ctx, sess)

	return token, &sess, nil
}

func (s *AuthService) getUserSessionByID(ctx context.Context, sessionId string) (*models.Session, error) {
	session := models.Session{}
	err := s.DB.Collection(models.SESSION_COLLECTION).FindOne(ctx, bson.M{"ID": sessionId}).Decode(&session)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *AuthService) SessionIsValid(ctx context.Context, req *auth.SessionIsValidRequest) (*auth.SessionIsValidResponse, error) {
	session := models.Session{}
	err := s.DB.Collection(models.SESSION_COLLECTION).FindOne(ctx, bson.M{"SessionId": req.GetSessionId()}).Decode(&session)

	if err != nil {
		return &auth.SessionIsValidResponse{
			IsValid: false,
		}, nil
	}

	isValid := session.ExpiresAt.Unix() > time.Now().UTC().Unix()

	return &auth.SessionIsValidResponse{
		IsValid: isValid,
	}, nil
}

func (s *AuthService) SignIn(ctx context.Context, req *auth.SignInRequest) (*auth.SignInResponse, error) {
	users := s.DB.Collection("user")

	var user models.User
	err := users.FindOne(ctx, bson.M{"Email": req.Email}).Decode(&user)

	if err != nil {
		return nil, status.Error(codes.NotFound, "No user found")
	}

	if user.PasswordHash != req.Password {
		return nil, status.Error(codes.Unauthenticated, "Incorrect password")
	}

	// create a new session
	token, sess, err := s.createNewUserSession(ctx, user.UserId)
	if err != nil {
		return nil, err
	}

	return &auth.SignInResponse{
		Session: &auth.UserSession{
			Token:        token,
			RefreshToken: sess.RefreshToken,
			ExpiresAt:    sess.ExpiresAt.Format(time.RFC3339),
			User:         user.AsProtoBuf(),
		},
	}, nil
}

func (s *AuthService) SignOut(ctx context.Context, req *auth.SignOutRequest) (*auth.SignOutResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing metadata")
	}

	sessionId := jwtauth.GetSessionIdFromMetadata(md)

	if sessionId == "" {
		return nil, status.Error(codes.Unauthenticated, "missing session id")
	}

	s.DB.Collection(models.SESSION_COLLECTION).DeleteOne(ctx, bson.M{"ID": sessionId})

	return &auth.SignOutResponse{}, nil
}

func (s *AuthService) InvalidateAllSessions(ctx context.Context, req *auth.InvalidateAllSessionsRequest) (*auth.InvalidateAllSessionsResponse, error) {
	userId := jwtauth.GetUserIdFromContext(ctx)

	s.DB.Collection(models.SESSION_COLLECTION).DeleteMany(ctx, bson.M{"SubjectId": userId})

	return &auth.InvalidateAllSessionsResponse{}, nil
}

func (s *AuthService) RefreshSession(ctx context.Context, req *auth.RefreshSessionRequest) (*auth.RefreshSessionResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing metadata")
	}

	sessionId := jwtauth.GetSessionIdFromMetadata(md)

	if sessionId == "" {
		return nil, status.Error(codes.Unauthenticated, "missing session id")
	}

	oldSession := models.Session{}
	err := s.DB.Collection(models.SESSION_COLLECTION).FindOneAndDelete(ctx, bson.M{"ID": sessionId}).Decode(&oldSession)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid session id")
	}

	token, sess, err := s.createNewUserSession(ctx, oldSession.Subject)
	if err != nil {
		return nil, err
	}

	user := s.getUserByUserID(ctx, oldSession.Subject)

	return &auth.RefreshSessionResponse{
		Session: &auth.UserSession{
			Token:        token,
			RefreshToken: sess.RefreshToken,
			ExpiresAt:    sess.ExpiresAt.Format(time.RFC3339),
			User:         user.AsProtoBuf(),
		},
	}, nil
}

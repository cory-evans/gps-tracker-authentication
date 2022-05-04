package jwtauth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/metadata"
)

func keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET")), nil
}

func MapJWT(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, fmt.Errorf("could not get metadata from context")
	}

	subject := md.Get(JWT_METADATA_SUB_KEY)
	if len(subject) == 1 {
		// the metadata is valid already
		return ctx, nil
	}

	tokenString, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return ctx, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, keyFunc)

	if err != nil {
		return ctx, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return ctx, fmt.Errorf("could not get claims")
	}

	md.Set(JWT_METADATA_SUB_KEY, claims.Subject)
	md.Set(JWT_METADATA_ID_KEY, claims.Id)

	return metadata.NewIncomingContext(ctx, md), nil
}

func CreateJWTSession(id, subject string, expiresAt time.Time) (string, error) {
	now := time.Now()
	claims := &jwt.StandardClaims{
		Id:        id,
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: expiresAt.Unix(),
		Subject:   subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingKey, err := keyFunc(token)

	ss, err := token.SignedString(signingKey)

	return ss, err
}

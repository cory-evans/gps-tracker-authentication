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
	// map incomming JWT to context metadata

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

	return metadata.NewIncomingContext(ctx, md), nil
}

func CreateJWTSession(id, subject string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &jwt.StandardClaims{
		Id:        id,
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(duration).Unix(),
		Subject:   subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingKey, err := keyFunc(token)

	ss, err := token.SignedString(signingKey)

	return ss, err
}

package jwtauthv1

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

func keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET")), nil
}

func MapJWT(ctx context.Context) (context.Context, error) {
	// map incomming JWT to context metadata

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, fmt.Errorf("could not get metadata")
	}

	log.Println("MapJWT md: ", md)

	tokenStrings := md.Get(JWT_METADATA_KEY)
	if len(tokenStrings) != 1 {
		return ctx, fmt.Errorf("could not get token")
	}

	tokenString := tokenStrings[0]

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

func CreateNewSession(userId, sessionId string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &jwt.StandardClaims{
		Id:        sessionId,
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(duration).Unix(),
		Subject:   userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingKey, err := keyFunc(token)

	ss, err := token.SignedString(signingKey)

	return ss, err
}

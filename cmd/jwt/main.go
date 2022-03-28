package main

import (
	"context"
	"fmt"
	"log"

	jwtauthv1 "github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth/v1"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()

	md := metadata.New(map[string]string{
		jwtauthv1.JWT_METADATA_KEY: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0IiwiZXhwIjoxNjQ4NTA2MzEzLCJpc3MiOiJ0ZXN0In0.B0QC-T9kCV2l37Vo5gL91FE0hEIFrg7Jbo9ONWKUvMg",
	})
	newCtx := metadata.NewIncomingContext(ctx, md)

	jwtauthv1.MapJWT(newCtx)

	log.Println(newCtx)
}

func createJWT() {
	signingKey := "secret"

	claims := &jwt.StandardClaims{
		ExpiresAt: jwt.TimeFunc().Unix() + 3600,
		Issuer:    "test",
		Audience:  "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(signingKey))
	fmt.Printf("%v %v\n", ss, err)
}

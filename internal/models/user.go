package models

import (
	auth "go.buf.build/grpc/go/corux/gps-tracker-auth/auth/v1"
)

const (
	USER_COLLECTION = "user"
)

type User struct {
	UserId       string `json:"user_id" bson:"UserId"`
	DisplayName  string `json:"display_name" bson:"DisplayName"`
	Email        string `json:"email" bson:"Email"`
	FirstName    string `json:"first_name" bson:"FirstName"`
	LastName     string `json:"last_name" bson:"LastName"`
	PasswordHash []byte
}

func (u *User) AsProtoBuf() *auth.User {
	return &auth.User{
		UserId:      u.UserId,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
	}
}

package modelsv1

import (
	authv1 "github.com/cory-evans/gps-tracker-authentication/pkg/auth/v1"
)

const (
	USER_COLLECTION = "user"
)

type User struct {
	UserId       string `json:"user_id" bson:"UserId"`
	DisplayName  string `json:"display_name" bson:"DisplayName"`
	Email        string `json:"email" bson:"Email"`
	PasswordHash string
}

func (u *User) AsProtoBuf() *authv1.User {
	return &authv1.User{
		UserId:      u.UserId,
		DisplayName: u.DisplayName,
		Email:       u.Email,
	}
}

package models

import (
	"github.com/cory-evans/gps-tracker-authentication/pkg/auth"
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

func (u *User) AsProtoBuf() *auth.User {
	return &auth.User{
		UserId:      u.UserId,
		DisplayName: u.DisplayName,
		Email:       u.Email,
	}
}

package modelsv1

import (
	authv1 "github.com/cory-evans/gps-tracker-authentication/pkg/auth/v1"
)

type Session struct {
	UserId       string `json:"user_id" bson:"UserId"`
	SessionId    string `json:"session_id" bson:"SessionId"`
	RefreshToken string `json:"refresh_token" bson:"RefreshToken"`
	ExpiresAt    int64  `json:"expires" bson:"Expires"`
}

func (s *Session) AsProtoBuf() *authv1.UserSession {
	return &authv1.UserSession{
		UserId:       s.UserId,
		SessionId:    s.SessionId,
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,
	}
}

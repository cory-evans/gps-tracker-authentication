package models

import "time"

const (
	SESSION_COLLECTION = "session"
	USER_TYPE          = "user"
	DEVICE_TYPE        = "device"
)

type Session struct {
	ID           string    `json:"session_id" bson:"SessionId"`
	Subject      string    `json:"subject_id" bson:"SubjectId"`
	RefreshToken string    `json:"refresh_token" bson:"RefreshToken"`
	ExpiresAt    time.Time `json:"expires_at" bson:"ExpiresAt"`
	Type         string    `json:"type" bson:"Type"`
}

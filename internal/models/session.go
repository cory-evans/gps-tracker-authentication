package models

const (
	SESSION_COLLECTION = "session"
	USER_TYPE          = "user"
	DEVICE_TYPE        = "device"
)

type Session struct {
	ID           string `json:"session_id" bson:"SessionId"`
	Subject      string `json:"subject_id" bson:"SubjectId"`
	RefreshToken string `json:"refresh_token" bson:"RefreshToken"`
	ExpiresAtUtc int64  `json:"expires_at_utc" bson:"ExpiresAtUtc"`
	Type         string `json:"type" bson:"Type"`
}

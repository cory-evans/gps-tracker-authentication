package models

const (
	USER_SESSION_COLLECTION   = "user_session"
	DEVICE_SESSION_COLLECTION = "device_session"
)

type Session struct {
	ID           string `json:"session_id" bson:"SessionId"`
	Subject      string `json:"subject_id" bson:"SubjectId"`
	RefreshToken string `json:"refresh_token" bson:"RefreshToken"`
}

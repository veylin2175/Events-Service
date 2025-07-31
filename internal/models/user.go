package models

type User struct {
	UserID int64   `json:"user_id"`
	Events []Event `json:"events"`
}

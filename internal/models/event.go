package models

import "time"

type Event struct {
	EventID int64     `json:"event_id"`
	Date    time.Time `json:"date"`
	Text    string    `json:"text"`
}

package models

import "time"

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Log       string    `json:"log"`
}

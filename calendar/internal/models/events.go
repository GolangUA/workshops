package models

import "time"

type Event struct {
	ID          string
	Title       string
	Description string
	TimeFrom    time.Time
	TimeTo      time.Time
	Notes       []string
}

package models

import (
	"time"
)

type Event struct {
	ID          string
	Title       string
	Description string
	TimeFrom    time.Time
	TimeTo      time.Time
	Notes       []string
}

func (e Event) WithTimezone(timezone string) (*Event, error) {
	var locFrom, locTo *time.Location
	if timezone == "" {
		locFrom = e.TimeFrom.Location()
		locTo = e.TimeTo.Location()
	} else if loc, err := time.LoadLocation(timezone); err != nil {
		return nil, err
	} else {
		locFrom = loc
		locTo = loc
	}

	return &Event{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		TimeFrom:    e.TimeFrom.In(locFrom),
		TimeTo:      e.TimeTo.In(locTo),
		Notes:       e.Notes,
	}, nil
}

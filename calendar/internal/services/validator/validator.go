package validator

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/api"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/calendar"
	"log"
	"strings"
	"time"
)

// Service validates structures
type Service struct {
}

type GetEvents struct {
	Title    string
	Timezone string
	DateFrom string
	TimeFrom string
	DateTo   string
	TimeTo   string
}

type CreateEvent struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Time        string   `json:"time"`
	Timezone    string   `json:"timezone"`
	Duration    int      `json:"duration"`
	Notes       []string `json:"notes"`
}

type UpdateEvent struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Time        string   `json:"time"`
	Timezone    string   `json:"timezone"`
	Duration    int      `json:"duration"`
	Notes       []string `json:"notes"`
}

func (s *Service) Validate(v interface{}) error {
	errors := make([]string, 0, 10)
	switch t := v.(type) {
	case *GetEvents:
		if _, err := time.LoadLocation(t.Timezone); err != nil {
			errors = append(errors, "timezone must be a valid timezone")
		}
		if _, err := time.Parse(calendar.DateLayout, t.DateFrom); t.DateFrom != "" && err != nil {
			errors = append(errors, "date_from must be a valid date")
		}
		if _, err := time.Parse(calendar.TimeLayout, t.TimeFrom); t.TimeFrom != "" && err != nil {
			errors = append(errors, "time_from must be a valid time")
		}
		if _, err := time.Parse(calendar.DateLayout, t.DateTo); t.DateTo != "" && err != nil {
			errors = append(errors, "date_to must be a valid date")
		}
		if _, err := time.Parse(calendar.TimeLayout, t.TimeTo); t.TimeTo != "" && err != nil {
			errors = append(errors, "time_to must be a valid time")
		}
		if t.Timezone != "" {
			if (t.DateFrom != "" && t.TimeFrom == "") || (t.DateFrom == "" && t.TimeFrom != "") {
				errors = append(errors, "both date_from and time_from should be present if timezone is specified")
			}
			if (t.DateTo != "" && t.TimeTo == "") || (t.DateTo == "" && t.TimeTo != "") {
				errors = append(errors, "both date_to and time_to should be present if timezone is specified")
			}
		}
	case *CreateEvent:
		errors = validateEventModification(errors, t.Title, t.Time, t.Timezone, t.Duration)
	case *UpdateEvent:
		if strings.TrimSpace(t.ID) == "" {
			errors = append(errors, "id must not be blank")
		}
		errors = validateEventModification(errors, t.Title, t.Time, t.Timezone, t.Duration)
	case *api.UserTimezone:
		if strings.TrimSpace(t.Username) == "" {
			errors = append(errors, "username must not be blank")
		}
		if strings.TrimSpace(t.Timezone) == "" {
			errors = append(errors, "timezone must not be blank")
		} else if _, err := time.LoadLocation(t.Timezone); err != nil {
			errors = append(errors, "timezone must be a valid timezone")
		}
	default:
		log.Printf("Validation of type %T is not supported\n", t)
		return nil
	}
	if len(errors) > 0 {
		return fmt.Errorf("validation errors: [%s]", strings.Join(errors, "; "))
	}
	return nil
}

func validateEventModification(errors []string, title, timeVal, timezone string, duration int) []string {
	if strings.TrimSpace(title) == "" {
		errors = append(errors, "title must not be blank")
	}
	if _, err := time.Parse(calendar.DateTimeLayout, timeVal); err != nil {
		errors = append(errors, "invalid time value")
	}
	if _, err := time.LoadLocation(timezone); err != nil || strings.TrimSpace(timezone) == "" {
		errors = append(errors, "timezone must not be blank and be a valid timezone")
	}
	if duration <= 0 {
		errors = append(errors, "duration must be greater than 0")
	}
	return errors
}

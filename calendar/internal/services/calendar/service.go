package calendar

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"time"
)

const DateTimeLayout = "2006-01-02 15:04"
const DateLayout = "2006-01-02"
const TimeLayout = "15:04"

type Repository interface {
	// Events
	GetEvents(username, title, dateFrom, timeFrom, dateTo, timeTo string) ([]*models.Event, error)
	GetEvent(id string) (*models.Event, error)
	GetEventsCount() (int, error)
	EventExists(id string) (bool, error)
	CreateEvent(username string, title, description string, from time.Time, to time.Time, notes []string) (*models.Event, error)
	UpdateEvent(id, title, description string, from time.Time, to time.Time, notes []string) (*models.Event, error)
	DeleteEvent(id string) (bool, error)
	// Users
	GetUser(username string) (*models.User, error)
	GetUsersCount() (int, error)
	UpdateUserTimezone(username, timezone string) (*models.User, error)
	// User's event
	GetEventOwner(eventId string) (string, error)
}

// Service holds calendar business logic and works with repository
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetEvents(username, title, dateFrom, timeFrom, dateTo, timeTo, timezone string) ([]*models.Event, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return nil, err
	}
	userTimezone := user.Timezone

	if timezone == "" {
		timezone = userTimezone
	}

	if dateFrom != "" {
		if timeFrom == "" {
			timeFrom = "00:00"
		}
		convertedDate, convertedTime, err := normalizeDateTime(dateFrom, timeFrom, timezone)
		if err != nil {
			return nil, fmt.Errorf("convert date=\"%s\" time=\"%s\" to timezone=\"%s\": %v", dateFrom, timeFrom, timezone, err)
		}
		dateFrom = convertedDate
		timeFrom = convertedTime
	}
	if dateTo != "" {
		if timeTo == "" {
			timeTo = "23:59"
		}
		convertedDate, convertedTime, err := normalizeDateTime(dateTo, timeTo, timezone)
		if err != nil {
			return nil, fmt.Errorf("convert date=\"%s\" time=\"%s\" to timezone=\"%s\": %v", dateTo, timeTo, timezone, err)
		}
		dateTo = convertedDate
		timeTo = convertedTime
	}

	events, err := s.repo.GetEvents(username, title, dateFrom, timeFrom, dateTo, timeTo)
	if err != nil {
		return nil, err
	}

	res := make([]*models.Event, 0, len(events))
	for _, e := range events {
		if conv, err := e.WithTimezone(userTimezone); err != nil {
			return nil, fmt.Errorf("convert event with ID=\"%s\" to timezone=\"%s\": %v", e.ID, timezone, err)
		} else {
			res = append(res, conv)
		}
	}

	return res, err
}

func (s *Service) GetEvent(id string) (*models.Event, error) {
	return s.repo.GetEvent(id)
}

func (s *Service) GetEventOwner(id string) (string, error) {
	return s.repo.GetEventOwner(id)
}

func (s *Service) GetEventsCount() (int, error) {
	return s.repo.GetEventsCount()
}

func (s *Service) CreateEvent(username string, title, description, timeVal, timezone string, duration time.Duration, notes []string) (*models.Event, error) {
	timeFrom, timeTo, err := timeFromTo(timeVal, timezone, duration)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateEvent(username, title, description, *timeFrom, *timeTo, notes)
}

func (s *Service) UpdateEvent(id, title, description, timeVal, timezone string, duration time.Duration, notes []string) (*models.Event, error) {
	if ok, err := s.repo.EventExists(id); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	timeFrom, timeTo, err := timeFromTo(timeVal, timezone, duration)
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateEvent(id, title, description, *timeFrom, *timeTo, notes)
}

func (s *Service) DeleteEvent(id string) (bool, error) {
	return s.repo.DeleteEvent(id)
}

func (s *Service) UpdateUserTimezone(username, timezone string) (*models.User, error) {
	return s.repo.UpdateUserTimezone(username, timezone)
}

func (s *Service) GetUsersCount() (int, error) {
	return s.repo.GetUsersCount()
}

func timeFromTo(timeVal, timezone string, duration time.Duration) (*time.Time, *time.Time, error) {
	l, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid location \"%s\": %v", timezone, err)
	}
	if duration <= 0 {
		return nil, nil, fmt.Errorf("duration must be greateer than 0")
	}
	timeFrom, err := time.ParseInLocation(DateTimeLayout, timeVal, l)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid datetime \"%s\": %v", timeVal, err)
	}
	timeTo := timeFrom.Add(duration)
	return &timeFrom, &timeTo, nil
}

var dbLocation = time.Local

func normalizeDateTime(date string, timev string, timezone string) (string, string, error) {
	if date == "" && timev == "" {
		return "", "", nil
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", "", fmt.Errorf("invalid timezone=\"%s\": %v", timezone, err)
	}

	if loc == nil {
		return date, timev, nil
	}

	dateTime := fmt.Sprintf("%s %s", date, timev)
	zoned, err := time.ParseInLocation(DateTimeLayout, dateTime, loc)
	if err != nil {
		return "", "", fmt.Errorf("convert datetime=\"%s\": %v", dateTime, err)
	}
	converted := zoned.In(dbLocation)

	return converted.Format(DateLayout), converted.Format(TimeLayout), nil
}

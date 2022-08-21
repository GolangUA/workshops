package grpc

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/Roma7-7-7/workshops/calendar/proto"
	"time"
)

type Validator interface {
	Validate(interface{}) error
}

type Service interface {
	GetEvents(username, title, dateFrom, timeFrom, dateTo, timeTo, timezone string) ([]*models.Event, error)
	GetEvent(id string) (*models.Event, error)
	GetEventOwner(id string) (string, error)
	GetEventsCount() (int, error)
	CreateEvent(username string, title, description, timeVal, timezone string, duration time.Duration, notes []string) (*models.Event, error)
	UpdateEvent(id, title, description, timeVal, timezone string, duration time.Duration, notes []string) (*models.Event, error)
	DeleteEvent(id string) (bool, error)
	GetUsersCount() (int, error)
	UpdateUserTimezone(username, timezone string) (*models.User, error)
}

type Server struct {
	valid   Validator
	service Service
	proto.UnimplementedServiceServer
}

func NewServer(valid Validator, service Service) *Server {
	return &Server{valid: valid, service: service}
}

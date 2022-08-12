package http

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/gin-gonic/gin"
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

type Auth interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Validate(c *gin.Context)
}

type Server struct {
	service Service
	valid   Validator
	auth    Auth
}

func NewServer(service Service, valid Validator, auth Auth) *Server {
	return &Server{
		service: service,
		valid:   valid,
		auth:    auth,
	}
}

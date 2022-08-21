package http

import (
	"context"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
	ValidateGin(c *gin.Context)
	ValidateGrpc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
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

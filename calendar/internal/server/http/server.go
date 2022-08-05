package http

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/calendar"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Validator interface {
	Validate(interface{}) error
}

type Server struct {
	service *calendar.Service
	valid   Validator
}

func NewServer(service *calendar.Service, valid Validator) *Server {
	return &Server{
		service: service,
		valid:   valid,
	}
}

type GenericResponse struct {
	Message string
}

func badRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, GenericResponse{
		Message: err.Error(),
	})
}

func notFound(c *gin.Context, entity string) {
	c.JSON(http.StatusNotFound, GenericResponse{
		Message: fmt.Sprintf("%s not found", entity),
	})
}

func serverError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, GenericResponse{
		Message: err.Error(),
	})
}

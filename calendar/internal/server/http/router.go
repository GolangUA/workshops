package http

import (
	"errors"
	"github.com/Roma7-7-7/workshops/calendar/api"
	"github.com/gin-gonic/gin"
	"log"
)

// will hold http routes and will registrate them

func (s *Server) Register(e *gin.Engine) {
	e.Use(gin.CustomRecovery(recoveryHandler))

	e.POST("/login", s.auth.Login)
	e.GET("/logout", s.auth.Logout)

	api := e.Group("/api")
	api.Use(s.auth.Validate)
	s.registerEvents(api.Group("/events"))
	api.PUT("/user", s.UpdateUserTimezone)
}

func recoveryHandler(c *gin.Context, err interface{}) {
	log.Printf("unexpected panic: %v\n", err)
	api.ServerErrorA(c, errors.New("unexpected error occurred"))
}

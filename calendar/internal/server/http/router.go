package http

import "github.com/gin-gonic/gin"

// will hold http routes and will registrate them

func (s *Server) Register(e *gin.Engine) {
	api := e.Group("/api")
	s.registerEvents(api.Group("/events"))
}

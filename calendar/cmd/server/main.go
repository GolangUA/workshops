package main

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/config"
	"github.com/Roma7-7-7/workshops/calendar/internal/repository/postgre"
	"github.com/Roma7-7-7/workshops/calendar/internal/server/http"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/calendar"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/validator"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()
	repo := postgre.NewRepository(cfg.DSN())
	service := calendar.NewService(repo)
	server := http.NewServer(service, &validator.Service{})

	r := gin.Default()
	server.Register(r)
	r.Run(":5000")
}

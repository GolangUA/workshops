package http

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/api"
	"github.com/Roma7-7-7/workshops/calendar/internal/logging"
	"github.com/Roma7-7-7/workshops/calendar/internal/middleware/auth"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (s *Server) GetEvents(c *gin.Context) {
	req := validator.GetEvents{
		Title:    c.Query("title"),
		Timezone: c.Query("timezone"),
		DateFrom: c.Query("dateFrom"),
		TimeFrom: c.Query("timeFrom"),
		DateTo:   c.Query("dateTo"),
		TimeTo:   c.Query("timeTo"),
	}
	logging.Logger.Debug("get events", zap.Any("payload", req))
	if err := s.valid.Validate(&req); err != nil {
		api.BadRequestA(c, err)
		return
	}
	if req.Timezone == "" {
		req.Timezone = auth.GetContext(c).UserTimezone()
	}

	events, err := s.service.GetEvents(auth.GetContext(c).Username(), req.Title, req.DateFrom, req.TimeFrom, req.DateTo, req.TimeTo, req.Timezone)
	if err != nil {
		logging.Logger.Error("get events", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	result := make([]*api.Event, len(events))
	for i, e := range events {
		result[i] = eventToApi(e)
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) GetEvent(c *gin.Context) {
	id := c.Param("id")
	logging.Logger.Debug("get event", zap.String("id", id))
	if !s.validateOwner(c, id) {
		return
	}

	event, err := s.service.GetEvent(id)
	if err != nil {
		logging.Logger.Error("get event", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	c.JSON(http.StatusOK, eventToApi(event))
}

func (s *Server) PostEvent(c *gin.Context) {
	var req validator.CreateEvent
	if err := c.BindJSON(&req); err != nil {
		api.BadJSONA(c)
		return
	}
	logging.Logger.Debug("post event", zap.Any("payload", req))
	if err := s.valid.Validate(&req); err != nil {
		api.BadRequestA(c, err)
		return
	}

	e, err := s.service.CreateEvent(auth.GetContext(c).Username(), req.Title, req.Description, req.Time, req.Timezone, time.Duration(req.Duration)*time.Minute, req.Notes)
	if err != nil {
		logging.Logger.Error("create event", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	c.JSON(http.StatusCreated, eventToApi(e))
}

func (s *Server) PutEvent(c *gin.Context) {
	id := c.Param("id")
	var req validator.UpdateEvent
	if err := c.BindJSON(&req); err != nil {
		api.BadJSONA(c)
		return
	}
	req.ID = id
	logging.Logger.Debug("put event", zap.Any("payload", req))

	if err := s.valid.Validate(&req); err != nil {
		api.BadRequestA(c, err)
		return
	} else if !s.validateOwner(c, req.ID) {
		return
	}

	e, err := s.service.UpdateEvent(req.ID, req.Title, req.Description, req.Time, req.Timezone, time.Duration(req.Duration)*time.Minute, req.Notes)
	if err != nil {
		logging.Logger.Error("update event", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	c.JSON(http.StatusOK, eventToApi(e))
}

func (s *Server) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	logging.Logger.Debug("delete event", zap.String("id", id))
	if !s.validateOwner(c, id) {
		return
	}

	if _, err := s.service.DeleteEvent(id); err != nil {
		logging.Logger.Error("delete event", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	} else {
		c.AbortWithStatus(http.StatusOK)
		return
	}
}

func (s *Server) registerEvents(group *gin.RouterGroup) {
	group.GET("/", s.GetEvents)
	group.GET("/:id", s.GetEvent)
	group.POST("/", s.PostEvent)
	group.PUT("/:id", s.PutEvent)
	group.DELETE("/:id", s.DeleteEvent)
}

func eventToApi(e *models.Event) *api.Event {
	var tz string
	if l := e.TimeFrom.Location(); l == nil {
		tz = "UTC"
	} else {
		tz = l.String()
	}
	return &api.Event{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Time:        e.TimeFrom.Format("2006-01-02 15:04"),
		TimeZone:    tz,
		Duration:    int(e.TimeTo.Sub(e.TimeFrom).Minutes()),
		Notes:       e.Notes,
	}
}

func (s *Server) validateOwner(c *gin.Context, eventId string) bool {
	if owner, err := s.service.GetEventOwner(eventId); err != nil {
		logging.Logger.Error("get event owner of event", zap.String("ID", eventId), zap.Error(err))
		api.ServerErrorA(c, err)
		return false
	} else if owner == "" {
		api.NotFoundA(c, fmt.Sprintf("event with ID=\"%s\"", eventId))
		return false
	} else if owner != auth.GetContext(c).Username() {
		api.ForbiddenA(c, fmt.Sprintf("event with ID=\"%s\"", eventId))
		return false
	}
	return true
}

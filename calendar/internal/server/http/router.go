package http

import (
	"errors"
	"github.com/Roma7-7-7/workshops/calendar/api"
	"github.com/Roma7-7-7/workshops/calendar/internal/logging"
	appMetrics "github.com/Roma7-7-7/workshops/calendar/internal/middleware/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
)

// will hold http routes and will registrate them

var doNothingMetricsHandler = promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

func (s *Server) Register(app *gin.Engine, metrics *gin.Engine) {
	app.Use(gin.CustomRecovery(recoveryHandler))
	// Do nothing metrics handler just to count app requests in prometheus
	app.Use(gin.WrapH(doNothingMetricsHandler))

	app.POST("/login", s.auth.Login)
	app.GET("/logout", s.auth.Logout)

	apiG := app.Group("/api")
	apiG.Use(s.auth.ValidateGin)
	s.registerEvents(apiG.Group("/events"))
	apiG.PUT("/user", s.UpdateUserTimezone)

	metrics.GET("/metrics/prometheus", gin.WrapH(promhttp.Handler()))
	appMetrics.Init(s.service.GetUsersCount, s.service.GetEventsCount)
}

func recoveryHandler(c *gin.Context, err interface{}) {
	logging.Logger.Error("unexpected panic", zap.Any("panic", err))
	api.ServerErrorA(c, errors.New("unexpected error occurred"))
}

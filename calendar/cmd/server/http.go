package main

import (
	"context"
	"github.com/Roma7-7-7/workshops/calendar/internal/config"
	"github.com/Roma7-7-7/workshops/calendar/internal/logging"
	"github.com/Roma7-7-7/workshops/calendar/internal/middleware/auth"
	"github.com/Roma7-7-7/workshops/calendar/internal/repository/postgre"
	appHttp "github.com/Roma7-7-7/workshops/calendar/internal/server/http"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/calendar"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const shutdownTimout = 5 * time.Second

func main() {
	cfg := config.GetConfig()
	logging.Init(cfg.Env)
	repo := postgre.NewRepository(cfg.DSN())
	aut := auth.NewMiddleware(repo, cfg.JWT.Secret)
	service := calendar.NewService(repo)
	server := appHttp.NewServer(service, &validator.Service{}, aut)

	app := gin.Default()
	metrics := gin.New()

	server.Register(app, metrics)

	appServer := &http.Server{
		Addr:    ":5000",
		Handler: app,
	}
	metricsServer := &http.Server{
		Addr:    ":5005",
		Handler: metrics,
	}
	go func() {
		if err := appServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Fatal("listen app", zap.Error(err))
		}
	}()
	go func() {
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Fatal("listen metrics", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logging.Logger.Info("Shutdown Servers ...")
	var wg sync.WaitGroup

	appCtx, appCancel := context.WithTimeout(context.Background(), shutdownTimout)
	defer appCancel()
	appShutdownFailed := true
	go func() {
		wg.Add(1)
		if err := appServer.Shutdown(appCtx); err != nil {
			logging.Logger.Error("shutdown app", zap.Error(err))
		} else {
			appShutdownFailed = false
		}
		wg.Done()
	}()

	metricsCtx, metricsCancel := context.WithTimeout(context.Background(), shutdownTimout)
	defer metricsCancel()
	metricsShutdownFailed := true
	go func() {
		wg.Add(1)
		if err := metricsServer.Shutdown(metricsCtx); err != nil {
			logging.Logger.Error("shutdown metrics", zap.Error(err))
		} else {
			metricsShutdownFailed = false
		}
		wg.Done()
	}()

	shutdownComplete := make(chan struct{})
	go func() {
		wg.Wait()
		close(shutdownComplete)
	}()

	select {
	case <-time.After(shutdownTimout):
	case <-shutdownComplete:
	}

	if !appShutdownFailed && !metricsShutdownFailed {
		logging.Logger.Info("Servers Exited")
		return
	}

	if appShutdownFailed {
		logging.Logger.Info("app server shutdown failed")
	}
	if metricsShutdownFailed {
		logging.Logger.Info("metrics server shutdown failed")
	}
	os.Exit(1)
}

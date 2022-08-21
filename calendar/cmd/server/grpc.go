package main

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/config"
	"github.com/Roma7-7-7/workshops/calendar/internal/logging"
	"github.com/Roma7-7-7/workshops/calendar/internal/middleware/auth"
	"github.com/Roma7-7-7/workshops/calendar/internal/repository/postgre"
	pb "github.com/Roma7-7-7/workshops/calendar/internal/server/grpc"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/calendar"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/validator"
	"github.com/Roma7-7-7/workshops/calendar/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {
	cfg := config.GetConfig()
	logging.Init(cfg.Env)
	repo := postgre.NewRepository(cfg.DSN())
	aut := auth.NewMiddleware(repo, cfg.JWT.Secret)
	service := calendar.NewService(repo)
	server := pb.NewServer(&validator.Service{}, service)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logging.Logger.Fatal("start listener", zap.Error(err))
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(aut.ValidateGrpc))
	proto.RegisterServiceServer(s, server)
	if err := s.Serve(listener); err != nil {
		logging.Logger.Fatal("start server", zap.Error(err))
	}
}

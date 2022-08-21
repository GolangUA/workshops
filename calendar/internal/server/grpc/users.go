package grpc

import (
	"context"
	"github.com/Roma7-7-7/workshops/calendar/internal/logging"
	"github.com/Roma7-7-7/workshops/calendar/internal/middleware/auth"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/validator"
	"github.com/Roma7-7-7/workshops/calendar/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUserTimezone(c context.Context, r *proto.UpdateUserTimezoneRequest) (*proto.UpdateUserTimezoneResponse, error) {
	req := validator.UserTimezone{
		Username: r.GetUserTimezone().GetUsername(),
		Timezone: r.GetUserTimezone().GetTimezone(),
	}
	logging.Logger.Debug("update user timezone", zap.Any("payload", req))

	if err := s.valid.Validate(&req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := auth.GetContextGrpc(c).Username()
	if user != req.Username {
		return nil, status.Error(codes.PermissionDenied, "You are not allowed to update this user")
	}

	u, err := s.service.UpdateUserTimezone(req.Username, req.Timezone)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.UpdateUserTimezoneResponse{UserTimezone: userToProto(u)}, nil
}

func userToProto(u *models.User) *proto.UserTimezone {
	return &proto.UserTimezone{
		Username: u.Name,
		Timezone: u.Timezone,
	}
}

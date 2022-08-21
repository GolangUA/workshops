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
	"time"
)

func (s *Server) GetEvents(c context.Context, r *proto.GetEventsRequest) (*proto.GetEventsResponse, error) {
	req := validator.GetEvents{
		Title:    r.GetTitle(),
		Timezone: r.GetTimezone(),
		DateFrom: r.GetDateFrom(),
		TimeFrom: r.GetTimeFrom(),
		DateTo:   r.GetDateTo(),
		TimeTo:   r.GetTimeTo(),
	}
	logging.Logger.Debug("get events", zap.Any("payload", req))
	if err := s.valid.Validate(&req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ctx := auth.GetContextGrpc(c)
	if req.Timezone == "" {
		req.Timezone = ctx.UserTimezone()
	}

	events, err := s.service.GetEvents(ctx.Username(), req.Title, req.DateFrom, req.TimeFrom, req.DateTo, req.TimeTo, req.Timezone)
	if err != nil {
		logging.Logger.Error("get events", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*proto.Event, len(events))
	for i, e := range events {
		result[i] = toProto(e)
	}
	return &proto.GetEventsResponse{Events: result}, nil
}

func (s *Server) GetEvent(c context.Context, r *proto.GetEventRequest) (*proto.GetEventResponse, error) {
	id := r.Id
	logging.Logger.Debug("get event", zap.String("id", id))
	if !s.validateOwner(c, id) {
		return nil, status.Error(codes.PermissionDenied, "not owner of event")
	}

	event, err := s.service.GetEvent(id)
	if err != nil {
		logging.Logger.Error("get event", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.GetEventResponse{Event: toProto(event)}, nil
}

func (s *Server) CreateEvent(c context.Context, r *proto.CreateEventRequest) (*proto.CreateEventResponse, error) {
	req := validator.CreateEvent{
		Title:       r.GetEvent().GetTitle(),
		Description: r.GetEvent().GetDescription(),
		Time:        r.GetEvent().GetTime(),
		Timezone:    r.GetEvent().GetTimeZone(),
		Duration:    int(r.GetEvent().GetDuration()),
		Notes:       r.GetEvent().GetNotes(),
	}

	logging.Logger.Debug("post event", zap.Any("payload", req))
	if err := s.valid.Validate(&req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	e, err := s.service.CreateEvent(auth.GetContextGrpc(c).Username(), req.Title, req.Description, req.Time, req.Timezone, time.Duration(req.Duration)*time.Minute, req.Notes)
	if err != nil {
		logging.Logger.Error("create event", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.CreateEventResponse{Event: toProto(e)}, nil
}

func (s *Server) UpdateEvent(c context.Context, r *proto.UpdateEventRequest) (*proto.UpdateEventResponse, error) {
	req := validator.UpdateEvent{
		ID:          r.GetEvent().GetId(),
		Title:       r.GetEvent().GetTitle(),
		Description: r.GetEvent().GetDescription(),
		Time:        r.GetEvent().GetTime(),
		Timezone:    r.GetEvent().GetTimeZone(),
		Duration:    int(r.GetEvent().GetDuration()),
		Notes:       r.GetEvent().GetNotes(),
	}
	logging.Logger.Debug("put event", zap.Any("payload", req))

	if err := s.valid.Validate(&req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if !s.validateOwner(c, req.ID) {
		return nil, status.Error(codes.PermissionDenied, "not owner of event")
	}

	e, err := s.service.UpdateEvent(req.ID, req.Title, req.Description, req.Time, req.Timezone, time.Duration(req.Duration)*time.Minute, req.Notes)
	if err != nil {
		logging.Logger.Error("update event", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.UpdateEventResponse{Event: toProto(e)}, nil
}

func (s *Server) DeleteEvent(c context.Context, r *proto.DeleteEventRequest) (*proto.EmptyResponse, error) {
	id := r.GetId()
	logging.Logger.Debug("delete event", zap.String("id", id))
	if !s.validateOwner(c, id) {
		return nil, status.Error(codes.PermissionDenied, "not owner of event")
	}

	if _, err := s.service.DeleteEvent(id); err != nil {
		logging.Logger.Error("delete event", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		return &proto.EmptyResponse{}, nil
	}
}

func (s *Server) validateOwner(c context.Context, eventId string) bool {
	if owner, err := s.service.GetEventOwner(eventId); err != nil {
		logging.Logger.Error("get event owner of event", zap.String("ID", eventId), zap.Error(err))
		return false
	} else if owner == "" {
		return false
	} else if owner != auth.GetContextGrpc(c).Username() {
		return false
	}
	return true
}

func toProto(e *models.Event) *proto.Event {
	var tz string
	if l := e.TimeFrom.Location(); l == nil {
		tz = "UTC"
	} else {
		tz = l.String()
	}
	return &proto.Event{
		Id:          &e.ID,
		Title:       e.Title,
		Description: &e.Description,
		Time:        e.TimeFrom.Format("2006-01-02 15:04"),
		TimeZone:    tz,
		Duration:    uint32(e.TimeTo.Sub(e.TimeFrom).Minutes()),
		Notes:       e.Notes,
	}
}

package interceptors

import (
	"context"
	"encoding/base64"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthMD struct{}

const (
	authHeader = "authorization"
	bearerAuth = "bearer"
)

func (a *AuthMD) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		credentials := a.getAuthCredentials(ctx)
		if credentials == "" {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		decoded, err := base64.StdEncoding.DecodeString(credentials)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "can't parse credentials")
		}

		data := strings.Split(string(decoded), ":")
		if len(data) != 2 {
			return nil, status.Error(codes.Unauthenticated, "can't parse credentials")
		}

		ctx = context.WithValue(ctx, "user", data[0])
		return handler(ctx, req)
	}
}

func (a *AuthMD) getAuthCredentials(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(authHeader)
	if len(values) == 0 {
		return ""
	}

	fields := strings.SplitN(values[0], " ", 2)
	if len(fields) < 2 {
		return ""
	}

	if !strings.EqualFold(fields[0], bearerAuth) {
		return ""
	}
	return fields[1]
}

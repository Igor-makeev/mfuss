package interceptors

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"mfuss/internal/grpc/auth"
)

// AuthUnaryInterceptor отвечает за аутентификацию grpc-клиентов.
func (i interceptors) AuthUnaryInterceptor(ctx context.Context,
	req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	var user uuid.UUID
	var signed string

	authHeader, ok := md["user"]
	if !ok {
		user, signed = i.a.Gen()
	} else {
		token := authHeader[0]
		user, err = i.a.Load(token)
		if errors.Is(err, auth.ErrUnauthorized) {
			user, signed = i.a.Gen()
		} else if err != nil {
			return nil, status.Errorf(codes.Internal, "Server error")
		}
	}

	if len(signed) > 0 {
		err = grpc.SendHeader(ctx, metadata.Pairs("user", signed))
		if err != nil {
			log.Printf("unable to send metadata: %v", err)
		}
	}

	ctx = context.WithValue(ctx, "userID", user)
	return handler(ctx, req)
}

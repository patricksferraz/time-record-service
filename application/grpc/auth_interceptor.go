package grpc

import (
	"context"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/entity"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/service"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	AuthService *service.AuthService
	Claims      *entity.Claims
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		span, ctx := apm.StartSpan(ctx, "Unary", "gRPC interceptor")
		defer span.End()

		log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
		log.WithField("method", info.FullMethod).Info("unary interceptor")

		err := a.authorize(ctx, info.FullMethod)
		if err != nil {
			log.WithError(err)
			apm.CaptureError(ctx, err).Send()
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		span, ctx := apm.StartSpan(ss.Context(), "Stream", "gRPC interceptor")
		defer span.End()

		log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
		log.WithField("method", info.FullMethod).Info("stream interceptor")

		err := a.authorize(ctx, info.FullMethod)
		if err != nil {
			log.WithError(err)
			apm.CaptureError(ctx, err).Send()
			return err
		}

		return handler(srv, ss)
	}
}

func (a *AuthInterceptor) authorize(ctx context.Context, method string) error {
	span, ctx := apm.StartSpan(ctx, "authorize", "gRPC interceptor")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("method", method).Info("authorize params")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := status.Error(codes.Unauthenticated, "metadata is not provided")
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}
	log.WithField("metadata", md).Info("handling metadata request")

	values := md["authorization"]
	if len(values) == 0 {
		err := status.Error(codes.Unauthenticated, "authorization token is not provided")
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}

	accessToken := values[0]
	claims, err := a.AuthService.Verify(ctx, accessToken)
	if err != nil {
		err := status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}
	log.WithField("claims", claims).Info("verify accessToken")

	a.Claims = claims

	for _, role := range claims.Roles {
		if role == method {
			return nil
		}
	}

	err = status.Error(codes.PermissionDenied, "no permission to access this RPC")
	log.WithError(err)
	apm.CaptureError(ctx, err).Send()
	return err
}

func NewAuthInterceptor(authService *service.AuthService) *AuthInterceptor {
	return &AuthInterceptor{
		AuthService: authService,
	}
}

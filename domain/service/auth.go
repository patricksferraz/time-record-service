package service

import (
	"context"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/application/grpc/pb"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/model"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type AuthService struct {
	Service pb.AuthServiceClient
}

func (a *AuthService) Verify(ctx context.Context, accessToken string) (*model.Claims, error) {
	span, ctx := apm.StartSpan(ctx, "Verify", "auth domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	req := &pb.FindClaimsByTokenRequest{
		AccessToken: accessToken,
	}
	log.WithField("req", req).Info("FindClaimsByToken request")

	_claims, err := a.Service.FindClaimsByToken(ctx, req)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("claims", _claims).Info("claims response")

	claims, err := model.NewClaims(_claims.EmployeeId, _claims.Roles)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("claims", claims).Info("claims created")

	return claims, nil
}

func NewAuthService(service pb.AuthServiceClient) *AuthService {
	return &AuthService{
		Service: service,
	}
}

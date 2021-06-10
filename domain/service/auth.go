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

func (a *AuthService) Verify(ctx context.Context, accessToken string) (*model.EmployeeClaims, error) {
	span, ctx := apm.StartSpan(ctx, "Verify", "auth domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	req := &pb.FindEmployeeClaimsByTokenRequest{
		AccessToken: accessToken,
	}
	log.WithField("req", req).Info("FindEmployeeClaimsByToken request")

	employee, err := a.Service.FindEmployeeClaimsByToken(ctx, req)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("employee", employee).Info("employee response")

	employeeClaims, err := model.NewEmployeeClaims(employee.Id, employee.Roles)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("employeeClaims", employeeClaims).Info("employee created")

	return employeeClaims, nil
}

func NewAuthService(service pb.AuthServiceClient) *AuthService {
	return &AuthService{
		Service: service,
	}
}

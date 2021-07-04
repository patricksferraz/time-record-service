package service

import (
	"context"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/application/grpc/pb"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/entity"
	"google.golang.org/grpc"
)

type AuthService struct {
	service pb.AuthServiceClient
}

func NewAuthService(cc *grpc.ClientConn) *AuthService {
	return &AuthService{
		service: pb.NewAuthServiceClient(cc),
	}
}

func (a *AuthService) Verify(ctx context.Context, accessToken string) (*entity.Claims, error) {
	req := &pb.FindClaimsByTokenRequest{
		AccessToken: accessToken,
	}

	_claims, err := a.service.FindClaimsByToken(ctx, req)
	if err != nil {
		return nil, err
	}

	claims, err := entity.NewClaims(_claims.EmployeeId, _claims.Roles)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

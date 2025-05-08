package external

import (
	"context"

	"github.com/patricksferraz/time-record-service/domain/entity"
	"github.com/patricksferraz/time-record-service/proto/pb"
	"google.golang.org/grpc"
)

type AuthClient struct {
	c pb.AuthKeycloakAclClient
}

func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	return &AuthClient{
		c: pb.NewAuthKeycloakAclClient(cc),
	}
}

func (a *AuthClient) Verify(ctx context.Context, accessToken string) (*entity.Claims, error) {
	req := &pb.FindClaimsByTokenRequest{
		AccessToken: accessToken,
	}

	_claims, err := a.c.FindClaimsByToken(ctx, req)
	if err != nil {
		return nil, err
	}

	claims, err := entity.NewClaims(_claims.EmployeeId, _claims.Roles)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/c-4u/time-record-service/domain/entity"
	"github.com/c-4u/time-record-service/infrastructure/external"
	"github.com/c-4u/time-record-service/logger"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type AuthMiddleware struct {
	AuthClient  *external.AuthClient
	Claims      *entity.Claims
	AccessToken *string
}

func (a *AuthMiddleware) Require() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

		accessToken := ctx.Request.Header.Get("Authorization")
		if accessToken == "" {
			err := errors.New("authorization token is not provided")
			log.WithError(err)
			apm.CaptureError(ctx, err).Send()
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		claims, err := a.AuthClient.Verify(ctx, accessToken)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			ctx.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("access token is invalid: %v", err)})
			ctx.Abort()
			return
		}
		log.WithField("claims", claims).Info("verify accessToken")

		a.Claims = claims
		a.AccessToken = &accessToken

		// TODO: adds retricted permissions
		// for _, role := range claims.Roles {
		// 	if role == method {
		// 		return nil
		// 	}
		// }

		// return status.Error(codes.PermissionDenied, "no permission to access this RPC")
	}
}

func NewAuthMiddleware(authClient *external.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		AuthClient: authClient,
	}
}

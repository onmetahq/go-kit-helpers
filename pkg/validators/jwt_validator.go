package validators

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt"
	metahttp "github.com/krishnateja262/meta-http/pkg/meta_http"
	"github.com/onmetahq/go-kit-helpers/pkg/logger"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
)

func JWTValidator(hmacSecret string, logger logger.CtxLogger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			tokenString, ok := ctx.Value(models.JWTContextKey).(string)
			if !ok {
				logger.Context(ctx).Error().Log("msg", "Invalid JWT", "token", tokenString)
				return nil, models.ErrUnauthorized
			}

			claims := &models.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					logger.Context(ctx).Error().Log("msg", "Invalid JWT header method", "token", tokenString, "error", token.Method.Alg())
					return nil, models.ErrUnexpectedSigningMethod
				}

				return []byte(hmacSecret), nil
			})

			if err != nil {
				if e, ok := err.(*jwt.ValidationError); ok {
					switch {
					case e.Errors&jwt.ValidationErrorMalformed != 0:
						logger.Context(ctx).Error().Log("msg", "Malformed JWT", "token", tokenString)
					case e.Errors&jwt.ValidationErrorExpired != 0:
						logger.Context(ctx).Error().Log("msg", "Expired JWT", "token", tokenString)
					case e.Errors&jwt.ValidationErrorNotValidYet != 0:
						logger.Context(ctx).Error().Log("msg", "Inactive JWT", "token", tokenString)
					case e.Inner != nil:
						logger.Context(ctx).Error().Log("msg", "Inner JWT", "token", tokenString)
					}
				}
				logger.Context(ctx).Error().Log("msg", "Error JWT", "token", tokenString, "error", err)
				return nil, models.ErrUnauthorized
			}

			if !token.Valid {
				logger.Context(ctx).Error().Log("msg", "Invalid Token", "token", tokenString)
				return nil, models.ErrUnauthorized
			}

			ctx = context.WithValue(ctx, models.JWTClaimsContextKey, claims)
			ctx = context.WithValue(ctx, metahttp.TenantID, claims.TenantID)
			ctx = context.WithValue(ctx, models.USERID, claims.UserId) // TODO: Remove it soon
			ctx = context.WithValue(ctx, metahttp.UserID, claims.UserId)

			return next(ctx, request)
		}
	}
}

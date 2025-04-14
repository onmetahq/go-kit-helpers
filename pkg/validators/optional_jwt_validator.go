package validators

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/golang-jwt/jwt"
	ctxLogger "github.com/onmetahq/go-kit-helpers/pkg/logger"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

func OptionalJWTValidator(hmacSecret string, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			lg := ctxLogger.NewCtxLogger(logger)

			tokenString, ok := ctx.Value(models.JWTContextKey).(string)
			if !ok {
				return next(ctx, request)
			}

			claims := &models.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					lg.Context(ctx).Error().Log("msg", "Invalid JWT header method", "token", tokenString, "error", token.Method.Alg())
					return nil, models.ErrUnexpectedSigningMethod
				}

				return []byte(hmacSecret), nil
			})

			if err != nil {
				if e, ok := err.(*jwt.ValidationError); ok {
					switch {
					case e.Errors&jwt.ValidationErrorMalformed != 0:
						lg.Context(ctx).Error().Log("msg", "Malformed JWT", "token", tokenString)
					case e.Errors&jwt.ValidationErrorExpired != 0:
						lg.Context(ctx).Error().Log("msg", "Expired JWT", "token", tokenString)
					case e.Errors&jwt.ValidationErrorNotValidYet != 0:
						lg.Context(ctx).Error().Log("msg", "Inactive JWT", "token", tokenString)
					case e.Inner != nil:
						lg.Context(ctx).Error().Log("msg", "Inner JWT", "token", tokenString)
					}
				}
				lg.Context(ctx).Error().Log("msg", "Error JWT", "token", tokenString, "error", err)
				return nil, models.ErrUnauthorized
			}

			if !token.Valid {
				lg.Context(ctx).Error().Log("msg", "Invalid Token", "token", tokenString)
				return nil, models.ErrUnauthorized
			}

			ctx = context.WithValue(ctx, models.JWTClaimsContextKey, claims)
			ctx = context.WithValue(ctx, ctxKeys.TenantID, claims.TenantID)
			ctx = context.WithValue(ctx, models.USERID, claims.UserId) // TODO: Remove it soon
			ctx = context.WithValue(ctx, ctxKeys.UserID, claims.UserId)

			return next(ctx, request)
		}
	}
}

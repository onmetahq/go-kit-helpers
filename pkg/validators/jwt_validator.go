package validators

import (
	"context"
	"log/slog"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

func JWTValidator(hmacSecret string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			tokenString, ok := ctx.Value(models.JWTContextKey).(string)
			if !ok {
				slog.ErrorContext(ctx, "JWT not found in context")
				return nil, models.ErrUnauthorized
			}

			claims := &models.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					slog.ErrorContext(ctx, "Invalid JWT header method", "error", token.Method.Alg())
					return nil, models.ErrUnexpectedSigningMethod
				}

				return []byte(hmacSecret), nil
			})

			if err != nil {
				if e, ok := err.(*jwt.ValidationError); ok {
					switch {
					case e.Errors&jwt.ValidationErrorMalformed != 0:
						slog.ErrorContext(ctx, "Malformed JWT")
					case e.Errors&jwt.ValidationErrorExpired != 0:
						slog.ErrorContext(ctx, "Expired JWT")
					case e.Errors&jwt.ValidationErrorNotValidYet != 0:
						slog.ErrorContext(ctx, "Inactive JWT")
					case e.Inner != nil:
						slog.ErrorContext(ctx, "Inner JWT error", "error", e.Inner)
					}
				}
				slog.ErrorContext(ctx, "JWT validation error", "error", err)
				return nil, models.ErrUnauthorized
			}

			if !token.Valid {
				slog.ErrorContext(ctx, "Invalid JWT")
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

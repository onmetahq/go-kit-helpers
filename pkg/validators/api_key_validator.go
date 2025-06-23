package validators

import (
	"context"
	"log/slog"

	"github.com/go-kit/kit/endpoint"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

func APIKeyValidator(key string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			apiKey, ok := ctx.Value(ctxKeys.APIContextKey).(string)
			if !ok {
				slog.ErrorContext(ctx, "API Key not found in context", "apikey", apiKey)
				return nil, models.ErrUnauthorized
			}

			if apiKey != key {
				slog.ErrorContext(ctx, "API key mismatch", "received", apiKey)
				return nil, models.ErrUnauthorized
			}

			return next(ctx, request)
		}
	}
}

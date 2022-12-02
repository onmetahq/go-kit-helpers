package validators

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/onmetahq/go-kit-helpers/pkg/logger"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	onmetamodels "github.com/onmetahq/meta-http/pkg/models"
)

func APIKeyValidator(key string, logger logger.CtxLogger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			apiKey, ok := ctx.Value(onmetamodels.APIContextKey).(string)
			if !ok {
				logger.Context(ctx).Error().Log("msg", "Invalid API Key", "apikey", apiKey)
				return nil, models.ErrUnauthorized
			}

			if apiKey != key {
				logger.Context(ctx).Error().Log("msg", "API key mismatch", "apikey", apiKey)
				return nil, models.ErrUnauthorized
			}

			return next(ctx, request)
		}
	}
}
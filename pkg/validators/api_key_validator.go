package validators

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	ctxLogger "github.com/onmetahq/go-kit-helpers/pkg/logger"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

func APIKeyValidator(key string, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			lg := ctxLogger.NewCtxLogger(logger)

			apiKey, ok := ctx.Value(ctxKeys.APIContextKey).(string)
			if !ok {
				lg.Context(ctx).Error().Log("msg", "Invalid API Key", "apikey", apiKey)
				return nil, models.ErrUnauthorized
			}

			if apiKey != key {
				lg.Context(ctx).Error().Log("msg", "API key mismatch", "apikey", apiKey)
				return nil, models.ErrUnauthorized
			}

			return next(ctx, request)
		}
	}
}

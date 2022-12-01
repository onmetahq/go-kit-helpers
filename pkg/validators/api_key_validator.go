package validators

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	metahttp "github.com/krishnateja262/meta-http/pkg/meta_http"
	"github.com/onmetahq/go-kit-helpers/pkg/logger"
)

func APIKeyValidator(key string, logger logger.CtxLogger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			apiKey, ok := ctx.Value(metahttp.APIContextKey).(string)
			if !ok {
				logger.Context(ctx).Error().Log("msg", "Invalid API Key", "apikey", apiKey)
				return nil, ErrUnauthorized
			}

			if apiKey != key {
				logger.Context(ctx).Error().Log("msg", "API key mismatch", "apikey", apiKey)
				return nil, ErrUnauthorized
			}

			return next(ctx, request)
		}
	}
}

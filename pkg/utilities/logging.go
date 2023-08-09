package utilities

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	ctxLogger "github.com/onmetahq/go-kit-helpers/pkg/logger"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
)

func Logger(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			lg := ctxLogger.NewCtxLogger(logger)

			startTime := time.Now()
			lg.Context(ctx).Info().Log("msg", "Server request started", "request", request, "path", ctx.Value(models.URLPath), "method", ctx.Value(models.HttpMethod))
			res, err := next(ctx, request)
			elapsedTime := time.Since(startTime)
			lg.Context(ctx).Info().Log("msg", "Server request ended", "response", res, "duration", elapsedTime.Milliseconds(), "path", ctx.Value(models.URLPath), "method", ctx.Value(models.HttpMethod))

			return res, err
		}
	}
}

package utilities

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
)

func Logger() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			startTime := time.Now()
			slog.InfoContext(ctx, "Server request started", "request", request, "path", ctx.Value(models.URLPath), "method", ctx.Value(models.HttpMethod))
			res, err := next(ctx, request)
			elapsedTime := time.Since(startTime)
			slog.InfoContext(ctx, "Server request ended", "response", res, "duration", elapsedTime.Milliseconds(), "path", ctx.Value(models.URLPath), "method", ctx.Value(models.HttpMethod))

			return res, err
		}
	}
}

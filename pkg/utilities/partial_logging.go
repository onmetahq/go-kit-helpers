package utilities

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
)

func PartialLogger() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			startTime := time.Now()
			res, err := next(ctx, request)
			elapsedTime := time.Since(startTime)
			if err == nil {
				slog.InfoContext(ctx, "Server request ended", "duration", elapsedTime.Milliseconds(), "path", ctx.Value(models.URLPath), "method", ctx.Value(models.HttpMethod))
			} else {
				slog.ErrorContext(ctx, "Server request ended with error", "request", request, "response", res, "error", err, "duration", elapsedTime.Milliseconds(), "path", ctx.Value(models.URLPath), "method", ctx.Value(models.HttpMethod))
			}
			return res, err
		}
	}
}

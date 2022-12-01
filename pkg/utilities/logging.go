package utilities

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/onmetahq/go-kit-helpers/pkg/logger"
)

func Logger(logger logger.CtxLogger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			startTime := time.Now()
			logger.Context(ctx).Debug().Log("msg", "Server request started", "request", request)
			res, err := next(ctx, request)
			elapsedTime := time.Since(startTime)
			logger.Context(ctx).Debug().Log("msg", "Server request ended", "response", res, "duration", elapsedTime.Milliseconds())
			return res, err
		}
	}
}

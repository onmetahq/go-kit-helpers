package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type Instrument struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
}

func NewInstrument(nameSpace string, subSystem string) Instrument {
	fieldKeys := []string{"path", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: nameSpace,
		Subsystem: subSystem,
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: nameSpace,
		Subsystem: subSystem,
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	return Instrument{
		RequestCount:   requestCount,
		RequestLatency: requestLatency,
	}
}

func MetricMiddleware(ins Instrument) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			path, ok := ctx.Value(models.URLPathTemplate).(string)
			if !ok {
				urlPath, ok1 := ctx.Value(models.URLPath).(string)
				if ok1 {
					path = urlPath
				} else {
					path = "/unable-to-find-path"
				}
			}
			defer func(begin time.Time) {
				lvs := []string{"path", path, "error", fmt.Sprint(err != nil)}
				ins.RequestCount.With(lvs...).Add(1)
				ins.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
			}(time.Now())

			res, err := next(ctx, request)
			return res, err
		}
	}
}

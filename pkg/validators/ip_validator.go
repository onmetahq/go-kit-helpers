package validators

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

func IPValidator() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			apikey, ok := ctx.Value(ctxKeys.MerchantAPIKey).(string)
			if !ok {
				slog.ErrorContext(ctx, "Invalid Merchant API key", "apikey", apikey)
				return nil, models.ErrUnauthorized
			}

			ip, ok1 := ctx.Value(ctxKeys.XForwardedFor).(string)
			if !ok1 {
				slog.DebugContext(ctx, "X-Forwarded-For header not found", "apikey", apikey)
				return next(ctx, request)
			}

			incomingIps := strings.Split(ip, ",")
			config := os.Getenv(apikey)
			if len(config) == 0 {
				slog.DebugContext(ctx, "Whitelisting not configured for API key", "apikey", apikey)
				return next(ctx, request)
			}

			eligibleIps := strings.Split(config, ",")
			eligibleIpMap := map[string]bool{}
			for _, eligibleIp := range eligibleIps {
				eligibleIpMap[strings.TrimSpace(eligibleIp)] = true
			}

			for _, incomingIp := range incomingIps {
				if eligibleIpMap[strings.TrimSpace(incomingIp)] {
					return next(ctx, request)
				}
			}

			slog.ErrorContext(ctx, "Unauthorized IP address", "incomingIP", ip, "validIps", config)
			return nil, models.ErrUnauthorized
		}
	}
}

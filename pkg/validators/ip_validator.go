package validators

import (
	"context"
	"os"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	ctxLogger "github.com/onmetahq/go-kit-helpers/pkg/logger"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

func IPValidator(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			lg := ctxLogger.NewCtxLogger(logger)

			apikey, ok := ctx.Value(ctxKeys.MerchantAPIKey).(string)
			if !ok {
				lg.Context(ctx).Error().Log("msg", "Invalid Merchant API key")
				return nil, models.ErrUnauthorized
			}

			ip, ok1 := ctx.Value(ctxKeys.XForwardedFor).(string)
			if !ok1 {
				lg.Context(ctx).Debug().Log("msg", "did not find client ip address, overriding ip check")
				return next(ctx, request)
			}
			incomingIps := strings.Split(ip, ",")

			config := os.Getenv(apikey)
			if len(config) == 0 {
				lg.Context(ctx).Debug().Log("msg", "whitelisting is not configured for the api key", "apiKey", apikey)
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

			lg.Context(ctx).Error().Log("msg", "did not find the IP address in whitelisted addresses", "incomingIP", ip, "validIps", config)
			return nil, models.ErrUnauthorized
		}
	}
}

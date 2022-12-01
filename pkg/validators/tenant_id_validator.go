package validators

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	metahttp "github.com/krishnateja262/meta-http/pkg/meta_http"
	"github.com/onmetahq/go-kit-helpers/pkg/logger"
)

func TenantIdValidator(pathParam string, logger logger.CtxLogger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			pathParams := ctx.Value(PathParamsContextKey).(map[string]string)
			id := pathParams[pathParam]

			tenantId := ctx.Value(metahttp.TenantID).(string)

			if tenantId != id {
				logger.Context(ctx).Error().Log("level", "error", "msg", "Tenant ID does not match with claims ID", "tenantId", id, "claimsId", tenantId)
				return nil, ErrUnauthorized
			}

			return next(ctx, request)
		}
	}
}

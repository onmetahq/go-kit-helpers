package validators

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/onmetahq/go-kit-helpers/pkg/logger"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

func TenantIdValidator(pathParam string, logger logger.CtxLogger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			pathParams := ctx.Value(models.PathParamsContextKey).(map[string]string)
			id := pathParams[pathParam]

			tenantId := ctx.Value(ctxKeys.TenantID).(string)

			if tenantId != id {
				logger.Context(ctx).Error().Log("level", "error", "msg", "Tenant ID does not match with claims ID", "tenantId", id, "claimsId", tenantId)
				return nil, models.ErrUnauthorized
			}

			return next(ctx, request)
		}
	}
}

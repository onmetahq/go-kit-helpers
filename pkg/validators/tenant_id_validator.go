package validators

import (
	"context"
	"log/slog"

	"github.com/go-kit/kit/endpoint"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

func TenantIdValidator(pathParam string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			pathParams := ctx.Value(models.PathParamsContextKey).(map[string]string)
			id := pathParams[pathParam]

			tenantId := ctx.Value(ctxKeys.TenantID).(string)

			if tenantId != id {
				slog.ErrorContext(ctx, "Tenant ID does not match with claims ID", "tenantId", id, "claimsId", tenantId)
				return nil, models.ErrUnauthorized
			}

			return next(ctx, request)
		}
	}
}

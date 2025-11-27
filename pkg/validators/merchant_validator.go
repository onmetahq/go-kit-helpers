package validators

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	metahttp "github.com/onmetahq/meta-http/pkg/meta_http"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

type GetKybStatusRequest struct {
	MerchantId string `json:"merchantId"`
}

type ValidateKybResponse struct {
	Success            bool   `json:"success"`
	IsApproved         bool   `json:"is_approved"`
	Message            string `json:"msg"`
	IsUnderMaintenance bool   `json:"isUnderMaintenance"`
}

func MerchantValidatorMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			tenantId, ok := ctx.Value(ctxKeys.TenantID).(string)
			if !ok || tenantId == "" {
				slog.ErrorContext(ctx, "Tenant ID not found in context")
				return nil, models.ErrUnauthorized
			}

			baseURL := os.Getenv("ENTITIES_URL")
			httpClient := metahttp.NewClient(baseURL, slog.Default(), 10*time.Second)

			headers := map[string]string{
				"Content-Type": "application/json",
				"apikey":       os.Getenv("REQUEST_KEY"),
			}

			requestBody := GetKybStatusRequest{
				MerchantId: tenantId,
			}

			var apiResponse ValidateKybResponse
			_, err = httpClient.Post(ctx, "/merchant/v1/kyb/validate", headers, requestBody, &apiResponse)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to fetch KYB status", "tenantId", tenantId, "err", err)
				return nil, models.ErrInternalServerError
			}

			if !apiResponse.Success {
				slog.ErrorContext(ctx, "Validate API returned unsuccessful response", "tenantId", tenantId, "message", apiResponse.Message)
				return nil, models.ErrInternalServerError
			}

			if !apiResponse.IsApproved || apiResponse.IsUnderMaintenance {
				slog.WarnContext(ctx, "Merchant KYB not approved", "tenantId", tenantId, "message", apiResponse.Message)
				return nil, models.ErrUnauthorized
			}

			return next(ctx, request)
		}
	}
}

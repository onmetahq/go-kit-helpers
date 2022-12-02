package options

import (
	"context"
	"net/http"
	"strings"

	metahttp "github.com/krishnateja262/meta-http/pkg/meta_http"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
)

func PutHeadersInCtx(ctx context.Context, r *http.Request) context.Context {
	ctx = metahttp.FetchContextFromHeaders(ctx, r)

	if r.Header.Get("Authorization") != "" {
		authHeaderStrings := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeaderStrings) > 1 {
			ctx = context.WithValue(ctx, models.JWTContextKey, authHeaderStrings[1])
		}
	}

	if r.Header.Get("apikey") != "" {
		ctx = context.WithValue(ctx, metahttp.APIContextKey, r.Header.Get("apikey"))
	}

	if r.Header.Get("x-api-key") != "" {
		ctx = context.WithValue(ctx, metahttp.MerchantAPIKey, r.Header.Get("x-api-key"))
	}

	return ctx
}

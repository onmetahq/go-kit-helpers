package options

import (
	"context"
	"net/http"
	"strings"

	"github.com/onmetahq/go-kit-helpers/pkg/models"
	onmetamodels "github.com/onmetahq/meta-http/pkg/models"
	"github.com/onmetahq/meta-http/pkg/utils"
)

func PutHeadersInCtx(ctx context.Context, r *http.Request) context.Context {
	ctx = utils.FetchContextFromHeaders(ctx, r)

	if r.Header.Get("Authorization") != "" {
		authHeaderStrings := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeaderStrings) > 1 {
			ctx = context.WithValue(ctx, models.JWTContextKey, authHeaderStrings[1])
		}
	}

	if r.Header.Get("apikey") != "" {
		ctx = context.WithValue(ctx, onmetamodels.APIContextKey, r.Header.Get("apikey"))
	}

	if r.Header.Get("x-api-key") != "" {
		ctx = context.WithValue(ctx, onmetamodels.MerchantAPIKey, r.Header.Get("x-api-key"))
	}

	return ctx
}

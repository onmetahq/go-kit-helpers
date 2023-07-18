package options

import (
	"context"
	"net/http"

	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

func PutIPInCtx(ctx context.Context, r *http.Request) context.Context {
	val := r.Header.Get(string(ctxKeys.XForwardedFor))
	if val != "" {
		ctx = context.WithValue(ctx, ctxKeys.XForwardedFor, val)
	}
	return ctx
}

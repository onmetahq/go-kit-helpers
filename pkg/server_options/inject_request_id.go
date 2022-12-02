package options

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	metahttp "github.com/krishnateja262/meta-http/pkg/meta_http"
)

func PutRequestIdInCtx(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, metahttp.RequestID, uuid.NewString())
}

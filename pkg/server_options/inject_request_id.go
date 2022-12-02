package options

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/onmetahq/meta-http/pkg/models"
)

func PutRequestIdInCtx(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, models.RequestID, uuid.NewString())
}

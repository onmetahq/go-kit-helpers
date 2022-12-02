package options

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
)

func PutPathParamsInCtx(ctx context.Context, r *http.Request) context.Context {
	vars := mux.Vars(r)
	return context.WithValue(ctx, models.PathParamsContextKey, vars)
}

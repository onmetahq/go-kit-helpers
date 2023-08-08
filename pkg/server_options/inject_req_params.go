package options

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	"github.com/onmetahq/meta-http/pkg/utils"
)

func PutReqInCtx(ctx context.Context, r *http.Request) context.Context {
	vars := mux.Vars(r)
	ctx = context.WithValue(ctx, models.PathParamsContextKey, vars)
	ctx = context.WithValue(ctx, models.HttpMethod, r.Method)
	ctx = context.WithValue(ctx, models.URLPath, r.URL.Path)
	route := mux.CurrentRoute(r)
	path, err := route.GetPathTemplate()
	if err == nil {
		ctx = context.WithValue(ctx, models.URLPathTemplate, path)
	}

	ctx = utils.FetchContextFromHeaders(ctx, r)

	if r.Header.Get("Authorization") != "" {
		authHeaderStrings := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeaderStrings) > 1 {
			ctx = context.WithValue(ctx, models.JWTContextKey, authHeaderStrings[1])
		}
	}

	return ctx
}

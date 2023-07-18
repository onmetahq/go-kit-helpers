package validators

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/kit/endpoint"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestIPValidator(t *testing.T) {

	logger := initLogger()

	type input struct {
		ctx  context.Context
		next endpoint.Endpoint
	}

	type output struct {
		err error
	}

	validIps := "10.20.1.3,asdaD,223.34.34.1,SADFASDF"
	os.Setenv("apiKey", validIps)

	tests := []struct {
		name   string
		input  input
		output output
	}{
		{
			name: "No Merchant API key in the context",
			input: input{
				ctx: func() context.Context {
					ctx := context.TODO()
					return ctx
				}(),
				next: func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
			},
			output: output{
				err: models.ErrUnauthorized,
			},
		},
		{
			name: "No x-forwarded-for IP address in the request",
			input: input{
				ctx: func() context.Context {
					ctx := context.TODO()
					return context.WithValue(ctx, ctxKeys.MerchantAPIKey, "apiKey")
				}(),
				next: func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
			},
			output: output{
				err: nil,
			},
		},
		{
			name: "valid matching IP address",
			input: input{
				ctx: func() context.Context {
					ctx := context.TODO()
					ctx = context.WithValue(ctx, ctxKeys.XForwardedFor, "10.20.1.3")
					return context.WithValue(ctx, ctxKeys.MerchantAPIKey, "apiKey")
				}(),
				next: func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
			},
			output: output{
				err: nil,
			},
		},
		{
			name: "multiple matching IP address",
			input: input{
				ctx: func() context.Context {
					ctx := context.TODO()
					ctx = context.WithValue(ctx, ctxKeys.XForwardedFor, "10.20.1.3, 223.34.34.1")
					return context.WithValue(ctx, ctxKeys.MerchantAPIKey, "apiKey")
				}(),
				next: func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
			},
			output: output{
				err: nil,
			},
		},
		{
			name: "one matching IP address",
			input: input{
				ctx: func() context.Context {
					ctx := context.TODO()
					ctx = context.WithValue(ctx, ctxKeys.XForwardedFor, "10.20.5.3, 223.34.34.1")
					return context.WithValue(ctx, ctxKeys.MerchantAPIKey, "apiKey")
				}(),
				next: func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
			},
			output: output{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		_, err := IPValidator(logger)(test.input.next)(test.input.ctx, map[string]string{})
		assert.Equal(t, test.output.err, err, test.name)
	}
}

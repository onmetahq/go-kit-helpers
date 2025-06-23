package validators

import (
	"context"
	"testing"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	"github.com/stretchr/testify/assert"
)

const SECRET = "abcdefg"

type customerClaims struct {
	UserID   string
	TenantID string
	jwt.StandardClaims
}

func generateCustomerClaims(userId string, tenantId string, expirationTime time.Time) *customerClaims {
	return &customerClaims{
		UserID:   userId,
		TenantID: tenantId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
}

func generateJWTToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET))
}

func TestOptionalJWTValidator(t *testing.T) {
	type input struct {
		ctx  context.Context
		next endpoint.Endpoint
	}

	type output struct {
		userId string
		err    error
	}

	presentTime := time.Now()
	accessTokenExpirationTime := presentTime.Add(time.Minute * 60)
	token, _ := generateJWTToken(generateCustomerClaims("user1", "tenant1", accessTokenExpirationTime))

	tests := []struct {
		name   string
		input  input
		output output
	}{
		{
			name: "Empty context with no value in JWT context key",
			input: input{
				ctx: context.Background(),
				next: func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
			},
			output: output{
				err: nil,
			},
		},
		{
			name: "Empty JWT token",
			input: input{
				ctx: context.WithValue(context.Background(), models.JWTContextKey, ""),
				next: func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
			},
			output: output{
				err: models.ErrUnauthorized,
			},
		},
		{
			name: "Valid JWT token",
			input: input{
				ctx: context.WithValue(context.Background(), models.JWTContextKey, token),
				next: func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
			},
			output: output{
				userId: "user1",
				err:    nil,
			},
		},
	}

	for _, test := range tests {
		_, err := OptionalJWTValidator(SECRET)(test.input.next)(test.input.ctx, map[string]string{})
		assert.Equal(t, test.output.err, err)
	}
}

func TestJWTValidator(t *testing.T) {
	type input struct {
		ctx  context.Context
		next endpoint.Endpoint
	}

	type output struct {
		userId string
		err    error
	}

	presentTime := time.Now()
	accessTokenExpirationTime := presentTime.Add(time.Minute * 60)
	token, _ := generateJWTToken(generateCustomerClaims("user1", "tenant1", accessTokenExpirationTime))

	tests := []struct {
		name   string
		input  input
		output output
	}{
		{
			name: "Empty context with no value in JWT context key",
			input: input{
				ctx: context.Background(),
				next: func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
			},
			output: output{
				err: models.ErrUnauthorized,
			},
		},
		{
			name: "Empty JWT token",
			input: input{
				ctx: context.WithValue(context.Background(), models.JWTContextKey, ""),
				next: func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
			},
			output: output{
				err: models.ErrUnauthorized,
			},
		},
		{
			name: "Valid JWT token",
			input: input{
				ctx: context.WithValue(context.Background(), models.JWTContextKey, token),
				next: func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
			},
			output: output{
				userId: "user1",
				err:    nil,
			},
		},
	}

	for _, test := range tests {
		_, err := JWTValidator(SECRET)(test.input.next)(test.input.ctx, map[string]string{})
		assert.Equal(t, test.output.err, err)
	}
}

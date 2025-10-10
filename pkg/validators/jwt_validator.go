package validators

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	metahttp "github.com/onmetahq/meta-http/pkg/meta_http"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

type BlacklistRequest struct {
	AccessToken string `json:"accessToken"`
}

type BlacklistResponse struct {
	Success bool `json:"success"`
	Data    struct {
		IsBlacklisted bool `json:"isBlacklisted"`
	} `json:"data"`
}

func checkTokenBlacklist(ctx context.Context, tokenString string) (bool, error) {	
	fullURL := os.Getenv("BLACKLIST_API_URL")
	if fullURL == "" {
		return false, fmt.Errorf("BLACKLIST_API_URL environment variable is not set")
	}
	
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return false, fmt.Errorf("API_KEY environment variable is not set")
	}
	
	// Parse the URL to get base URL and path
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return false, fmt.Errorf("invalid BLACKLIST_API_URL: %w", err)
	}
	
	// Create base URL (scheme + host + port)
	baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	path := parsedURL.Path
	
	headers := map[string]string{
		"Content-Type": "application/json",
		"apikey":       apiKey,
	}
	
	payload := BlacklistRequest{AccessToken: tokenString}
	var response BlacklistResponse
	
	logger := slog.Default()
	client := metahttp.NewClient(baseURL, logger, 10*time.Second)
	
	_, err = client.Post(ctx, path, headers, payload, &response)
	if err != nil {
		return false, fmt.Errorf("failed to call blacklist API: %w", err)
	}
	
	return response.Success && response.Data.IsBlacklisted, nil
}

func JWTValidator(hmacSecret string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			tokenString, ok := ctx.Value(models.JWTContextKey).(string)
			if !ok {
				slog.ErrorContext(ctx, "JWT not found in context")
				return nil, models.ErrUnauthorized
			}

			claims := &models.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					slog.ErrorContext(ctx, "Invalid JWT header method", "error", token.Method.Alg())
					return nil, models.ErrUnexpectedSigningMethod
				}

				return []byte(hmacSecret), nil
			})

			if err != nil {
				if e, ok := err.(*jwt.ValidationError); ok {
					switch {
					case e.Errors&jwt.ValidationErrorMalformed != 0:
						slog.ErrorContext(ctx, "Malformed JWT")
					case e.Errors&jwt.ValidationErrorExpired != 0:
						slog.ErrorContext(ctx, "Expired JWT")
					case e.Errors&jwt.ValidationErrorNotValidYet != 0:
						slog.ErrorContext(ctx, "Inactive JWT")
					case e.Inner != nil:
						slog.ErrorContext(ctx, "Inner JWT error", "error", e.Inner)
					}
				}
				slog.ErrorContext(ctx, "JWT validation error", "error", err)
				return nil, models.ErrUnauthorized
			}

			if !token.Valid {
				slog.ErrorContext(ctx, "Invalid JWT")
				return nil, models.ErrUnauthorized
			}

			// Check if token is blacklisted
			isBlacklisted, err := checkTokenBlacklist(ctx, tokenString)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to check token blacklist", "error", err)
			} else if isBlacklisted {
				slog.ErrorContext(ctx, "Token is blacklisted")
				return nil, models.ErrUnauthorized
			}

			ctx = context.WithValue(ctx, models.JWTClaimsContextKey, claims)
			ctx = context.WithValue(ctx, ctxKeys.TenantID, claims.TenantID)
			ctx = context.WithValue(ctx, models.USERID, claims.UserId) // TODO: Remove it soon
			ctx = context.WithValue(ctx, ctxKeys.UserID, claims.UserId)

			return next(ctx, request)
		}
	}
}

package validators

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
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

func checkTokenBlacklist(tokenString string) (bool, error) {
	slog.Debug("Checking token blacklist", "token", tokenString)
	url := os.Getenv("BLACKLIST_API_URL")
	if url == "" {
		return false, fmt.Errorf("BLACKLIST_API_URL environment variable is not set")
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return false, fmt.Errorf("API_KEY environment variable is not set")
	}
	
	payload, _ := json.Marshal(BlacklistRequest{AccessToken: tokenString})
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return false, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}
	
	var blacklistResp BlacklistResponse
	if err := json.NewDecoder(resp.Body).Decode(&blacklistResp); err != nil {
		return false, err
	}
	
	return blacklistResp.Success && blacklistResp.Data.IsBlacklisted, nil
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
			isBlacklisted, err := checkTokenBlacklist(tokenString)
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

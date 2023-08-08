package models

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	// JWTContextKey holds the key used to store a JWT in the context.
	JWTContextKey        contextKey = "JWTToken"
	JWTTokenContextKey              = JWTContextKey
	JWTClaimsContextKey  contextKey = "JWTClaims"
	PathParamsContextKey contextKey = "PathParams"
	USERID               contextKey = "UserID"
	URLPath              contextKey = "URLPath"
	HttpMethod           contextKey = "HttpMethod"
	URLPathTemplate      contextKey = "URLPathTemplate"
)

type Claims struct {
	TenantID string `json:"tenantId"`
	APIKey   string `json:"apiKey"`
	UserId   string `json:"userId"`
	jwt.StandardClaims
}

var (
	// ErrTokenContextMissing denotes a token was not passed into the parsing
	// middleware's context.
	ErrTokenContextMissing = errors.New("token up for parsing was not passed through the context")

	// ErrTokenInvalid denotes a token was not able to be validated.
	ErrTokenInvalid = errors.New("JWT was invalid")

	// ErrTokenExpired denotes a token's expire header (exp) has since passed.
	ErrTokenExpired = errors.New("JWT is expired")

	// ErrTokenMalformed denotes a token was not formatted as a JWT.
	ErrTokenMalformed = errors.New("JWT is malformed")

	// ErrTokenNotActive denotes a token's not before header (nbf) is in the
	// future.
	ErrTokenNotActive = errors.New("token is not valid yet")

	// ErrUnexpectedSigningMethod denotes a token was signed with an unexpected
	// signing method.
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")

	//Unauthorized error
	ErrUnauthorized = errors.New("unauthorized to access")

	ErrNotFound = errors.New("data not found")
)

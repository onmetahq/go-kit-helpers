package utilities

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// validating a temporary API token to extract tenant ID
func ValidateTempAPI(tokenString string, secret string) (string, error) {
	var secretKey = []byte(secret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return "", fmt.Errorf("token has expired")
			}
		}
		if tid, ok := claims["tenantId"].(string); ok {
			return tid, nil
		}
	}

	return "", fmt.Errorf("invalid token")
}

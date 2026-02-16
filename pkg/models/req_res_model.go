package models

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	http_models "github.com/onmetahq/meta-http/pkg/models"
)

type Response struct {
	Success    bool      `json:"success"`
	Message    string    `json:"message,omitempty"`
	Data       any       `json:"data,omitempty"`
	Error      ErrorInfo `json:"error,omitempty"`
	XRequestid string    `json:"X-Request-ID,omitempty"`
}

type ErrorInfo struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrResponse struct {
	Success bool      `json:"success"`
	Error   ErrorInfo `json:"error"`
}

func EncodeJSONResponseWithRequestID(ctx context.Context, w http.ResponseWriter, response any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	requestID := ""
	if id, ok := ctx.Value(http_models.RequestID).(string); ok {
		requestID = id
	}

	var finalResponse any

	switch resp := response.(type) {
	case Response:
		finalResponse = map[string]any{
			"success":   resp.Success,
			"requestId": requestID,
		}
		if resp.Message != "" {
			finalResponse.(map[string]any)["message"] = resp.Message
		}
		if resp.Data != nil {
			finalResponse.(map[string]any)["data"] = resp.Data
		}

	case ErrResponse:
		finalResponse = map[string]any{
			"success":   resp.Success,
			"error":     resp.Error,
			"requestId": requestID,
		}

	case map[string]any:
		resp["requestId"] = requestID
		finalResponse = resp

	default:
		finalResponse = map[string]any{
			"success":   true,
			"data":      response,
			"requestId": requestID,
		}
	}

	return json.NewEncoder(w).Encode(finalResponse)
}

func EnsureRequestID(ctx context.Context, r *http.Request) context.Context {
	requestID := extractRequestIDFromHeader(r)

	if requestID == "" {
		requestID = extractRequestIDFromBody(r)
	}

	if requestID == "" {
		requestID = uuid.NewString()
	}

	return context.WithValue(ctx, http_models.RequestID, requestID)
}

func extractRequestIDFromHeader(r *http.Request) string {
	for key, values := range r.Header {
		normalized := strings.ToLower(strings.ReplaceAll(key, "_", "-"))

		if normalized == "x-request-id" {
			for _, v := range values {
				if v = strings.TrimSpace(v); v != "" {
					return v
				}
			}
		}
	}
	return ""
}

func extractRequestIDFromBody(r *http.Request) string {
	if r == nil || r.Body == nil {
		return ""
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil || len(bodyBytes) == 0 {
		return ""
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var payload map[string]any
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return ""
	}

	for key, value := range payload {
		normalized := strings.ToLower(strings.ReplaceAll(key, "_", "-"))

		if normalized == "x-request-id" || normalized == "request-id" {
			if v, ok := value.(string); ok {
				if v = strings.TrimSpace(v); v != "" {
					return v
				}
			}
		}
	}

	return ""
}

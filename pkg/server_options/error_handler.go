package options

import (
	"context"
	"log/slog"

	"github.com/go-kit/kit/transport"
)

type errorLogger struct {
}

func NewErrorLogger() transport.ErrorHandler {
	return &errorLogger{}
}

func (l *errorLogger) Handle(ctx context.Context, err error) {
	slog.ErrorContext(ctx, "Error occurred in transport layer", "error", err)
}

package options

import (
	"context"

	"github.com/go-kit/kit/transport"
	"github.com/onmetahq/go-kit-helpers/pkg/logger"
)

type errorLogger struct {
	logger logger.CtxLogger
}

func NewErrorLogger(logger logger.CtxLogger) transport.ErrorHandler {
	return &errorLogger{
		logger: logger,
	}
}

func (l *errorLogger) Handle(ctx context.Context, err error) {
	l.logger.Context(ctx).Error().Log("err", err)
}

package options

import (
	"context"

	"github.com/go-kit/kit/transport"
	"github.com/go-kit/log"
	"github.com/onmetahq/go-kit-helpers/pkg/logger"
)

type errorLogger struct {
	logger log.Logger
}

func NewErrorLogger(logger log.Logger) transport.ErrorHandler {
	return &errorLogger{
		logger: logger,
	}
}

func (l *errorLogger) Handle(ctx context.Context, err error) {
	logger.NewCtxLogger(l.logger).Context(ctx).Error().Log("err", err)
}

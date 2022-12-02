package logger

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/onmetahq/meta-http/pkg/utils"
)

type CtxLogger interface {
	Context(ctx context.Context) CtxLogger

	Level(f func(log.Logger) log.Logger) log.Logger
	Info() log.Logger
	Error() log.Logger
	Warn() log.Logger
	Debug() log.Logger

	Logger() log.Logger
}

type customLogger struct {
	logger log.Logger
}

func NewCtxLogger(logger log.Logger) CtxLogger {
	return &customLogger{
		logger: logger,
	}
}

func (l *customLogger) Context(ctx context.Context) CtxLogger {
	data := utils.FetchHeadersFromContext(ctx)
	args := []interface{}{}
	for key, value := range data {
		args = append(args, key)
		args = append(args, value)
	}
	l.logger = log.With(l.logger, args...)
	return l
}

func (l *customLogger) Info() log.Logger {
	return level.Info(l.logger)
}

func (l *customLogger) Error() log.Logger {
	return level.Error(l.logger)
}

func (l *customLogger) Warn() log.Logger {
	return level.Warn(l.logger)
}

func (l *customLogger) Debug() log.Logger {
	return level.Debug(l.logger)
}

func (l *customLogger) Level(f func(log.Logger) log.Logger) log.Logger {
	return f(l.logger)
}

func (l *customLogger) Logger() log.Logger {
	return l.logger
}

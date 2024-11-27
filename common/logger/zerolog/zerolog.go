package zerolog

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

const (
	ContextKeyRequestId = "requestID"
	LogKeyRequestId     = "request_id"
	ContextKeyEventId   = "eventID"
	LogKeyEventId       = "event_id"
)

type Zerolog struct {
	logger zerolog.Logger
}

func NewZeroLog() Zerolog {
	return Zerolog{
		logger: zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
}

func (z Zerolog) Error(err error, message string) {
	z.logger.Error().Err(err).Msg(message)
}

func (z Zerolog) ErrorWithContext(ctx context.Context, err error, message string) {
	logger := z.withContext(ctx).logger
	logger.Error().Err(err).Msg(message)
}

func (z Zerolog) ErrorfWithContext(ctx context.Context, err error, format string, v ...interface{}) {
	logger := z.withContext(ctx).logger
	logger.Error().Err(err).Msgf(format, v...)
}

func (z Zerolog) WarnfWithContext(ctx context.Context, err error, format string, v ...interface{}) {
	logger := z.withContext(ctx).logger
	logger.Warn().Err(err).Msgf(format, v...)
}

func (z Zerolog) InfofWithContext(ctx context.Context, format string, v ...interface{}) {
	logger := z.withContext(ctx).logger
	logger.Info().Msgf(format, v...)
}

func (z Zerolog) Info(message string) {
	z.logger.Info().Msg(message)
}

func (z Zerolog) Infof(format string, v ...interface{}) {
	z.logger.Info().Msgf(format, v...)
}

func (z Zerolog) WithHandlerName(name string) Zerolog {
	z.logger = z.logger.With().Str("handlerName", name).Logger()
	return z
}

func (z Zerolog) WithServiceName(name string) Zerolog {
	z.logger = z.logger.With().Str("serviceName", name).Logger()
	return z
}

func (z Zerolog) WithRepositoryName(name string) Zerolog {
	z.logger = z.logger.With().Str("repositoryName", name).Logger()
	return z
}

func (z Zerolog) withContext(ctx context.Context) Zerolog {
	if value, ok := ctx.Value(ContextKeyRequestId).(string); ok {
		z.logger = z.logger.With().Str(LogKeyRequestId, value).Logger()
	}

	if value, ok := ctx.Value(ContextKeyEventId).(string); ok {
		z.logger = z.logger.With().Str(LogKeyEventId, value).Logger()
	}

	return z
}

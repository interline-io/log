package log

import (
	"context"

	"github.com/rs/zerolog"
)

func For(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

func WithLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return logger.WithContext(ctx)
}

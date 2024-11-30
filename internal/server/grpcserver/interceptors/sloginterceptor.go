package interceptors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func NewSlogInterceptor() grpc.UnaryServerInterceptor {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}
	slogger := logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		slog.Default().Log(ctx, slog.Level(lvl), msg, fields...)
	})
	return logging.UnaryServerInterceptor(slogger, opts...)
}

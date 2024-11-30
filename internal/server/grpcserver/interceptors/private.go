package interceptors

import (
	"context"
	"google.golang.org/grpc/metadata"
)

// getFirstMd silently gets first metadata value as string
func getFirstMd(ctx context.Context, key string) (s string) {
	values := metadata.ValueFromIncomingContext(ctx, key)
	if len(values) > 0 {
		s = values[0]
	}
	return
}

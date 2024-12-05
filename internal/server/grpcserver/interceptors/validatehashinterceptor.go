package interceptors

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/hasher"
)

var (
	hashPresentKey = "hashpresentkey"
)

func NewValidateHashInterceptor(key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if key != "" {
			hHash := getFirstMd(ctx, constants.HeaderHashSHA256)
			rHash := hasher.Hash(req.(proto.Message), []byte(key))
			if hHash != rHash {
				slog.Error("metadata value is invalid", "realHash", rHash, "requestHash", hHash)
				msg := fmt.Sprintf("%s header value is invalid: %s", constants.HeaderHashSHA256, hHash)
				return nil, status.Error(codes.Unauthenticated, msg)
			}
			ctx = makeContextWithHashValidMark(ctx)
		}
		return handler(ctx, req)
	}
}

func makeContextWithHashValidMark(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	md.Set(hashPresentKey, "1")
	ctx = metadata.NewIncomingContext(ctx, md)
	return ctx
}

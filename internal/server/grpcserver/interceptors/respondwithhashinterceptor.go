package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/cryptofacade"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/hasher"
)

func NewRespondWithHashInterceptor(key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if key != "" && getIsContextWithHashValidMark(ctx) {
			var rHash string
			if err != nil {
				rHash = cryptofacade.Hash([]byte(err.Error()), []byte(key))
			} else {
				rHash = hasher.Hash(req.(proto.Message), []byte(key))
			}
			md := metadata.New(map[string]string{
				constants.HeaderHashSHA256: rHash,
			})
			_ = grpc.SetTrailer(ctx, md)
		}
		return
	}
}

func getIsContextWithHashValidMark(ctx context.Context) bool {
	values := metadata.ValueFromIncomingContext(ctx, hashPresentKey)
	return len(values) > 0
}

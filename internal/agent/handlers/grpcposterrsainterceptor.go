package handlers

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/rsacrypt"
)

func newGrpcPosterRSAInterceptor(pub *rsa.PublicKey) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if pub != nil {
			var err error
			var buf []byte
			m, ok := req.(proto.Message)
			if !ok {
				err = fmt.Errorf("request is not a proto.Message: %T", req)
			}
			if err == nil {
				buf, err = proto.Marshal(m)
			}
			if err == nil {
				buf, err = rsacrypt.EncryptWithPublicKey(buf, pub)
			}
			if err == nil {
				req = &pb.UpdateMetricsRequest{Raw: buf}
			}
			if err != nil {
				slog.Info("newGrpcPosterRSAInterceptor error", "err", err)
				return status.Errorf(codes.Internal, "newGrpcPosterRSAInterceptor")
			}
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

package interceptors

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

func NewDecryptBodyInterceptor(priv *rsa.PrivateKey) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if priv != nil {
			var buf []byte
			m, ok := req.(*pb.UpdateMetricsRequest)
			if !ok {
				err = fmt.Errorf("request is not a *pb.UpdateMetricsRequest: %T", req)
			}
			if err == nil && (m == nil || m.Raw == nil || len(m.Raw) == 0) {
				err = fmt.Errorf("raw request body is empty")
			}
			if err == nil {
				buf, err = rsacrypt.DecryptWithPrivateKey(m.Raw, priv)
			}
			if err == nil {
				err = proto.Unmarshal(buf, m)
			}
			if err == nil {
				req = m
			}
			if err != nil {
				slog.Info("body decryption error", "error", err)
				return nil, status.Errorf(codes.Unauthenticated, "body decryption error")
			}
		}
		return handler(ctx, req)
	}
}

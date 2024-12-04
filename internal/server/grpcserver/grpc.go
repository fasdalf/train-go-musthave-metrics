package grpcserver

import (
	"crypto/rsa"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"

	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
	gc "github.com/fasdalf/train-go-musthave-metrics/internal/server/grpcserver/controller"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/grpcserver/interceptors"
	hh "github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
)

func NewGrpcServer(ms hh.Storage, db hh.Pingable, retryer hh.Retryer, key string, decryptionKey *rsa.PrivateKey, tr *net.IPNet) *grpc.Server {
	options := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors.NewSlogInterceptor()),
		grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor()),
		// IRL use github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/realip
		// or at least peer.FromContext(ctx).Addr
		grpc.ChainUnaryInterceptor(interceptors.NewValidateIPInterceptor(tr)),
		grpc.ChainUnaryInterceptor(interceptors.NewValidateHashInterceptor(key)),
		grpc.ChainUnaryInterceptor(interceptors.NewRespondWithHashInterceptor(key)),
		// TODO: ##@@ add gzip+rsa home-made compressor and validate metadata for "grpc-accept-encoding": ["gziprsa"]
		// This one grpc.ForceServerCodecV2(),
		// or use interceptor with decrypt from message body
	}

	s := grpc.NewServer(options...)
	mServer := gc.NewMetricsServer(ms, retryer)
	pb.RegisterMetricsServer(s, mServer)

	return s
}

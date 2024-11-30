package grpcserver

import (
	"crypto/rsa"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"

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
		// TODO: ##@@ add hash, gzip, rsa interceptors
	}

	s := grpc.NewServer(options...)
	// TODO: ##@@ add retryer
	mServer := &gc.MetricsServer{
		Storage: ms,
	}
	pb.RegisterMetricsServer(s, mServer)

	return s
}

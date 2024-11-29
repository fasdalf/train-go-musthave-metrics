package grpcserver

import (
	"crypto/rsa"
	"net"

	"google.golang.org/grpc"

	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
	gc "github.com/fasdalf/train-go-musthave-metrics/internal/server/grpcserver/controller"
	hh "github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
)

func NewGrpcServer(ms hh.Storage, db hh.Pingable, retryer hh.Retryer, key string, decryptionKey *rsa.PrivateKey, tr *net.IPNet) *grpc.Server {
	s := grpc.NewServer()
	mServer := &gc.MetricsServer{
		Storage: ms,
	}
	// TODO: ##@@ add retryer, hash, gzip, rsa, trusted subnet interceptors
	pb.RegisterMetricsServer(s, mServer)

	return s
}

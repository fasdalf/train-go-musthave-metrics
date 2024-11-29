package grpcserver

import (
	"crypto/rsa"
	"net"
	"testing"

	"google.golang.org/grpc"

	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
)

func TestNewGrpcServer_Created(t *testing.T) {
	type args struct {
		ms            handlers.Storage
		db            handlers.Pingable
		retryer       handlers.Retryer
		key           string
		decryptionKey *rsa.PrivateKey
		tr            *net.IPNet
	}
	tests := []struct {
		name string
		args args
		want *grpc.Server
	}{
		{
			name: "success",
			args: args{
				ms:            nil,
				db:            nil,
				retryer:       nil,
				key:           "",
				decryptionKey: nil,
				tr:            nil,
			},
			want: nil,
		},
		// TODO: ##@@ Add real test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGrpcServer(tt.args.ms, tt.args.db, tt.args.retryer, tt.args.key, tt.args.decryptionKey, tt.args.tr); got == nil {
				t.Error("NewGrpcServer() returned nil")
			}
		})
	}
}

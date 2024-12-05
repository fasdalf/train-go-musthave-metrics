package interceptors

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/rsacrypt"
)

func TestNewDecryptBodyInterceptor(t *testing.T) {
	m := &pb.UpdateMetricsRequest{}
	type args struct {
		ctx     context.Context
		req     any
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	}
	tests := []struct {
		name     string
		args     args
		wantResp any
		wantErr  bool
	}{
		{
			name:     "succeess",
			args:     args{req: m},
			wantResp: 10,
			wantErr:  false,
		},
		{
			name:     "error",
			args:     args{req: "invalid request"},
			wantResp: nil,
			wantErr:  true,
		},
		{
			name:     "empty",
			args:     args{req: &pb.UpdateMetricsRequest{}},
			wantResp: nil,
			wantErr:  true,
		},
	}
	usi := grpc.UnaryServerInfo{}
	nh := func(ctx context.Context, req any) (any, error) {
		return 10, nil
	}
	ctx := context.Background()
	priv, pub := rsacrypt.GenerateKeyPair(2048)
	interceptor := NewDecryptBodyInterceptor(priv)
	buf, _ := proto.Marshal(&pb.UpdateMetricsRequest{Metrics: []*pb.Metric{{Id: "mock"}}})
	buf, _ = rsacrypt.EncryptWithPublicKey(buf, pub)
	m.Raw = buf

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := interceptor(ctx, tt.args.req, &usi, nh)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDecryptBodyInterceptor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("NewDecryptBodyInterceptor() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

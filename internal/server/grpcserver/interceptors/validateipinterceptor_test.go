package interceptors

import (
	"context"
	"net"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
)

func TestNewValidateIPInterceptor(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		wantResp any
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.New(map[string]string{
						constants.HeaderRealIP: "192.168.50.48",
					}),
				),
			},
			wantResp: 10,
			wantErr:  false,
		},
		{
			name: "success",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.New(map[string]string{
						constants.HeaderRealIP: "8.8.8.8",
					}),
				),
			},
			wantResp: nil,
			wantErr:  true,
		},
	}
	_, tr, _ := net.ParseCIDR("192.168.50.0/24")
	usi := grpc.UnaryServerInfo{}
	nh := func(ctx context.Context, req any) (any, error) {
		return 10, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewValidateIPInterceptor(tr)
			gotResp, err := got(tt.args.ctx, nil, &usi, nh)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIPInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ValidateIPInterceptor() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

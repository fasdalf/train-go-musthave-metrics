package interceptors

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc"
)

func TestNewSlogInterceptor(t *testing.T) {
	usi := grpc.UnaryServerInfo{}
	nh := func(ctx context.Context, req any) (any, error) {
		return 10, nil
	}

	got := NewSlogInterceptor()
	if got == nil {
		t.Errorf("NewSlogInterceptor() is nil")
	}
	gotResp, err := got(context.Background(), nil, &usi, nh)
	if err != nil {
		t.Errorf("NewSlogInterceptor() error = %v, should be empty", err)
		return
	}
	if !reflect.DeepEqual(gotResp, 10) {
		t.Errorf("NewSlogInterceptor() gotResp = %v, want %v", gotResp, 10)
	}
}

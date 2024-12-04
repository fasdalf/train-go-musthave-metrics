package hasher

import (
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
)

func TestHash(t *testing.T) {
	type args struct {
		m   proto.Message
		key []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "one",
			args: args{&pb.Metric{Value: 20}, []byte("key")},
			want: "6b657975b3480986181d640ebe11e89a92b140733d11693f86f05f1c4af079f21c460a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.args.m, tt.args.key); got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

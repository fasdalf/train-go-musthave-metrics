package interceptors

import (
	"context"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/localip"
)

func NewValidateIPInterceptor(tr *net.IPNet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if tr != nil {
			ipString := getFirstMd(ctx, constants.HeaderRealIP)
			if err := localip.ValidateIPStringInSubnet(ipString, tr); err != nil {
				slog.Error("metadata value is invalid", "header", constants.HeaderRealIP, "value", ipString, "error", err)
				return nil, status.Error(codes.Unauthenticated, "not a trusted IP")
			}
		}
		return handler(ctx, req)
	}
}

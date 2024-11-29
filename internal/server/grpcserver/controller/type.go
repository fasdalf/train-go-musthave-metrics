package controller

import (
	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
	hh "github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
)

type Foo struct {
	Bar int
}

// MetricsServer поддерживает все необходимые методы сервера.
type MetricsServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedMetricsServer
	Storage hh.Storage
}

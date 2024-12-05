package hasher

import (
	"google.golang.org/protobuf/proto"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/cryptofacade"
)

func Hash(m proto.Message, key []byte) string {
	buf, _ := proto.Marshal(m)
	return cryptofacade.Hash(buf, key)
}

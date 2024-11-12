package mockreaders

import (
	"bytes"
	"io"
)

// ErrCloser — специальная структура, эмулирующая ошибку закрытия
type ErrCloser struct {
	*bytes.Buffer
}

func NewErrCloser(buf []byte) *ErrCloser {
	return &ErrCloser{bytes.NewBuffer(buf)}
}
func (e *ErrCloser) Close() error {
	return io.ErrClosedPipe
}

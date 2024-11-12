package mockreaders

import (
	"errors"
)

// ErrReader — специальная структура, эмулирующая ошибку чтения
type ErrReader struct{}

func (e *ErrReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("эмулированная ошибка чтения")
}

func (e *ErrReader) Close() error {
	return nil
}

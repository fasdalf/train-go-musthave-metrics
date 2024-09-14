package cryptofacade

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(body, key []byte) string {
	// создаём новый hash.Hash, вычисляющий контрольную сумму SHA-256
	h := sha256.New()
	// передаём байты для хеширования
	h.Write(body)
	// получаем хеш в виде строки
	return hex.EncodeToString(h.Sum(key))
}

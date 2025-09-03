package utils

import (
	"hash/crc32"
	"math"
)

func HashInRange(key string) float32 {
	hash := crc32.ChecksumIEEE([]byte(key))

	return float32(hash) / float32(math.MaxUint32)
}

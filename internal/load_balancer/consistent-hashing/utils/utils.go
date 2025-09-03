package utils

import (
	"hash/crc32"
	"math"

	"github.com/Marie20767/load-balancer/internal/load_balancer/consistent-hashing/config"
)

func HashInRange(key string) float32 {
	hash := crc32.ChecksumIEEE([]byte(key))

	return (float32(hash) / float32(math.MaxUint32)) * config.HashRingValue
}

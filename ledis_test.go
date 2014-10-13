package cache

import (
	"testing"
)

func BenchmarkLedisGet(b *testing.B) {
	MakeBenchmarkGet(LedisProviderDefault)(b)
}

func BenchmarkLedisSet(b *testing.B) {
	MakeBenchmarkSet(LedisProviderDefault)(b)
}

func BenchmarkLedisSetParallel(b *testing.B) {
	MakeBenchmarkSetParallel(LedisProviderDefault)(b)
}

func BenchmarkLedisToRedisGet(b *testing.B) {
	MakeBenchmarkGet(LedisProviderToRedisDefault)(b)
}

func BenchmarkLedisToRedisSet(b *testing.B) {
	MakeBenchmarkSet(LedisProviderToRedisDefault)(b)
}

func BenchmarkLedisToRedisParallel(b *testing.B) {
	MakeBenchmarkSetParallel(LedisProviderToRedisDefault)(b)
}

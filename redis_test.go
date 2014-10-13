package cache

import (
	"testing"
)

func BenchmarkRedisGet(b *testing.B) {
	MakeBenchmarkGet(RedisProviderDefault)(b)
}

func BenchmarkRedisSet(b *testing.B) {
	MakeBenchmarkSet(RedisProviderDefault)(b)
}

func BenchmarkRedisSetParallel(b *testing.B) {
	MakeBenchmarkSetParallel(RedisProviderDefault)(b)
}

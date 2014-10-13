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

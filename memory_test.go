package cache

import (
	"testing"
)

func BenchmarkMemoryGet(b *testing.B) {
	MakeBenchmarkGet(MemoryProvider)(b)
}

func BenchmarkMemorySet(b *testing.B) {
	MakeBenchmarkSet(MemoryProvider)(b)
}

func BenchmarkMemorySetParallel(b *testing.B) {
	MakeBenchmarkSetParallel(MemoryProvider)(b)
}

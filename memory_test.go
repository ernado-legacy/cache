package cache

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMemory(t *testing.T) {
	Convey("Memory", t, func() {
		memory := MemoryProvider()
		Convey("Set", func() {
			v := "data"
			key := "key"
			So(memory.Set(key, v), ShouldBeNil)
			Convey("Get", func() {
				var value string
				So(memory.Get(key, &value), ShouldBeNil)
				So(value, ShouldEqual, v)
			})
			Convey("Remove", func() {
				So(memory.Remove(key), ShouldBeNil)
				So(memory.Get(key, &v), ShouldEqual, ErrorNotExist)
			})
			Convey("Wrong type", func() {
				var value int
				So(memory.Get(key, value), ShouldNotBeNil)
				So(memory.Get(key, &value), ShouldNotBeNil)
			})
		})
	})
}

func BenchmarkMemorySet(b *testing.B) {
	m := MemoryProvider()
	for i := 0; i < b.N; i++ {
		m.Set("key", "data")
	}
}

func BenchmarkMemorySetParallel(b *testing.B) {
	m := MemoryProvider()
	b.RunParallel((func(pb *testing.PB) {
		for pb.Next() {
			m.Set("key", "data")
		}
	}))
}

func BenchmarkMemoryGet(b *testing.B) {
	m := MemoryProvider()
	m.Set("key", "data")
	var v string
	for i := 0; i < b.N; i++ {
		m.Set("key", &v)
	}
}

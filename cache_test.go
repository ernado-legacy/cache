package cache

import (
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestCache(t *testing.T) {
	Convey("Cache", t, func() {
		TestProvider := func(f func() Provider) {
			provider := f()
			name := reflect.ValueOf(provider).Elem().Type().Name()
			Convey(name, func() {
				Convey("Set", func() {
					v := "testing:data"
					key := "key"
					So(provider.Set(key, v), ShouldBeNil)
					Convey("Set structure", func() {
						type Data struct {
							A int
							B string
						}
						v := Data{6, "data"}
						So(provider.Set(key, v), ShouldBeNil)
					})
					Convey("Get", func() {
						var value string
						So(provider.Get(key, &value), ShouldBeNil)
						So(value, ShouldEqual, v)
					})
					Convey("Remove", func() {
						So(provider.Remove(key), ShouldBeNil)
						So(provider.Get(key, &v), ShouldEqual, ErrorNotExist)
					})
					Convey("Wrong type", func() {
						var value int
						So(provider.Get(key, value), ShouldNotBeNil)
						So(provider.Get(key, &value), ShouldNotBeNil)
					})
				})
			})
		}
		TestProvider(MemoryProvider)
		TestProvider(RedisProviderDefault)
		TestProvider(LedisProviderDefault)
	})
}

func MakeBenchmarkSet(callback func() Provider) func(*testing.B) {
	provider := callback()
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			provider.Set("key", "data")
		}
	}
}

func MakeBenchmarkSetParallel(callback func() Provider) func(*testing.B) {
	provider := callback()
	return func(b *testing.B) {
		b.RunParallel((func(pb *testing.PB) {
			for pb.Next() {
				provider.Set("key", "data")
			}
		}))
	}
}

func MakeBenchmarkGet(callback func() Provider) func(*testing.B) {
	provider := callback()
	return func(b *testing.B) {
		provider.Set("key", "data")
		var v string
		for i := 0; i < b.N; i++ {
			provider.Set("key", &v)
		}
	}
}

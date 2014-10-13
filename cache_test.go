package cache

import (
	"bytes"
	"encoding/gob"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestReflection(t *testing.T) {
	Convey("Reflection", t, func() {
		var (
			buff      = new(bytes.Buffer)
			encoder   = gob.NewEncoder(buff)
			decoder   = gob.NewDecoder(buff)
			value     string
			v         interface{} = &value
			valueType             = reflect.TypeOf(v).Elem()
			valuePtr              = reflect.New(valueType)
		)
		So(encoder.Encode(value), ShouldBeNil)
		So(decoder.Decode(valuePtr.Interface()), ShouldBeNil)
		So(valuePtr.Elem().Interface().(string), ShouldEqual, value)
	})
}

func TestCache(t *testing.T) {
	Convey("Cache", t, func() {
		TestProvider := func(f func() Provider) {
			provider := f()
			var name string
			providerValue := reflect.ValueOf(provider)
			if providerValue.Type().Kind() == reflect.Ptr {
				name = providerValue.Elem().Type().Name()
			} else {
				name = providerValue.Type().Name()
			}
			Convey(name, func() {
				Convey("Set "+name, func() {
					v := "testing:data:" + name
					key := v
					So(provider.Set(key, v), ShouldBeNil)
					Convey("Set structure "+name, func() {
						type Data struct {
							A int
							B string
							C []byte
						}
						v := Data{6, "data", []byte("test")}
						So(provider.Set(key, v), ShouldBeNil)
						value := new(Data)
						So(provider.Get(key, value), ShouldBeNil)
						So(value.A, ShouldEqual, v.A)
						So(value.B, ShouldEqual, v.B)
						So(string(value.C), ShouldEqual, string(v.C))
					})
					Convey("Get "+name, func() {
						var value string
						So(provider.Get(key, &value), ShouldBeNil)
						So(value, ShouldEqual, v)
					})
					Convey("Remove "+name, func() {
						So(provider.Remove(key), ShouldBeNil)
						So(provider.Get(key, &v), ShouldEqual, ErrorNotExist)
					})
					Convey("Wrong type "+name, func() {
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
		TestProvider(newClientAsProviderDefault)
		TestProvider(LedisProviderToRedisDefault)
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

func BenchmarkClientGet(b *testing.B) {
	MakeBenchmarkGet(newClientAsProviderDefault)(b)
}

func BenchmarkClientSet(b *testing.B) {
	MakeBenchmarkSet(newClientAsProviderDefault)(b)
}

func BenchmarkClientSetParallel(b *testing.B) {
	MakeBenchmarkSetParallel(newClientAsProviderDefault)(b)
}

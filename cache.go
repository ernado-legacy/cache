package cache

import (
	"errors"
	"reflect"
)

var (
	ErrorNotExist        = errors.New("Key not exists")
	ErrorShouldBePointer = errors.New("Should be pointer")
	ErrorInvalidType     = errors.New("Unable to set value: invalid type")
	ErrorNoProviders     = errors.New("No cache backends")
	ErrorInvalidProvider = errors.New("Invalid provider")
)

type Provider interface {
	Get(key string, v interface{}) error
	Set(key string, v interface{}) error
	TTL(key string, ttl uint64) error
	Remove(key string) error
}

type Client interface {
	Provider
	AddProvider(p interface{}) error
}

type defaultClient struct {
	providers []Provider
}

func (c *defaultClient) AddProvider(i interface{}) error {
	p, ok := i.(Provider)
	if ok {
		c.providers = append(c.providers, p)
		return nil
	}
	callback, ok := i.(func() Provider)
	if ok {
		c.providers = append(c.providers, callback())
		return nil
	}
	return ErrorInvalidProvider
}

func NewClient() Client {
	return new(defaultClient)
}

func newClientAsProviderDefault() Provider {
	c := NewClient()
	c.AddProvider(MemoryProvider)
	c.AddProvider(RedisProviderDefault)
	c.AddProvider(LedisProviderDefault)
	c.AddProvider(LedisProviderToRedisDefault)
	return c
}

func (c *defaultClient) Get(key string, v interface{}) (err error) {
	type Result struct {
		Value interface{}
		Error error
	}
	var (
		count = len(c.providers)
		rv    = reflect.ValueOf(v)
	)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrorShouldBePointer
	}
	if count == 0 {
		return ErrorNoProviders
	}
	var (
		results   = make(chan Result)
		valueType = reflect.TypeOf(v).Elem()
	)
	for _, provider := range c.providers {
		go func() {
			value := reflect.New(valueType).Interface()
			results <- Result{value, provider.Get(key, value)}
		}()
	}
	for i := 0; i < count; i++ {
		if data := <-results; data.Error == nil {
			rv.Elem().Set(reflect.ValueOf(data.Value).Elem())
			return
		} else {
			err = data.Error
		}
	}
	return err
}

func (c *defaultClient) Remove(key string) (err error) {
	if len(c.providers) == 0 {
		return ErrorNoProviders
	}
	for _, provider := range c.providers {
		err = provider.Remove(key)
		if err == ErrorNotExist {
			continue
		}
		if err != nil {
			return
		}
	}
	return nil
}

func (c *defaultClient) Set(key string, v interface{}) error {
	var (
		count = len(c.providers)
		errs  = make(chan error, count)
		err   error
	)
	if count == 0 {
		return ErrorNoProviders
	}
	for _, provider := range c.providers {
		go func() {
			errs <- provider.Set(key, v)
		}()
	}
	for i := 0; i < count; i++ {
		if err = <-errs; err != nil {
			return err
		}
	}
	return nil
}

func (c *defaultClient) TTL(key string, ttl uint64) (err error) {
	if len(c.providers) == 0 {
		return ErrorNoProviders
	}
	for _, provider := range c.providers {
		if err = provider.TTL(key, ttl); err != nil {
			return
		}
	}
	return
}

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
)

type Provider interface {
	Get(key string, v interface{}) error
	Set(key string, v interface{}) error
	TTL(key string, ttl uint64) error
	Remove(key string) error
}

type Client interface {
	Provider
	AddProvider(Provider)
}

type DefaultClient struct {
	providers []Provider
}

func (c *DefaultClient) AddProvider(p Provider) {
	c.providers = append(c.providers, p)
}

func NewClient() Client {
	return new(DefaultClient)
}

func newClientAsProviderDefault() Provider {
	c := NewClient()
	c.AddProvider(MemoryProvider())
	c.AddProvider(RedisProviderDefault())
	c.AddProvider(LedisProviderDefault())
	return c
}

func (c *DefaultClient) Get(key string, v interface{}) error {
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
			err := provider.Get(key, value)
			results <- Result{value, err}
		}()
	}
	var err error
	for i := 0; i < count; i++ {
		data := <-results
		if data.Error == nil {
			rv.Elem().Set(reflect.ValueOf(data.Value).Elem())
			return nil
		}
		err = data.Error
	}
	return err
}

func (c *DefaultClient) Remove(key string) (err error) {
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

func (c *DefaultClient) Set(key string, v interface{}) error {
	var (
		count = len(c.providers)
		errs  = make(chan error, count)
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
		err := <-errs
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DefaultClient) TTL(key string, ttl uint64) (err error) {
	if len(c.providers) == 0 {
		return ErrorNoProviders
	}
	for _, provider := range c.providers {
		err = provider.TTL(key, ttl)
		if err != nil {
			return
		}
	}
	return
}

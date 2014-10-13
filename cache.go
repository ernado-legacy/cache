package cache

import (
	"errors"
	// "sync"
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
	c := new(DefaultClient)
	return c
}

func newClientAsProviderDefault() Provider {
	c := NewClient()
	c.AddProvider(MemoryProvider())
	c.AddProvider(RedisProviderDefault())
	c.AddProvider(LedisProviderDefault())
	return c
}

func (c *DefaultClient) Get(key string, v interface{}) (err error) {
	for _, provider := range c.providers {
		err = provider.Get(key, v)
		if err == ErrorNotExist {
			continue
		}
		return
	}
	if err == ErrorNotExist {
		return
	}
	return ErrorNoProviders
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
	count := len(c.providers)
	if count == 0 {
		return ErrorNoProviders
	}
	errs := make(chan error, count)
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

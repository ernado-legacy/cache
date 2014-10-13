package cache

import (
	"bytes"
	"encoding/gob"
	"github.com/siddontang/ledisdb/client/go/ledis"
)

type LedisCache struct {
	client *ledis.Client
}

func LedisProvider(config *ledis.Config) Provider {
	return &LedisCache{ledis.NewClient(config)}
}

func LedisProviderDefault() Provider {
	cfg := new(ledis.Config)
	cfg.Addr = "127.0.0.1:6380"
	cfg.MaxIdleConns = 5
	return LedisProvider(cfg)
}

func (c LedisCache) Set(key string, value interface{}) error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(value); err != nil {
		return err
	}
	_, err := c.client.Do("set", key, buffer.Bytes())
	return err
}

func (c LedisCache) Get(key string, v interface{}) error {
	data, err := ledis.Bytes(c.client.Do("get", key))
	if err == ledis.ErrNil {
		return ErrorNotExist
	}
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(v)
}

func (c LedisCache) Remove(key string) error {
	_, err := c.client.Do("del", key)
	if err == ledis.ErrNil {
		return ErrorNotExist
	}
	return err
}

func (c LedisCache) TTL(key string, ttl uint64) error {
	_, err := c.client.Do("expire", key, ttl)
	if err == ledis.ErrNil {
		return ErrorNotExist
	}
	return err
}

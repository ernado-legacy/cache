package cache

import (
	"bytes"
	"encoding/gob"
	"github.com/siddontang/ledisdb/client/go/ledis"
)

const (
	ledisDefaultAddr = "127.0.0.1:6380"
	ledisMaxIdle     = 5
)

type ledisCache struct {
	client *ledis.Client
}

func LedisProvider(config *ledis.Config) Provider {
	return ledisCache{ledis.NewClient(config)}
}

func LedisProviderDefault() Provider {
	cfg := new(ledis.Config)
	cfg.Addr = ledisDefaultAddr
	cfg.MaxIdleConns = ledisMaxIdle
	return LedisProvider(cfg)
}

func LedisProviderToRedisDefault() Provider {
	cfg := new(ledis.Config)
	cfg.Addr = redisDefaultAddr
	cfg.MaxIdleConns = ledisMaxIdle
	return LedisProvider(cfg)
}

func (c ledisCache) Set(key string, value interface{}) error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(value); err != nil {
		return err
	}
	_, err := c.client.Do("set", key, buffer.Bytes())
	return err
}

func (c ledisCache) Get(key string, v interface{}) error {
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

func (c ledisCache) Remove(key string) error {
	_, err := c.client.Do("del", key)
	if err == ledis.ErrNil {
		return ErrorNotExist
	}
	return err
}

func (c ledisCache) TTL(key string, ttl uint64) error {
	_, err := c.client.Do("expire", key, ttl)
	if err == ledis.ErrNil {
		return ErrorNotExist
	}
	return err
}

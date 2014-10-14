package cache

import (
	"encoding/json"
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

func (c ledisCache) Set(key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = c.client.Do("set", key, data)
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
	return json.Unmarshal(data, v)
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

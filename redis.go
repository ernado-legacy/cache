package cache

import (
	"bytes"
	"encoding/gob"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisCache struct {
	pool *redis.Pool
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if len(password) != 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (c *RedisCache) Set(key string, v interface{}) error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(v); err != nil {
		return err
	}
	conn := c.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("SET", key, buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) Get(key string, v interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		return ErrorNotExist
	}
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(v)
}

func (c *RedisCache) Remove(key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	if err == redis.ErrNil {
		return ErrorNotExist
	}
	return err
}

func (c *RedisCache) TTL(key string, ttl uint64) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("EXPIRE", key, ttl)
	if err == redis.ErrNil {
		return ErrorNotExist
	}
	return err
}

func RedisProvider(args ...string) (Provider, error) {
	var (
		server   = ":6379"
		password string
	)

	if len(args) > 0 {
		server = args[0]
		if len(args) > 1 {
			password = args[1]
		}
	}

	pool := newPool(server, password)
	conn, err := pool.Dial()
	if err != nil {
		return nil, err
	}
	conn.Close()
	return &RedisCache{pool}, err
}
func RedisProviderDefault() Provider {
	p, err := RedisProvider()
	if err != nil {
		panic(err)
	}
	return p
}

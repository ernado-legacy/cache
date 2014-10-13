package cache

import (
	"reflect"
	"time"
)

type memoryEntry struct {
	value     interface{}
	ttl       uint64
	temporary bool
}

type MemoryCache struct {
	data map[string]*memoryEntry
}

func MemoryProvider() Provider {
	m := new(MemoryCache)
	m.data = make(map[string]*memoryEntry)
	go m.cycle()
	return m
}

func (c *MemoryCache) Set(key string, v interface{}) error {
	c.data[key] = &memoryEntry{value: v}
	return nil
}

func (c *MemoryCache) tick() {
	for k := range c.data {
		if !c.data[k].temporary {
			continue
		}
		c.data[k].ttl -= 1
		if c.data[k].ttl <= 0 {
			c.Remove(k)
		}
	}
}

func (c *MemoryCache) cycle() {
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		c.tick()
	}
}

func (c *MemoryCache) Get(key string, v interface{}) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err = ErrorInvalidType
	}()
	value, ok := c.data[key]
	if !ok {
		return ErrorNotExist
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrorShouldBePointer
	}
	rv.Elem().Set(reflect.ValueOf(value.value))
	return nil
}

func (c *MemoryCache) Remove(key string) error {
	_, ok := c.data[key]
	if !ok {
		return ErrorNotExist
	}
	delete(c.data, key)
	return nil
}

func (c *MemoryCache) TTL(key string, ttl uint64) error {
	_, ok := c.data[key]
	if !ok {
		return ErrorNotExist
	}
	c.data[key].ttl = ttl
	c.data[key].temporary = true
	return nil
}

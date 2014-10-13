package cache

import (
	"errors"
)

var (
	ErrorNotExist        = errors.New("Key not exists")
	ErrorShouldBePointer = errors.New("Should be pointer")
	ErrorInvalidType     = errors.New("Unable to set value: invalid type")
)

type Provider interface {
	Get(key string, v interface{}) error
	Set(key string, v interface{}) error
	TTL(key string, ttl uint64) error
	Remove(key string) error
}

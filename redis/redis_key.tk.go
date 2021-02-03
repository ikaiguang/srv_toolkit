package tkredis

import (
	"github.com/go-kratos/kratos/pkg/conf/env"
	"sync"
)

// key
var (
	_keyPrefix string
	_keyOnce   sync.Once
)

// KeyPrefix .
func KeyPrefix() string {
	_keyOnce.Do(func() {
		_keyPrefix = env.AppID
	})
	return _keyPrefix
}

// Key .
func Key(key string) string {
	return KeyPrefix() + ":" + key
}

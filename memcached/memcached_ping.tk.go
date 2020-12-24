package tkmc

import (
	"context"
	"github.com/go-kratos/kratos/pkg/cache/memcache"
	"github.com/pkg/errors"
)

// Ping .
func Ping() (err error) {
	err = mc.Set(context.Background(), &memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0})
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

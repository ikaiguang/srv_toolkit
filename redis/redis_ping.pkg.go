package tkredis

import (
	"context"
	"github.com/pkg/errors"
)

// Ping .
func Ping() (reply interface{}, err error) {
	reply, err = redisConn.Do(context.Background(), "ping")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

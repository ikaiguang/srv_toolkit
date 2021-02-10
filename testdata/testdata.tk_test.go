package testdata

import (
	tkredis "github.com/ikaiguang/srv_toolkit/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetup(t *testing.T) {
	Setup()

	tkredis.Setup("redis.toml", "Client")
	assert.NotNil(t, tkredis.Redis())
}

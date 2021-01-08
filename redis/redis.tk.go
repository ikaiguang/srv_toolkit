package tkredis

import (
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/env"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
)

// redis
var (
	redisConn *redis.Redis
)

// Redis *redis.Redis
func Redis() *redis.Redis {
	return redisConn
}

// Setup .
func Setup(redisConf, section string) {
	var err error
	defer func() {
		if err != nil {
			tk.Exit(err)
		}
	}()

	redisConn, err = NewRedis(redisConf, section)
	if err != nil {
		return
	}
	//defer redisConn.Close()

	// ping
	if _, err = Ping(); err != nil {
		return
	}
}

// Close .
func Close() (err error) {
	err = redisConn.Close()
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return
}

// NewRedis *redis.Redis
func NewRedis(redisConf, section string) (conn *redis.Redis, err error) {
	cfg, err := getConfig(redisConf, section)
	if err != nil {
		return
	}
	conn = NewRedisWithConfig(cfg)
	//cf = func() {conn.Close()}
	return
}

// NewRedisWithConfig *redis.Redis
func NewRedisWithConfig(cfg *redis.Config) (conn *redis.Redis) {
	conn = redis.NewRedis(cfg)
	return
}

// getConfig .
func getConfig(redisConf, section string) (cfg *redis.Config, err error) {
	var ct paladin.TOML
	if err = paladin.Get(redisConf).Unmarshal(&ct); err != nil {
		err = errors.WithStack(err)
		return
	}

	cfg = &redis.Config{}
	if err = ct.Get(section).UnmarshalTOML(cfg); err != nil {
		err = errors.WithStack(err)
		return
	}
	if env.AppID != "" {
		cfg.Name = env.AppID
	}
	return
}

package tkmc

import (
	"github.com/go-kratos/kratos/pkg/cache/memcache"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
)

// memcached
var (
	mc *memcache.Memcache
)

// MC *memcache.Memcache
func MC() *memcache.Memcache {
	return mc
}

// Setup .
func Setup(mcConf, section string) {
	var err error
	defer func() {
		if err != nil {
			tk.Fatal(err)
		}
	}()

	mc, err = NewMC(mcConf, section)
	if err != nil {
		return
	}
	//defer mc.Close()

	// ping
	if err = Ping(); err != nil {
		return
	}
}

// Close .
func Close() (err error) {
	err = mc.Close()
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return
}

// NewMC *memcache.Memcache
func NewMC(mcConf, section string) (conn *memcache.Memcache, err error) {
	cfg, err := getConfig(mcConf, section)
	if err != nil {
		return
	}
	conn = NewMCWithConfig(cfg)
	//cf = func() {conn.Close()}
	return
}

// NewMCWithConfig *memcache.Memcache
func NewMCWithConfig(cfg *memcache.Config) (conn *memcache.Memcache) {
	conn = memcache.New(cfg)
	return
}

// getConfig .
func getConfig(mcConf, section string) (cfg *memcache.Config, err error) {
	var ct paladin.TOML
	if err = paladin.Get(mcConf).Unmarshal(&ct); err != nil {
		err = errors.WithStack(err)
		return
	}

	cfg = &memcache.Config{}
	if err = ct.Get(section).UnmarshalTOML(cfg); err != nil {
		err = errors.WithStack(err)
		return
	}
	//if env.AppID != "" {
	//	cfg.Name = env.AppID
	//}
	return
}

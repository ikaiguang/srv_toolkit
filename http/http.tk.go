package tkhttp

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	tklog "github.com/ikaiguang/srv_toolkit/log"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"time"
)

// http
var (
	conf blademaster.ServerConfig
)

// initConfig .
func initConfig(httpConf, section string) (err error) {
	var ct paladin.TOML
	if err = paladin.Get(httpConf).Unmarshal(&ct); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = ct.Get(section).UnmarshalTOML(&conf); err != nil {
		err = errors.WithStack(err)
		return
	}
	if conf.Network == "" {
		conf.Network = "tcp"
	}
	return
}

// New .
func New(httpConf, section string) (engine *blademaster.Engine, err error) {
	err = initConfig(httpConf, section)
	if err != nil {
		return
	}

	//engine = blademaster.DefaultServer(&conf)
	engine = blademaster.NewServer(&conf)
	engine.Use(blademaster.Recovery(), blademaster.Trace(), Logger())
	return
}

// Start run server
func Start(engine *blademaster.Engine) (err error) {
	//err = engine.Start()
	l, err := net.Listen(conf.Network, conf.Addr)
	if err != nil {
		err = errors.Wrapf(err, "http listen tcp: %s", conf.Addr)
		return
	}

	tklog.Info(fmt.Sprintf("http listen addr: %s", l.Addr().String()))
	server := &http.Server{
		ReadTimeout:  time.Duration(conf.ReadTimeout),
		WriteTimeout: time.Duration(conf.WriteTimeout),
	}
	go func() {
		if err := engine.RunServer(server, l); err != nil {
			if errors.Cause(err) == http.ErrServerClosed {
				tklog.Info("http server closed")
				return
			}
			tk.Fatal(errors.Wrapf(err, "http engine.ListenServer(%+v, %+v)", server, l))
		}
	}()
	return
}

// Close .
func Close(engine *blademaster.Engine) {
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	if err := engine.Shutdown(ctx); err != nil {
		tklog.Error(fmt.Sprintf("httpSrv.Shutdown error(%v)", err))
	}
	cancel()
}

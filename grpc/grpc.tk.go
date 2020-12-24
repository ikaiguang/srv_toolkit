package tkgrpc

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	tklog "github.com/ikaiguang/srv_toolkit/log"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

// http
var (
	conf warden.ServerConfig
)

// initConfig .
func initConfig(rpcConf, section string) (err error) {
	var ct paladin.TOML
	if err = paladin.Get(rpcConf).Unmarshal(&ct); err != nil {
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

// New new a grpc server.
func New(rpcConf, section string) (ws *warden.Server, err error) {
	err = initConfig(rpcConf, section)
	if err != nil {
		return
	}
	ws = warden.NewServer(&conf)
	return
}

// Start create a new goroutine run server with configured listen addr
// will panic if any error happend
// return server itself
func Start(ws *warden.Server) (err error) {
	//ws, err = ws.Start()
	lis, err := net.Listen(conf.Network, conf.Addr)
	if err != nil {
		return errors.WithStack(err)
	}
	tklog.Info(fmt.Sprintf("grpc listen addr: %v", lis.Addr()))
	reflection.Register(ws.Server())
	go func() {
		if err := ws.Serve(lis); err != nil {
			tk.Fatal(err)
		}
	}()
	return
}

// Close .
func Close(ws *warden.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	if err := ws.Shutdown(ctx); err != nil {
		tklog.Error(fmt.Sprintf("grpcSrv.Shutdown error(%v)", err))
	}
	cancel()
}

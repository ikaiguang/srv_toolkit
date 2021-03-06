package setup

import (
	"fmt"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	tkapp "github.com/ikaiguang/srv_toolkit/app"
	tkdb "github.com/ikaiguang/srv_toolkit/db"
	tkflag "github.com/ikaiguang/srv_toolkit/flag"
	tkgrpc "github.com/ikaiguang/srv_toolkit/grpc"
	tkhttp "github.com/ikaiguang/srv_toolkit/http"
	tkinit "github.com/ikaiguang/srv_toolkit/initialize"
	tkjwt "github.com/ikaiguang/srv_toolkit/jwt"
	tklog "github.com/ikaiguang/srv_toolkit/log"
	tkredis "github.com/ikaiguang/srv_toolkit/redis"
	tkru "github.com/ikaiguang/srv_toolkit/redis/utils"
	dbmr "github.com/ikaiguang/srv_toolkit/srv_hello/internal/migration"
	models "github.com/ikaiguang/srv_toolkit/srv_hello/internal/model"
	routes "github.com/ikaiguang/srv_toolkit/srv_hello/internal/route"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"os"
	"os/signal"
	"syscall"
)

// Serve .
func Serve() {
	var err error
	defer func() {
		if err != nil {
			tk.Fatal(err)
		}
	}()

	Setup()
	defer Close()

	// http server
	engine, err := StartHTTP()
	if err != nil {
		return
	}
	defer tkhttp.Close(engine)
	// 注册路由
	routes.RegisterHTTP(engine)

	// grpc
	ws, err := StartGRPC()
	if err != nil {
		return
	}
	defer tkgrpc.Close(ws)
	// 注册路由
	routes.RegisterGRPC(ws)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		tklog.Info(fmt.Sprintf("get a signal %s", s.String()))
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			tklog.Info("server exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

// Setup .
func Setup() {
	// 初始化
	// 初始化参数
	tkflag.Setup()
	// 初始化配置
	tkinit.Setup("application.toml", "App")
	if tkinit.IsTest() {
		tkapp.OmitDetail()
	}

	// 程序
	// 初始化日志
	tklog.Setup("log.toml", "Log")
	tkapp.SetLogger(&tklog.AppLogger{})
	tkjwt.SetLogger(&tklog.AppLogger{})

	// 初始化redis
	tkredis.Setup("redis.toml", "Client")
	tkru.Init(tkredis.Redis())
	// 初始化memcached
	//tkmc.Setup("memcached.toml", "Client")

	// 初始化数据库迁移 & 创建数据库
	dbmr.Setup("db.toml", "Migration")
	createDB()

	// 初始化数据库
	tkdb.Setup("db.toml", "Client")
	models.Setup()

	// 数据迁移
	migrationDB()
}

// createDB 数据库 运行迁移，请初始化数据库迁移连接
func createDB() {
	var err error

	// 数据库
	err = dbmr.CreateDB()
	if err != nil {
		tk.Fatal(err)
	}
}

// migrationDB 数据库 运行迁移，请初始化数据库迁移连接
func migrationDB() {
	var err error

	// 数据库迁移
	err = dbmr.Migrations()
	if err != nil {
		tk.Fatal(err)
	}
}

// StartHTTP http
func StartHTTP() (engine *blademaster.Engine, err error) {
	engine, err = tkhttp.New("http.toml", "Server")
	if err != nil {
		return
	}
	err = tkhttp.Start(engine)
	if err != nil {
		return
	}
	return
}

// StartGRPC grpc
func StartGRPC() (ws *warden.Server, err error) {
	ws, err = tkgrpc.New("grpc.toml", "Server")
	if err != nil {
		return
	}
	err = tkgrpc.Start(ws)
	if err != nil {
		return
	}
	return
}

// Close .
func Close() {
	var err error

	err = tkinit.Close()
	if err != nil {
		tklog.ERROR(err)
	}

	//err = tkmc.Close()
	//if err != nil {
	//	tklog.ERROR(err)
	//}

	err = tkredis.Close()
	if err != nil {
		tklog.ERROR(err)
	}

	err = tkdb.Close()
	if err != nil {
		tklog.ERROR(err)
	}

	err = tklog.Close()
	if err != nil {
		tklog.ERROR(err)
	}
}

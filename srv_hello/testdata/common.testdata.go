package testdata

import (
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	tkinit "github.com/ikaiguang/srv_toolkit/initialize"
	"github.com/ikaiguang/srv_toolkit/srv_hello/internal/setup"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

// var
var (
	onceServer  sync.Once
	pwdPath     string
	toolkitPath string
	confGRPC    warden.ServerConfig
	confHTTP    blademaster.ServerConfig
	addrZero    = "0.0.0.0"
	addrLocal   = "127.0.0.1"
)

func init() {
	var err error

	pwdPath, err = os.Getwd()
	if err != nil {
		tk.Fatal(err)
		return
	}

	toolkitPath, err = filepath.Abs(filepath.Join(".", "../.."))
	if err != nil {
		tk.Fatal(err)
		return
	}
}

// lazyServe .
func lazyServe() {
	onceServer.Do(func() {
		copyConf()
		go setup.Serve()
		time.Sleep(time.Second)
		initConfig()
	})
}

func initConfig() {
	var err error

	// paladin 先启动服务，无需再次监听配置目录
	//paladin.DefaultClient, err = paladin.NewFile("")
	//if err != nil {
	//	err = errors.WithStack(err)
	//	return
	//}

	if err = initConfigGRPC("grpc.toml", "Server"); err != nil {
		tk.Fatal(err)
		return
	}

	if err = initConfigHTTP("http.toml", "Server"); err != nil {
		tk.Fatal(err)
		return
	}
}

func initConfigGRPC(grpcConf, section string) (err error) {
	var ct paladin.TOML
	if err = paladin.Get(grpcConf).Unmarshal(&ct); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = ct.Get(section).UnmarshalTOML(&confGRPC); err != nil {
		err = errors.WithStack(err)
		return
	}
	if confGRPC.Network == "" {
		confGRPC.Network = "tcp"
	}
	host, port, err := net.SplitHostPort(confGRPC.Addr)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if host == addrZero {
		confGRPC.Addr = net.JoinHostPort(addrLocal, port)
	}
	return
}

func initConfigHTTP(httpConf, section string) (err error) {
	var ct paladin.TOML
	if err = paladin.Get(httpConf).Unmarshal(&ct); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = ct.Get(section).UnmarshalTOML(&confHTTP); err != nil {
		err = errors.WithStack(err)
		return
	}
	if confHTTP.Network == "" {
		confGRPC.Network = "tcp"
	}
	host, port, err := net.SplitHostPort(confHTTP.Addr)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if host == addrZero {
		confHTTP.Addr = net.JoinHostPort(addrLocal, port)
	}
	return
}

func copyConf() {
	if pwdPath == toolkitPath {
		return
	}

	src := filepath.Join(toolkitPath, tkinit.RelPathConfig)
	dst := filepath.Join(pwdPath, tkinit.RelPathConfig)

	if err := copyDir(src, dst); err != nil {
		tk.Fatal(err)
		return
	}
}

// copyDir copies a whole directory recursively
func copyDir(src string, dst string) (err error) {
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	// 原文件夹信息
	if srcInfo, err = os.Stat(src); err != nil {
		err = errors.WithStack(err)
		return
	}

	// 创建目标信息
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		err = errors.WithStack(err)
		return
	}

	// 读取原文件所有信息
	if fds, err = ioutil.ReadDir(src); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, fd := range fds {
		srcFp := path.Join(src, fd.Name())
		dstFp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcFp, dstFp); err != nil {
				err = errors.WithStack(err)
				return err
			}
		} else {
			if err = copyFile(srcFp, dstFp); err != nil {
				err = errors.WithStack(err)
				return err
			}
		}
	}
	return
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) (err error) {
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	// 读源文件
	if srcFd, err = os.Open(src); err != nil {
		err = errors.WithStack(err)
		return err
	}
	defer srcFd.Close()

	// 写目标文件
	if dstFd, err = os.Create(dst); err != nil {
		err = errors.WithStack(err)
		return err
	}
	defer dstFd.Close()

	// 复制
	if _, err = io.Copy(dstFd, srcFd); err != nil {
		err = errors.WithStack(err)
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		err = errors.WithStack(err)
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

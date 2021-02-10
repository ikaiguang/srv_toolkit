package testdata

import (
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	tkflag "github.com/ikaiguang/srv_toolkit/flag"
	tkinit "github.com/ikaiguang/srv_toolkit/initialize"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// var
var (
	pwdPath     string
	toolkitPath string
	configPath  string
)

func init() {
	toolkitPath = filepath.Join(os.Getenv("GOPATH"), "src/github.com/ikaiguang/srv_toolkit")
	pwdPath = filepath.Join(toolkitPath, "testdata")
	configPath = filepath.Join(pwdPath, "conf.d")

	copyConf()
}

// Setup .
func Setup() {
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	// 初始化
	// 初始化参数
	tkflag.Setup()
	// paladin
	paladin.DefaultClient, err = paladin.NewFile(configPath)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
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

package tkinit

import (
	"github.com/go-kratos/kratos/pkg/conf/env"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// path
var (
	appPath         string
	configPath      string
	runtimePath     string
	logPath         string
	staticPath      string
	attachmentsPath string
	version         string
)

// config
const (
	defaultFileMode os.FileMode = 0755
)

// Config .
type Config struct {
	AppID   string
	Version string
	// dev/fat1/uat/pre/prod
	DeployEnv string
	//Region         string
	//Zone           string
	//Hostname       string
	//Color          string
	//DiscoveryNodes string
}

// Setup .
func Setup(appConf, section string) {
	var err error
	defer func() {
		if err != nil {
			tk.Exit(err)
		}
	}()

	// path
	err = InitPath()
	if err != nil {
		return
	}

	// paladin
	paladin.DefaultClient, err = paladin.NewFile(ConfigPath())
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// app
	cfg, err := getConfig(appConf, section)
	if err != nil {
		return
	}
	SetEnv(cfg)
}

// Close .
func Close() (err error) {
	//return paladin.Close()
	paladin.DefaultClient = nil
	return
}

// ConfigPath config path
func ConfigPath() string {
	return configPath
}

// AppPath app path
func AppPath() string {
	return appPath
}

// Version .
func Version() string {
	return version
}

// RuntimePath runtime path
func RuntimePath() string {
	return runtimePath
}

// LogPath log path
func LogPath() string {
	return logPath
}

// StaticPath static path
func StaticPath() string {
	return staticPath
}

// AttachmentPath attachments path
func AttachmentPath() string {
	return attachmentsPath
}

// SetEnv .
func SetEnv(cfg *Config) {
	version = cfg.Version

	if len(env.AppID) == 0 {
		env.AppID = cfg.AppID
	}
	if len(env.DeployEnv) == 0 {
		env.DeployEnv = cfg.DeployEnv
	}
}

// IsTest test
func IsTest() bool {
	switch env.DeployEnv {
	case env.DeployEnvPre, env.DeployEnvProd:
		return false
	default:
		return true
	}
}

// InitPath .
func InitPath() (err error) {
	// app
	appPath, err = os.Getwd()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// 目录信息
	var dirInfo os.FileInfo

	// config
	configPath = filepath.Join(appPath, "config")
	dirInfo, err = os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = errors.Errorf("cannot find config path : ./config;\n\terror : %s", err.Error())
			return
		}
		err = errors.WithStack(err)
		return
	}
	if !dirInfo.IsDir() {
		err = errors.Errorf("Not a directory : %s", configPath)
		return
	}

	// static
	staticPath = filepath.Join(appPath, "static")
	dirInfo, err = os.Stat(staticPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(staticPath, defaultFileMode)
			if err != nil {
				err = errors.WithStack(err)
				return
			}
		}
		err = errors.WithStack(err)
		return
	}
	if !dirInfo.IsDir() {
		err = errors.Errorf("Not a directory : %s", staticPath)
		return
	}

	// runtime
	runtimePath = filepath.Join(appPath, "runtime")
	dirInfo, err = os.Stat(runtimePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(runtimePath, defaultFileMode)
			if err != nil {
				err = errors.WithStack(err)
				return
			}
		}
		err = errors.WithStack(err)
		return
	}
	if !dirInfo.IsDir() {
		err = errors.Errorf("Not a directory : %s", runtimePath)
		return
	}

	// log
	logPath = filepath.Join(runtimePath, "logs")
	dirInfo, err = os.Stat(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(logPath, defaultFileMode)
			if err != nil {
				err = errors.WithStack(err)
				return
			}
		}
		err = errors.WithStack(err)
		return
	}
	if !dirInfo.IsDir() {
		err = errors.Errorf("Not a directory : %s", logPath)
		return
	}

	// attachments
	attachmentsPath = filepath.Join(appPath, "attachments")
	dirInfo, err = os.Stat(attachmentsPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(attachmentsPath, defaultFileMode)
			if err != nil {
				err = errors.WithStack(err)
				return
			}
		}
		err = errors.WithStack(err)
		return
	}
	if !dirInfo.IsDir() {
		err = errors.Errorf("Not a directory : %s", attachmentsPath)
		return
	}
	return
}

// getConfig .
func getConfig(filename, section string) (cfg *Config, err error) {
	var ct paladin.TOML
	if err = paladin.Get(filename).Unmarshal(&ct); err != nil {
		err = errors.WithStack(err)
		return
	}

	cfg = &Config{}
	if err = ct.Get(section).UnmarshalTOML(cfg); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

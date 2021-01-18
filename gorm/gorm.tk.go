package tkgorm

import (
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// db
var (
	dbConn *gorm.DB
)

// DB *sql.DB
func DB() *gorm.DB {
	return dbConn
}

// Config .
type Config struct {
	gorm.Config
	Dsn string
}

// Setup .
func Setup(gormConf, section string) {
	var err error
	defer func() {
		if err != nil {
			tk.Fatal(err)
		}
	}()

	cfg, err := getConfig(gormConf, section)
	if err != nil {
		return
	}
	dbConn, err = gorm.Open(mysql.Open(cfg.Dsn), &cfg.Config)
	if err != nil {
		return
	}
}

// Close .
func Close() (err error) {
	return
}

// getConfig .
func getConfig(gormConf, section string) (cfg *Config, err error) {
	var ct paladin.TOML
	if err = paladin.Get(gormConf).Unmarshal(&ct); err != nil {
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

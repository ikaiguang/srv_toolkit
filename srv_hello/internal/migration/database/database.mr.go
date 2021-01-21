package dbmrdb

import (
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	tkdb "github.com/ikaiguang/srv_toolkit/db"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

// migration
var (
	dbConn *gorm.DB
	dbConf *Config
)

// Config .
type Config struct {
	tkdb.Config
	CreateDb     bool
	DatabaseName string
}

// DB return dbConn
func DB() *gorm.DB {
	return dbConn
}

// DBName .
func DBName() string {
	return dbConf.DatabaseName
}

// Setup .
func Setup(dbConfFile, section string) {
	var err error
	defer func() {
		if err != nil {
			tk.Fatal(err)
		}
	}()

	// 配置
	dbConf, err = getConfig(dbConfFile, section)
	if err != nil {
		return
	}
	dbConn, err = newDB(dbConf)
	if err != nil {
		return
	}
}

// Close .
func Close() (err error) {
	// kratos
	// gorm 无Close方法
	//err = dbConn.Close()
	//if err != nil {
	//	err = errors.WithStack(err)
	//	return err
	//}
	dbConn = nil
	return
}

// RebuildDBConn .
func RebuildDBConn() (err error) {
	dbConn, err = rebuildDBConn(dbConf)
	if err != nil {
		return
	}
	return
}

// rebuildDBConn .
func rebuildDBConn(cfg *Config) (db *gorm.DB, err error) {
	if dbConf == nil {
		err = errors.New("dbConf is nil")
		return
	}
	// 重建DSN
	conf, err := mysqlDriver.ParseDSN(cfg.DSN)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	conf.DBName = cfg.DatabaseName
	cfg.DSN = conf.FormatDSN()

	return newDB(cfg)
}

// newDB .
func newDB(cfg *Config) (db *gorm.DB, err error) {
	// 连接
	dbLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Info,
		Colorful:      true,
	})
	//dbLogger.LogMode(logger.Info)
	dbConf := &gorm.Config{
		Logger: dbLogger,
	}
	db, err = gorm.Open(mysql.Open(cfg.DSN), dbConf)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// ping
	// gorm 打开成功即可用连接，无需ping测试
	//err = dbConn.Ping(context.Background())
	//if err != nil {
	//	err = errors.WithStack(err)
	//	return
	//}
	return
}

// getConfig .
func getConfig(dbConfFile, section string) (cfg *Config, err error) {
	var ct paladin.TOML
	if err = paladin.Get(dbConfFile).Unmarshal(&ct); err != nil {
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

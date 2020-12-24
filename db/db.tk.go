package dbpkg

import (
	"context"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
)

// db
var (
	dbConn *sql.DB
)

// DB *sql.DB
func DB() *sql.DB {
	return dbConn
}

// Setup .
func Setup(dbConf, section string) {
	var err error
	defer func() {
		if err != nil {
			tk.Exit(err)
		}
	}()

	dbConn, err = NewDB(dbConf, section)
	if err != nil {
		return
	}
	//defer dbConn.Close()

	// ping
	err = dbConn.Ping(context.Background())
	if err != nil {
		err = errors.WithStack(err)
		return
	}
}

// Close .
func Close() (err error) {
	err = dbConn.Close()
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return
}

// NewDB *sql.DB
func NewDB(dbConf, section string) (conn *sql.DB, err error) {
	cfg, err := getConfig(dbConf, section)
	if err != nil {
		return
	}
	conn = NewDBWithConfig(cfg)
	//cf = func() {conn.Close()}
	return
}

// NewDBWithConfig *sql.DB
func NewDBWithConfig(cfg *sql.Config) (conn *sql.DB) {
	conn = sql.NewMySQL(cfg)
	return
}

// getConfig .
func getConfig(dbConf, section string) (cfg *sql.Config, err error) {
	var ct paladin.TOML
	if err = paladin.Get(dbConf).Unmarshal(&ct); err != nil {
		err = errors.WithStack(err)
		return
	}

	cfg = &sql.Config{}
	if err = ct.Get(section).UnmarshalTOML(cfg); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// Begin transaction
func Begin(ctx context.Context) (tx *sql.Tx, err error) {
	return dbConn.Begin(ctx)
}

// Commit .
func Commit(tx *sql.Tx) error {
	return tx.Commit()
}

// Rollback .
func Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

package dbmrdb

import (
	"fmt"
	tklog "github.com/ikaiguang/srv_toolkit/log"
	"github.com/pkg/errors"
)

// DBInfo .
type DBInfo struct {
	SchemaName      string `gorm:"column:schema_name"`
	SchemaNameUpper string `gorm:"column:SCHEMA_NAME"`
}

// CreateDB 创建数据库
func CreateDB() (err error) {
	// 需要创建数据库？
	if !dbConf.CreateDb {
		return
	}
	// 检查数据库
	dbName := DBName()
	tklog.Info("数据库迁移：检查数据库：" + dbName)
	var infos []*DBInfo
	queryDBSQL := "SELECT schema_name FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?"
	err = DB().Raw(queryDBSQL, dbName).Scan(&infos).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(infos) > 0 {
		return
	}

	// 事务
	tx := DB().Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		}
	}()

	// 创建数据库
	tklog.Info("数据库迁移：创建数据库：" + dbName)
	execSQL := `CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARSET utf8mb4;`
	err = tx.Exec(fmt.Sprintf(execSQL, dbName)).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// 分配权限
	//

	// 提交修改
	err = tx.Commit().Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

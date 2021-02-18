package dbmrv1p0p1_test

import (
	tklog "github.com/ikaiguang/srv_toolkit/log"
	dbmrdb "github.com/ikaiguang/srv_toolkit/srv_hello/internal/migration/database"
	testmodels "github.com/ikaiguang/srv_toolkit/srv_hello/internal/model/test"
	"github.com/pkg/errors"
)

// Test .
type Test struct{}

// CreateTable
func (t *Test) CreateTable() (err error) {
	// table name
	tableName := testmodels.TestModel.TableName()
	tklog.Info("数据库迁移：检查表：" + tableName)
	if dbmrdb.DB().Migrator().HasTable(tableName) {
		return
	}

	// 事务
	tx := dbmrdb.DB().Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		}
	}()

	// 创建数据库
	tklog.Info("数据库迁移：创建表：" + tableName)
	err = tx.Exec(testmodels.TestModel.TableSQL()).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// 提交修改
	err = tx.Commit().Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

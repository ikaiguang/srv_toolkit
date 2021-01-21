package dbmrv1p0p1_test

import (
	"fmt"
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

	execSQL := `
CREATE TABLE %s
(
    id            BIGINT AUTO_INCREMENT COMMENT 'id',
    column_string VARCHAR(255) NOT NULL DEFAULT '' COMMENT '字符串',
    column_int    INT          NOT NULL DEFAULT '0' COMMENT '数字',
    PRIMARY KEY (id)
) ENGINE InnoDB
  DEFAULT CHARSET utf8mb4
    COMMENT '测试表';
`

	// 事务
	tx := dbmrdb.DB().Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		}
	}()

	// 创建数据库
	tklog.Info("数据库迁移：创建表：" + tableName)
	err = tx.Exec(fmt.Sprintf(execSQL, tableName)).Error
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

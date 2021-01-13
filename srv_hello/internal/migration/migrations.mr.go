package dbmr

import (
	dbmrdb "github.com/ikaiguang/srv_toolkit/srv_hello/internal/migration/database"
	dbmrv1p0p1 "github.com/ikaiguang/srv_toolkit/srv_hello/internal/migration/v1.0.1"
)

// Migrations 运行迁移，请初始化下面的配置
// flagpkg.Setup()   // 初始化参数
// configpkg.Setup() // 初始化配置
func Migrations() (err error) {
	// 创建数据库
	err = dbmrdb.CreateDB()
	if err != nil {
		return
	}

	// 重建数据库连接：切换到数据库
	err = dbmrdb.RebuildDBConn()
	if err != nil {
		return
	}
	defer func() { _ = dbmrdb.Close() }()

	// 迁移 v1.0.1
	err = dbmrv1p0p1.RunMigrations()
	if err != nil {
		return
	}
	return
}

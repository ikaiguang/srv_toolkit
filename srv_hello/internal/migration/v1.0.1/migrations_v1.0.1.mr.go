package dbmrv1p0p1

import dbmrv1p0p1_test "github.com/ikaiguang/srv_toolkit/srv_hello/internal/migration/v1.0.1/test"

// RunMigrations .
func RunMigrations() (err error) {
	// test
	err = dbmrv1p0p1_test.Migrations()
	if err != nil {
		return
	}
	return
}

package dbmrv1p0p1_test

// Migrations .
func Migrations() (err error) {
	// test
	testHandler := &Test{}
	err = testHandler.CreateTable()
	if err != nil {
		return
	}
	return
}

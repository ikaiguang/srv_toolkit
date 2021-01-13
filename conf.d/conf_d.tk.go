package tkconfd

import (
	"path/filepath"
	"runtime"
)

// var
var (
	currentPath string
)

func init() {
	_, file, _, _ := runtime.Caller(0)

	currentPath = filepath.Dir(file)
}

// Path return current path
func Path() string {
	return currentPath
}

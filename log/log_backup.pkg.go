package tklog

import "github.com/go-kratos/kratos/pkg/log"

// SetupBackup .
func Setup2() {
	cfg := &log.Config{
		Stdout: true,
	}
	log.Init(cfg)
	//defer log.Close()
}

// CloseBackup .
func CloseBackup() (err error) {
	return log.Close()
}

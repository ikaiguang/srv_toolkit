package tklog

import "github.com/go-kratos/kratos/pkg/log"

// SetupBackup .
func SetupBackup(logConf, section string) {
	cfg, err := getConfig(logConf, section)
	if err != nil {
		return
	}
	log.Init(&cfg.Config)
	//log.SetFormat("[%D %T] [%L] [%S] %M")
	//defer log.Close()
}

// CloseBackup .
func CloseBackup() (err error) {
	return log.Close()
}

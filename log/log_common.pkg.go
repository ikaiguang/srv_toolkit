package tklog

import (
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/pkg/errors"
)

// Config .
type Config struct {
	log.Config
	StdOutLog      bool
	LogFilename    string
	LogFormat      string
	LogShortCaller bool
}

// log
const (
	// stack trace
	_stackTrace = "stack_trace"
	// stack trace depth
	_stackTracerDepth = 6
	// log format
	_logFormatJSON = "json"
	_logFormatText = "text"
)

// getConfig .
func getConfig(filename, section string) (cfg *Config, err error) {
	var ct paladin.TOML
	if err = paladin.Get(filename).Unmarshal(&ct); err != nil {
		err = errors.WithStack(err)
		return
	}

	cfg = &Config{}
	if err = ct.Get(section).UnmarshalTOML(cfg); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

package tklog

import (
	tkinit "github.com/ikaiguang/srv_toolkit/initialize"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"time"
)

// log
const (
	_addCallerSkip = 1
)

// log
var (
	logger *zap.Logger
)

// Setup .
func Setup(logConf, section string) {
	var err error
	defer func() {
		if err != nil {
			tk.Exit(err)
		}
	}()

	//err = initDevelopment()
	//err = initConsole()
	err = initProduction(logConf, section)
	if err != nil {
		return
	}
	logger = logger.WithOptions(zap.AddCallerSkip(_addCallerSkip))
}

// Close .
func Close() (err error) {
	_ = logger.Sync()
	//err = logger.Sync()
	//if err != nil {
	//	err = errors.WithStack(err)
	//	return
	//}
	err = initConsole()
	if err != nil {
		return
	}
	return
}

// initProduction .
func initProduction(logConf, section string) (err error) {
	cfg, err := getConfig(logConf, section)
	if err != nil {
		return
	}
	var cores []zapcore.Core

	// log
	logEncoderCfg := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "tracer",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.999"))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	if cfg.LogShortCaller {
		logEncoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	}
	logEncoder := zapcore.NewConsoleEncoder(logEncoderCfg)
	if cfg.LogFormat == _logFormatJSON {
		logEncoder = zapcore.NewJSONEncoder(logEncoderCfg)
	}
	logLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	logWriter, err := getLogWriter(cfg.LogFilename)
	if err != nil {
		return
	}

	// log core
	cores = append(cores, zapcore.NewCore(logEncoder, zapcore.AddSync(logWriter), logLevel))

	// std
	if cfg.StdOutLog {
		stdEncoderCfg := logEncoderCfg
		stdEncoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		stdEncoderCfg.EncodeCaller = zapcore.FullCallerEncoder
		stdEncoder := zapcore.NewConsoleEncoder(stdEncoderCfg)
		stdWriter := os.Stderr
		cores = append(cores, zapcore.NewCore(stdEncoder, zapcore.AddSync(stdWriter), logLevel))
	}

	// logger
	core := zapcore.NewTee(cores...)
	logger = zap.New(core, zap.AddCaller())

	return
}

// getLogWriter log writer
func getLogWriter(filename string) (writer io.Writer, err error) {
	writer, err = rotatelogs.New(
		filepath.Join(tkinit.LogPath(), filename+"app.log_%Y%m%d.log"),
		rotatelogs.WithRotationTime(time.Hour*24),  // one day
		rotatelogs.WithMaxAge(time.Hour*24*365*10), // ten years
	)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// initConsole .
func initConsole() (err error) {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.999"))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	cfg := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err = cfg.Build()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// initDevelopment .
func initDevelopment() (err error) {
	//logger, err = zap.NewProduction()
	logger, err = zap.NewDevelopment()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
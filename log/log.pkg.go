package tklog

import (
	"github.com/go-kratos/kratos/pkg/log"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
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
			tk.Fatal(err)
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
	_ = log.Close()
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

	// 初始化日志
	log.Init(&cfg.Config)
	if len(cfg.LogFormat) > 0 {
		log.SetFormat(cfg.LogFormat)
	}

	// logger
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
	if cfg.LogFormatter == _logFormatJSON {
		logEncoder = zapcore.NewJSONEncoder(logEncoderCfg)
	}
	logLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= convertZapLv(cfg.V)
	})
	logWriter, err := getLogWriter(cfg)
	if err != nil {
		return
	}

	// log core
	cores = append(cores, zapcore.NewCore(logEncoder, zapcore.AddSync(logWriter), logLevel))

	// std
	if cfg.Stdout {
		stdEncoderCfg := logEncoderCfg
		stdEncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		stdEncoderCfg.EncodeCaller = zapcore.FullCallerEncoder
		stdEncoder := zapcore.NewConsoleEncoder(stdEncoderCfg)
		stdWriter := colorable.NewColorableStderr()
		cores = append(cores, zapcore.NewCore(stdEncoder, zapcore.AddSync(stdWriter), logLevel))
	}

	// logger
	core := zapcore.NewTee(cores...)
	logger = zap.New(core, zap.AddCaller())

	return
}

// convertZapLv .
func convertZapLv(logLv int32) zapcore.Level {
	switch logLv {
	case 0:
		return zapcore.DebugLevel
	case 1:
		return zapcore.InfoLevel
	case 2:
		return zapcore.WarnLevel
	case 3:
		return zapcore.ErrorLevel
	case 4:
		return zapcore.FatalLevel
	default:
		return zapcore.DPanicLevel
	}
}

// getLogWriter log writer
func getLogWriter(cfg *Config) (writer io.Writer, err error) {
	opts := []rotatelogs.Option{
		rotatelogs.WithRotationTime(time.Hour * 24), // one day
	}
	// 文件大小
	if cfg.RotateSize > 0 {
		opts = append(opts, rotatelogs.WithRotationSize(cfg.RotateSize))
	}
	// 存储 n 个 或 n 久
	if cfg.MaxLogFile > 0 {
		opts = append(opts, rotatelogs.WithRotationCount(uint(cfg.MaxLogFile))) // n 个
	} else {
		opts = append(opts, rotatelogs.WithMaxAge(time.Hour*24*365*10)) // ten years
	}
	writer, err = rotatelogs.New(
		filepath.Join(cfg.Dir, cfg.LogFilename+"_app.%Y%m%d.log"),
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

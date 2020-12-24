package tklog

import (
	"context"
	"fmt"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ERROR .
func ERROR(err error) {
	//logger.Error(err.Error())
	logger.Error(err.Error(), zap.String(_stackTrace, ErrStackTrace(err)))
}

// ERRORC .
func ERRORC(ctx context.Context, err error) {
	args := append(AddExtraField(ctx), zap.String(_stackTrace, ErrStackTrace(err)))
	logger.Error(err.Error(), args...)
}

// INFO .
func INFO(err error) {
	//logger.Info(err.Error())
	logger.Info(err.Error(), zap.String(_stackTrace, ErrStackTrace(err)))
}

// stackTracer errors.StackTrace
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// ErrStackTrace error stack trace
func ErrStackTrace(err error) (est string) {
	if err == nil {
		file, line := tk.File(3)
		est = fmt.Sprintf("\n%s:%d\n", file, line)
		return
	}

	trace, ok := err.(stackTracer)
	if !ok {
		file, line := tk.File(3)
		est = fmt.Sprintf("\n%s:%d\n", file, line)
		return
	}
	st := trace.StackTrace()

	if len(st) > _stackTracerDepth {
		est = fmt.Sprintf("%+v", st[0:_stackTracerDepth])
		return
	}
	est = fmt.Sprintf("%+v", st)
	return
}

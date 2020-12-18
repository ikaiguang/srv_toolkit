package tke

import (
	"fmt"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
)

// error
const (
	_stackTracerDepth = 7
)

// Status .
type Status struct {
	Code   Code
	Msg    string
	Detail string
}

// Error .
func (s *Status) Error() string {
	//return fmt.Sprintf(`{"code":%d,"msg":"%s","detail":"%s"}`, s.Code, s.Msg, s.Detail)
	return fmt.Sprintf(`code:%d,msg:%s,detail:%s`, s.Code, s.Msg, s.Detail)
}

// New .
func New(code Code) error {
	return errors.WithStack(&Status{Code: code, Msg: Msg(code)})
}

// Newf .
func Newf(code Code, err error) error {
	return errors.WithStack(&Status{Code: code, Msg: Msg(code), Detail: err.Error()})
}

// FromError 解析错误
func FromError(err error) (s *Status, ok bool) {
	err = errors.Cause(err)
	s, ok = err.(*Status)
	return
}

// stackTracer errors.StackTrace
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Stacktrace error stack trace
func Stacktrace(err error, stackTracerDepths ...int) (est string) {
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

	depth := _stackTracerDepth
	if len(stackTracerDepths) > 0 {
		depth = stackTracerDepths[0]
	}
	if len(st) > depth {
		est = fmt.Sprintf("%+v", st[0:depth])
		return
	}
	est = fmt.Sprintf("%+v", st)
	return
}

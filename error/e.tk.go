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
	code    Code
	msg     string
	details []interface{}
}

// Error .
func (s *Status) Error() string {
	return fmt.Sprintf(`code:%d | msg:%s | details:%v`, s.Code(), s.Message(), s.Details())
}

// Code return error code
func (s *Status) Code() int { return int(s.code.Code()) }

// Code32 return error code
func (s *Status) Code32() int32 { return s.code.Code() }

// Message return error message
func (s *Status) Message() string {
	return s.msg
}

// Details return details.
func (s *Status) Details() []interface{} { return s.details }

// New .
func New(code Code, details ...interface{}) error {
	return errors.WithStack(&Status{code: code, msg: Msg(code), details: details})
}

// Newf .
func Newf(code Code, errs ...error) error {
	s := &Status{code: code, msg: Msg(code), details: make([]interface{}, len(errs))}
	for i := range errs {
		s.details[i] = errs[i]
	}
	return errors.WithStack(s)
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

package apptk

import (
	"bytes"
	"github.com/golang/protobuf/jsonpb"
	"sync"
)

// var
var (
	// bind
	PBBind   = protobufBinding{}
	JSONBind = jsonBinding{}

	// JSONHandler emit empty
	JSONHandler = &jsonpb.Marshaler{
		OrigName:     true,
		EmitDefaults: true,
	}

	// JSONOmitHandler omit empty
	JSONOmitHandler = &jsonpb.Marshaler{
		OrigName:     true,
		EmitDefaults: false,
	}
)

// context type
const (
	ContentTypeKey  = "Content-Type"
	ContentTypeJSON = "application/json;charset=utf8"
	ContentTypePB   = "application/x-protobuf"
)

// LoggerInterface .
type LoggerInterface interface {
	INFO(err error)
	ERROR(err error)
}

// Log LoggerInterface impl
type Log struct{}

// INFO LoggerInterface impl
func (Log) INFO(err error) {
	return
}

// ERROR LoggerInterface impl
func (Log) ERROR(err error) {
	return
}

// bufPool .
var bufPool = &sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// getBuf .
func getBuf() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

// putBuf .
func putBuf(buf *bytes.Buffer) {
	buf.Reset()
	bufPool.Put(buf)
}

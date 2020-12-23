package tkapp

import (
	"bytes"
	"github.com/golang/protobuf/jsonpb"
	"log"
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
	Error(err error)
}

// Log LoggerInterface impl
type Log struct{}

// Error LoggerInterface impl
func (Log) Error(err error) {
	log.Println(err)
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

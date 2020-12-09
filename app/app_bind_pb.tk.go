package apptk

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
)

// protobufBinding .
type protobufBinding struct{}

// Name .
func (protobufBinding) Name() string {
	return "protobuf"
}

// Bind .
func (b protobufBinding) Bind(req *http.Request, obj interface{}) error {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(buf, obj.(proto.Message)); err != nil {
		return err
	}
	return nil
}

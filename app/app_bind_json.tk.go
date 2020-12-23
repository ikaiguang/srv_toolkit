package tkapp

import (
	"encoding/json"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
	"io/ioutil"
	"net/http"
)

// jsonBinding .
type jsonBinding struct{}

// Name .
func (jsonBinding) Name() string {
	return "json"
}

// Bind .
func (b jsonBinding) Bind(req *http.Request, obj interface{}) (err error) {
	// read
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		if err == io.EOF {
			err = nil
			return err
		}
		err = errors.WithStack(err)
		return err
	}

	// unmarshal
	if in, ok := obj.(proto.Message); ok {
		err = protojson.Unmarshal(buf, in)
	} else {
		err = json.Unmarshal(buf, obj)
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

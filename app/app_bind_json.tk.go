package apptk

import (
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"net/http"
)

// jsonBinding .
type jsonBinding struct{}

// Name .
func (jsonBinding) Name() string {
	return "json"
}

// Bind .
func (b jsonBinding) Bind(req *http.Request, obj interface{}) error {
	// pb
	if in, ok := obj.(proto.Message); ok {
		if err := jsonpb.Unmarshal(req.Body, in); err != nil {
			return err
		}
		return nil
	}

	// json
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

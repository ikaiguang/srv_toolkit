package tkapp

import (
	"bytes"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	tkpb "github.com/ikaiguang/srv_toolkit/api"
	"github.com/pkg/errors"
	"io"
)

// UnmarshalPB .
func UnmarshalPB(buf []byte, result proto.Message) (resp *tkpb.Response, err error) {
	resp = &tkpb.Response{}
	err = proto.Unmarshal(buf, resp)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = ptypes.UnmarshalAny(resp.Data, result)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UnmarshalJSON .
func UnmarshalJSON(reader io.Reader, result proto.Message) (resp *tkpb.Response, err error) {
	resp = &tkpb.Response{}
	err = jsonpb.Unmarshal(reader, resp)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = ptypes.UnmarshalAny(resp.Data, result)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UnmarshalJSONBytes .
func UnmarshalJSONBytes(jsonBytes []byte, result proto.Message) (resp *tkpb.Response, err error) {
	resp = &tkpb.Response{}
	err = jsonpb.Unmarshal(bytes.NewReader(jsonBytes), resp)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = ptypes.UnmarshalAny(resp.Data, result)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

package tkapp

import (
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
	etk "github.com/ikaiguang/srv_toolkit/error"
	"google.golang.org/protobuf/proto"
	"net/http"
)

// Bind .
func Bind(ctx *blademaster.Context, in proto.Message) (err error) {
	if ctx.Request.Method == http.MethodGet {
		err = binding.Form.Bind(ctx.Request, in)
		if err != nil {
			err = etk.Newf(etk.InvalidParameters, err)
		}
		return
	}

	// bind
	switch contentType := ctx.Request.Header.Get(ContentTypeKey); contentType {
	case ContentTypePB:
		err = PBBind.Bind(ctx.Request, in)
	case ContentTypeJSON:
		err = JSONBind.Bind(ctx.Request, in)
	case binding.MIMEJSON:
		err = JSONBind.Bind(ctx.Request, in)
	case binding.MIMEXML, binding.MIMEXML2:
		err = binding.XML.Bind(ctx.Request, in)
	default:
		err = binding.Form.Bind(ctx.Request, in)
	}
	if err != nil {
		err = etk.Newf(etk.InvalidParameters, err)
		return
	}
	return
}

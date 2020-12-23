package tkapp

import (
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
	etk "github.com/ikaiguang/srv_toolkit/error"
	"google.golang.org/protobuf/proto"
	"strings"
)

// Bind .
func Bind(ctx *blademaster.Context, in proto.Message) (err error) {
	// bind
	switch contentType := ctx.Request.Header.Get(ContentTypeKey); contentType {
	case ContentTypePB:
		err = ctx.BindWith(in, PBBind)
	case ContentTypeJSON:
		err = ctx.BindWith(in, JSONBind)
	case binding.MIMEJSON:
		err = ctx.BindWith(in, JSONBind)
	default:
		if strings.HasPrefix(contentType, binding.MIMEJSON) {
			err = ctx.BindWith(in, JSONBind)
		} else {
			err = ctx.BindWith(in, binding.Default(ctx.Request.Method, contentType))
		}
	}
	if err != nil {
		err = etk.Newf(etk.InvalidParameters, err)
		return
	}
	return
}

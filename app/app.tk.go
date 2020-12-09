package apptk

import (
	"encoding/json"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	tkpb "github.com/ikaiguang/srv_toolkit/api"
	etk "github.com/ikaiguang/srv_toolkit/error"
	tk "github.com/ikaiguang/srv_toolkit/toolkit"
	"github.com/pkg/errors"
	"net/http"
)

// init
var (
	// error detail
	omitDetail bool

	// logger
	logger LoggerInterface = Log{}
)

// SetLogger
func SetLogger(handler LoggerInterface) {
	logger = handler
}

// OmitDetail set omit error detail
func OmitDetail() {
	omitDetail = true
}

// PB response
func PB(ctx *blademaster.Context, data proto.Message) {
	// any data
	anyData, err := ptypes.MarshalAny(data)
	if err != nil {
		err = errors.WithStack(err)
		PBError(ctx, err)
		return
	}
	resp := &tkpb.Response{
		Code: etk.SUCCESS.Code(),
		Msg:  etk.Msg(etk.SUCCESS),
		Data: anyData,
	}

	// resp
	bodyBytes, err := proto.Marshal(resp)
	if err != nil {
		err = errors.WithStack(err)
		PBError(ctx, err)
		return
	}
	ctx.Bytes(http.StatusOK, ContentTypePB, bodyBytes)
	ctx.Abort()
	return
}

// PBError response
func PBError(ctx *blademaster.Context, err error) {
	loggingError(err)

	resp := errorRes(err)

	// resp
	bodyBytes, err := proto.Marshal(resp)
	if err != nil {
		err = errors.WithStack(err)
		responseErrorPB(ctx, err)
		return
	}
	ctx.Bytes(http.StatusOK, ContentTypePB, bodyBytes)
	ctx.Abort()
	return
}

// responseErrorPB response
func responseErrorPB(ctx *blademaster.Context, err error) {
	loggingError(err)

	resp := errorInit(err)

	// resp
	bodyBytes, _ := proto.Marshal(resp)
	ctx.Bytes(http.StatusOK, ContentTypePB, bodyBytes)
	ctx.Abort()
}

// JSON response emit defaults
func JSON(ctx *blademaster.Context, data proto.Message) {
	// any data
	anyData, err := ptypes.MarshalAny(data)
	if err != nil {
		err = errors.WithStack(err)
		JSONError(ctx, err)
		return
	}
	resp := &tkpb.Response{
		Code: etk.SUCCESS.Code(),
		Msg:  etk.Msg(etk.SUCCESS),
		Data: anyData,
	}

	// marshal
	var buf = getBuf()
	defer putBuf(buf)
	if err := JSONHandler.Marshal(buf, resp); err != nil {
		err = errors.WithStack(err)
		JSONError(ctx, err)
		return
	}

	// resp
	ctx.Bytes(http.StatusOK, ContentTypeJSON, buf.Bytes())
	ctx.Abort()
	return
}

// JSONOmitempty response omit empty
func JSONOmitempty(ctx *blademaster.Context, data proto.Message) {
	// any data
	anyData, err := ptypes.MarshalAny(data)
	if err != nil {
		err = errors.WithStack(err)
		JSONError(ctx, err)
		return
	}

	resp := &tkpb.Response{
		Code: etk.SUCCESS.Code(),
		Msg:  etk.Msg(etk.SUCCESS),
		Data: anyData,
	}

	// marshal
	var buf = getBuf()
	defer putBuf(buf)
	if err := JSONOmitHandler.Marshal(buf, resp); err != nil {
		err = errors.WithStack(err)
		JSONError(ctx, err)
		return
	}

	// resp
	ctx.Bytes(http.StatusOK, ContentTypeJSON, buf.Bytes())
	ctx.Abort()
	return
}

// JSONError response
func JSONError(ctx *blademaster.Context, err error) {
	loggingError(err)

	resp := errorRes(err)

	// marshal
	var buf = getBuf()
	defer putBuf(buf)
	if err := JSONHandler.Marshal(buf, resp); err != nil {
		err = errors.WithStack(err)
		responseErrorJSON(ctx, err)
		return
	}

	// resp
	ctx.Bytes(http.StatusOK, ContentTypeJSON, buf.Bytes())
	ctx.Abort()
	return
}

// responseErrorJSON response
func responseErrorJSON(ctx *blademaster.Context, err error) {
	loggingError(err)

	resp := errorInit(err)

	// resp
	bodyBytes, _ := json.Marshal(resp)
	ctx.Bytes(http.StatusOK, ContentTypeJSON, bodyBytes)
	ctx.Abort()
}

// errorRes .
func errorRes(err error) (resp *tkpb.Response) {
	resp = errorInit(err)

	s, ok := etk.FromError(err)
	if ok {
		resp.Code = s.Code.Code()
		resp.Msg = s.Msg
	}
	return
}

// errorInit .
func errorInit(err error) (resp *tkpb.Response) {
	resp = &tkpb.Response{
		Code: etk.ERROR.Code(),
		Msg:  err.Error(),
	}
	if !omitDetail {
		resp.Detail = err.Error()
	}
	return
}

// loggingError .
func loggingError(err error) {
	if logger == nil {
		err = errors.New("请先实例化logger，实例化方法：appkt.SetLogger()；")
		tk.Fatal(err)
		return
	}
	logger.ERROR(err)
}

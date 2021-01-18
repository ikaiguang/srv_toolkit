package testdata

import (
	"bytes"
	"context"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	xtime "github.com/go-kratos/kratos/pkg/time"
	tkapp "github.com/ikaiguang/srv_toolkit/app"
	tkcurl "github.com/ikaiguang/srv_toolkit/curl"
	pingpb "github.com/ikaiguang/srv_toolkit/srv_hello/api/ping"
	pinghanlder "github.com/ikaiguang/srv_toolkit/srv_hello/internal/handler/ping"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

// pingURL .
func pingURL() string {
	return "http://" + confHTTP.Addr + "/ping"
}

func TestHTTPProtobuf(t *testing.T) {
	lazyServe()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ccConf := &blademaster.ClientConfig{
		Dial:      xtime.Duration(3 * time.Second),
		Timeout:   xtime.Duration(40 * time.Second),
		KeepAlive: xtime.Duration(40 * time.Second),
	}
	cc := blademaster.NewClient(ccConf)

	httParam := tkcurl.ProtobufRequestParam()
	httParam.URL = pingURL()
	httpReq, err := tkcurl.NewRequest(http.MethodGet, httParam)
	assert.Nil(t, err)

	bodyBytes, err := cc.Raw(ctx, httpReq)
	assert.Nil(t, err)

	ret := &pingpb.PingResp{}
	_, err = tkapp.UnmarshalPB(bodyBytes, ret)
	assert.Nil(t, err)
	assert.Equal(t, ret.String(), pinghanlder.PingPong.String())
}

func TestHTTPJSON(t *testing.T) {
	lazyServe()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ccConf := &blademaster.ClientConfig{
		Dial:      xtime.Duration(3 * time.Second),
		Timeout:   xtime.Duration(40 * time.Second),
		KeepAlive: xtime.Duration(40 * time.Second),
	}
	cc := blademaster.NewClient(ccConf)

	httParam := tkcurl.JSONRequestParam()
	httParam.URL = pingURL()
	httpReq, err := tkcurl.NewRequest(http.MethodGet, httParam)
	assert.Nil(t, err)

	bodyBytes, err := cc.Raw(ctx, httpReq)
	assert.Nil(t, err)

	ret := &pingpb.PingResp{}
	_, err = tkapp.UnmarshalJSON(bytes.NewReader(bodyBytes), ret)
	assert.Nil(t, err)
	assert.Equal(t, ret.String(), pinghanlder.PingPong.String())
}

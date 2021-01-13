package pingroute

import (
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	tkapp "github.com/ikaiguang/srv_toolkit/app"
	pingpb "github.com/ikaiguang/srv_toolkit/srv_hello/api/ping"
	pinghanlder "github.com/ikaiguang/srv_toolkit/srv_hello/internal/handler/ping"
)

// PingRoute route
var PingRoute = &ping{}

// ping route
type ping struct{}

// Ping .
func (s *ping) Ping(ctx *blademaster.Context) {
	// input parameters
	var (
		in  = &pingpb.PingReq{}
		err error
	)

	err = tkapp.Bind(ctx, in)
	if err != nil {
		tkapp.Error(ctx, err)
		return
	}

	// data
	data, err := pinghanlder.PingHandler.Ping(ctx, in)
	if err != nil {
		tkapp.Error(ctx, err)
		return
	}
	tkapp.Success(ctx, data)
	return
}

package pinghanlder

import (
	"context"
	pingpb "github.com/ikaiguang/srv_toolkit/srv_hello/api/ping"
)

// Ping
var (
	// PingHandler handler
	PingHandler pingpb.CkgPingServer = &ping{}

	// PingPong pong
	PingPong = pingpb.PingResp{Message: "pong"}
)

// ping handler
type ping struct{ pingpb.UnimplementedCkgPingServer }

// Ping .
func (s *ping) Ping(ctx context.Context, in *pingpb.PingReq) (res *pingpb.PingResp, err error) {
	res = &PingPong
	//err = tke.New(e.HelloTestError)
	return
}

package routes

import (
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	pingpb "github.com/ikaiguang/srv_toolkit/srv_hello/api/ping"
	pinghanlder "github.com/ikaiguang/srv_toolkit/srv_hello/internal/handler/ping"
)

// RegisterGRPC .
func RegisterGRPC(ws *warden.Server) {
	// ping
	pingpb.RegisterCkgPingServer(ws.Server(), pinghanlder.PingHandler)
}

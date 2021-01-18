package testdata

import (
	"context"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	pingpb "github.com/ikaiguang/srv_toolkit/srv_hello/api/ping"
	pinghanlder "github.com/ikaiguang/srv_toolkit/srv_hello/internal/handler/ping"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGRPCBasic(t *testing.T) {
	lazyServe()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cc := warden.DefaultClient()
	conn, err := cc.Dial(ctx, confGRPC.Addr)
	assert.Nil(t, err)

	// request
	pingClient := pingpb.NewCkgPingClient(conn)
	pingReq := &pingpb.PingReq{}
	pingRes, err := pingClient.Ping(ctx, pingReq)
	assert.Nil(t, err)

	// result
	assert.Equal(t, pingRes.String(), pinghanlder.PingPong.String())
}

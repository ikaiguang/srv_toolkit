package testdata

import (
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	"testing"
)

func TestGRPCBasic(t *testing.T) {
	lazyServe()

	clientConf := &warden.ClientConfig{

	}
	warden.DefaultClient()
	client := warden.NewClient(clientConf)
	_ = client
}

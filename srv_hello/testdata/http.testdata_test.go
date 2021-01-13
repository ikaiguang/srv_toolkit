package testdata

import (
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"testing"
)

func TestHTTPProtobuf(t *testing.T) {
	lazyServe()

	clientConf := &blademaster.ClientConfig{

	}
	client := blademaster.NewClient(clientConf)
	_ = client
}

func TestHTTPJSON(t *testing.T) {
	lazyServe()

	clientConf := &blademaster.ClientConfig{

	}
	client := blademaster.NewClient(clientConf)
	_ = client
}

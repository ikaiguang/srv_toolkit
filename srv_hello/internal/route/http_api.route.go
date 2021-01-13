package routes

import (
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	pingroute "github.com/ikaiguang/srv_toolkit/srv_hello/internal/route/ping"
)

// RegisterApi .
func RegisterApi(engine *blademaster.Engine) {
	// api v1
	apiV1 := "/api/v1"
	apiV1Group := engine.Group(apiV1)

	pingroute.RegisterApi(apiV1Group)
}

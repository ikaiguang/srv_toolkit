package routes

import (
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	pingroute "github.com/ikaiguang/srv_toolkit/srv_hello/internal/route/ping"
)

// RegisterWeb .
func RegisterWeb(engine *blademaster.Engine) {
	// web v1
	webV1 := "/web/v1"
	webV1Group := engine.Group(webV1)

	pingroute.RegisterWeb(webV1Group)
}


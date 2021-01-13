package routes

import (
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	pingroute "github.com/ikaiguang/srv_toolkit/srv_hello/internal/route/ping"
)

// RegisterHTTP .
func RegisterHTTP(engine *blademaster.Engine) {
	RegisterGlobal(engine)
	RegisterWeb(engine)
	RegisterApi(engine)
}

// RegisterGlobal .
func RegisterGlobal(engine *blademaster.Engine) {
	// global
	engine.GET("/ping", pingroute.PingRoute.Ping)
}

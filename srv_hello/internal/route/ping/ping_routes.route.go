package pingroute

import "github.com/go-kratos/kratos/pkg/net/http/blademaster"

// RegisterWeb .
func RegisterWeb(group *blademaster.RouterGroup) {
	group.GET("/ping", PingRoute.Ping)
}

// RegisterApi .
func RegisterApi(group *blademaster.RouterGroup) {
	group.GET("/ping", PingRoute.Ping)
}

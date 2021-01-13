package e

import tke "github.com/ikaiguang/srv_toolkit/error"

func init() {
	tke.Register(msg)
}

var msg = map[tke.Code]string{
	HelloTestError: "HelloWorld",
}

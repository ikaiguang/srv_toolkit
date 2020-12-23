package tkflag

import (
	"flag"
	"fmt"
)

// Setup .
func Setup() {
	var usage = `app version : app/1.0.1
Usage: app [-h]
`
	// 重写 -h 说明
	flag.CommandLine.Usage = func() {
		fmt.Println(usage)
	}
	flag.Parse()
}

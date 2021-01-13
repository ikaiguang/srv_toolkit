package testdata

import (
	"github.com/ikaiguang/srv_toolkit/srv_hello/internal/setup"
	"sync"
	"time"
)

// var
var (
	onceServer sync.Once
)

// lazyServe .
func lazyServe() {
	onceServer.Do(func() {
		go setup.Serve()
		time.Sleep(time.Second)
	})
}

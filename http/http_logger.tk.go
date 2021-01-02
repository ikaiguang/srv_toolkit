package tkhttp

import (
	"fmt"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/metadata"
	"github.com/go-kratos/kratos/pkg/stat/metric"
	"strconv"
	"time"
)
const (
	clientNamespace = "http_client"
	serverNamespace = "http_server"
)

var (
	_metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Labels:    []string{"path", "caller", "method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})
	_metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "caller", "method", "code"},
	})
)

// Logger is logger  middleware
func Logger() blademaster.HandlerFunc {
	const noUser = "no_user"
	return func(c *blademaster.Context) {
		now := time.Now()
		req := c.Request
		path := req.URL.Path
		params := req.Form
		var quota float64
		if deadline, ok := c.Context.Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}

		c.Next()

		err := c.Error
		cerr := ecode.Cause(err)
		dt := time.Since(now)
		caller := metadata.String(c, metadata.Caller)
		if caller == "" {
			caller = noUser
		}

		if len(c.RoutePath) > 0 {
			_metricServerReqCodeTotal.Inc(c.RoutePath[1:], caller, req.Method, strconv.FormatInt(int64(cerr.Code()), 10))
			_metricServerReqDur.Observe(int64(dt/time.Millisecond), c.RoutePath[1:], caller, req.Method)
		}

		lf := log.Infov
		errmsg := ""
		isSlow := dt >= (time.Millisecond * 500)
		if err != nil {
			errmsg = err.Error()
			lf = log.Errorv
			if cerr.Code() > 0 {
				lf = log.Warnv
			}
		} else {
			if isSlow {
				lf = log.Warnv
			}
		}
		lf(c,
			log.KVString("method", req.Method),
			log.KVString("ip", c.RemoteIP()),
			log.KVString("user", caller),
			log.KVString("path", path),
			log.KVString("params", params.Encode()),
			log.KVInt("ret", cerr.Code()),
			log.KVString("msg", cerr.Message()),
			log.KVString("stack", fmt.Sprintf("%+v", err)),
			log.KVString("err", errmsg),
			log.KVFloat64("timeout_quota", quota),
			log.KVFloat64("ts", dt.Seconds()),
			log.KVString("source", "http-access-log"),
		)
	}
}

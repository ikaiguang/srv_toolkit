package tklog

import (
	"context"
	"github.com/go-kratos/kratos/pkg/conf/env"
	"github.com/go-kratos/kratos/pkg/net/metadata"
	"github.com/go-kratos/kratos/pkg/net/trace"
	"go.uber.org/zap"
)

// kratos
const (
	_timeFormat = "2006-01-02T15:04:05.999999"

	// log level defined in level.go.
	_levelValue = "level_value"
	//  log level name: INFO, WARN...
	_level = "level"
	// log time.
	_time = "time"
	// request path.
	// _title = "title"
	// log file.
	_source = "source"
	// common log filed.
	_log = "log"
	// app name.
	_appID = "app_id"
	// container ID.
	_instanceID = "instance_id"
	// uniq ID from trace.
	_tid = "traceid"
	// request time.
	// _ts = "ts"
	// requester.
	_caller = "caller"
	// container environment: prod, pre, uat, fat.
	_deplyEnv = "env"
	// container area.
	_zone = "zone"
	// mirror flag
	_mirror = "mirror"
	// color.
	_color = "color"
	// env_color
	_envColor = "env_color"
	// cluster.
	_cluster = "cluster"
)

// AddExtraField .
func AddExtraField(ctx context.Context) (fields []zap.Field) {
	if t, ok := trace.FromContext(ctx); ok {
		fields = append(fields, zap.String(_tid, t.TraceID()))
	}
	if caller := metadata.String(ctx, metadata.Caller); caller != "" {
		fields = append(fields, zap.String(_caller, caller))
	}
	if color := metadata.String(ctx, metadata.Color); color != "" {
		fields = append(fields, zap.String(_color, color))
	}
	if env.Color != "" {
		fields = append(fields, zap.String(_envColor, env.Color))
	}
	if cluster := metadata.String(ctx, metadata.Cluster); cluster != "" {
		fields = append(fields, zap.String(_cluster, cluster))
	}
	fields = append(fields, zap.String(_deplyEnv, env.DeployEnv))
	fields = append(fields, zap.String(_zone, env.Zone))
	//fields = append(fields, zap.String(_appID, c.Family))
	fields = append(fields, zap.String(_appID, env.AppID))
	//fields = append(fields, zap.String(_instanceID, c.Host))
	fields = append(fields, zap.String(_instanceID, env.Hostname))
	if metadata.String(ctx, metadata.Mirror) != "" {
		fields = append(fields, zap.Bool(_mirror, true))
	}
	return
}

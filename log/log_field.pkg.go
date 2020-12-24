package tklog

import (
	"go.uber.org/zap"
	"time"
)

// KVString construct Field with string value.
func KVString(key string, value string) zap.Field {
	return zap.String(key, value)
}

// KVInt construct Field with int value.
func KVInt(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// KVInt64 construct D with int64 value.
func KVInt64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

// KVUint construct Field with uint value.
func KVUint(key string, value uint) zap.Field {
	return zap.Uint(key, value)
}

// KVUint64 construct Field with uint64 value.
func KVUint64(key string, value uint64) zap.Field {
	return zap.Uint64(key, value)
}

// KVFloat32 construct Field with float32 value.
func KVFloat32(key string, value float32) zap.Field {
	return zap.Float32(key, value)
}

// KVFloat64 construct Field with float64 value.
func KVFloat64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

// KVDuration construct Field with Duration value.
func KVDuration(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

// KV return a log kv for logging field.
// NOTE: use KV{type name} can avoid object alloc and get better performance. []~(￣▽￣)~*干杯
func KV(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

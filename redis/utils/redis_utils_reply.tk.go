package tkru

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

// IsRedisNil .
func IsRedisNil(err error) bool {
	return err == redis.ErrNil
}

// Int .
func Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

// Int64 .
func Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

// Uint64 .
func Uint64(reply interface{}, err error) (uint64, error) {
	return redis.Uint64(reply, err)
}

// Float64 .
func Float64(reply interface{}, err error) (float64, error) {
	return redis.Float64(reply, err)
}

// String .
func String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

// Bytes .
func Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}

// Bool .
func Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

// Values .
func Values(reply interface{}, err error) ([]interface{}, error) {
	return redis.Values(reply, err)
}

// Float64s .
func Float64s(reply interface{}, err error) ([]float64, error) {
	return redis.Float64s(reply, err)
}

// Strings .
func Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

// ByteSlices .
func ByteSlices(reply interface{}, err error) ([][]byte, error) {
	return redis.ByteSlices(reply, err)
}

// Int64s .
func Int64s(reply interface{}, err error) ([]int64, error) {
	return redis.Int64s(reply, err)
}

// Ints .
func Ints(reply interface{}, err error) ([]int, error) {
	return redis.Ints(reply, err)
}

// BytesMap .
func BytesMap(result interface{}, err error) (map[string][]byte, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: BytesMap expects even number of values result")
	}
	m := make(map[string][]byte, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		value, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: BytesMap key not a bulk string value")
		}
		m[string(key)] = value
	}
	return m, nil
}

// StringMap .
func StringMap(result interface{}, err error) (map[string]string, error) {
	return redis.StringMap(result, err)
}

// IntMap .
func IntMap(result interface{}, err error) (map[string]int, error) {
	return redis.IntMap(result, err)
}

// Int64Map .
func Int64Map(result interface{}, err error) (map[string]int64, error) {
	return redis.Int64Map(result, err)
}

// Positions .
func Positions(result interface{}, err error) ([]*[2]float64, error) {
	return redis.Positions(result, err)
}

// Uint64s .
func Uint64s(reply interface{}, err error) ([]uint64, error) {
	return redis.Uint64s(reply, err)
}

// Uint64Map .
func Uint64Map(result interface{}, err error) (map[string]uint64, error) {
	return redis.Uint64Map(result, err)
}

// SlowLogs .
func SlowLogs(result interface{}, err error) ([]redis.SlowLog, error) {
	return redis.SlowLogs(result, err)
}

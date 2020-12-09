package tk

import "runtime"

// File file line
func File(skips ...int) (file string, line int) {
	skip := 1
	if len(skips) > 0 {
		skip = skips[0]
	}
	if skip < 0 {
		skip = 0
	}
	_, file, line, _ = runtime.Caller(skip)
	return
}

package common

import "time"

// CurrentTimeMillis calculates the millis elapsed since Unix Epoch
func CurrentTimeMillis() int64 {
	return time.Now().UnixNano() / time.Millisecond.Nanoseconds()
}

package common

import "time"

func CurrentTimeMillis() int64 {
    return time.Now().UnixNano() / time.Millisecond.Nanoseconds()
}
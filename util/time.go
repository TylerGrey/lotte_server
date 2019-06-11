package util

import "time"

// LocalTimeUnix unixtime
func LocalTimeUnix() int64 {
	loc, _ := time.LoadLocation("Asia/Seoul")
	t := time.Now().In(loc)
	return t.Unix()
}

package util

import (
	"strconv"
	"time"
)

type Stamp func() string

func Unixtime() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

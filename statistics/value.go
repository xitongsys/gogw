package statistics

import (
	"time"
)

type ValueWithTime struct {
	timestamp time.Time
	value     interface{}
}

func NewValueWithTime(v interface{}, t time.Time) *ValueWithTime {
	return &ValueWithTime{
		timestamp: t,
		value:     v,
	}
}

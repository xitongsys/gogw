package statistics

import (
	"time"
)

type ValueWithTime struct {
	Timestamp time.Time
	Value     interface{}
}

func NewValueWithTime(v interface{}, t time.Time) *ValueWithTime {
	return &ValueWithTime{
		Timestamp: t,
		Value:     v,
	}
}

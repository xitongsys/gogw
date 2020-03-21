package statistics

import (
	"time"
)

type WindowStat struct {
	length time.Duration

	start, end time.Time
	curCount int64
	curValue interface{}

	op func(v interface{}, cnt int64, nv interface{}) interface{}
}

func NewWindowStat(length time.Duration, 
	op func(v interface{}, cnt int64, nv interface{}) interface{}, capacity int) *WindowStat {
	return & WindowStat {
		length: length,
		start: time.Now(),
		end: time.Now().Add(length),

		curCount: 0,
		curValue: nil,

		op: op,
	}
}

func (w *WindowStat) Add(v interface{}) {
	t := time.Now()
	if t.Unix() >= w.start.Unix() && t.Unix() < w.end.Unix() {
		w.curValue = w.op(w.curValue, w.curCount, v)
		w.curCount++

	}else if t.Unix() >= w.end.Unix() {
		w.start = t
		w.end = w.end.Add(w.length)
		w.curValue = w.op(nil, 0, v)
		w.curCount = 1
	}
}

func (w *WindowStat) GetLatest() interface{} {
	return w.curValue
}
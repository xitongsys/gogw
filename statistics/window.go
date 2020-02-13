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

	histValues *Queue
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

		histValues: NewQueue(capacity),
	}
}

func (w *WindowStat) Add(v interface{}) {
	t := time.Now()
	if t.Unix() >= w.start.Unix() && t.Unix() < w.end.Unix() {
		w.curValue = w.op(w.curValue, w.curCount, v)
		w.curCount++

	}else if t.Unix() >= w.end.Unix() {
		for t.Unix() >= w.end.Unix(){
			w.histValues.Push(NewValueWithTime(w.curValue, w.end))
			w.curCount, w.curValue = 0, nil
			w.start = w.start.Add(w.length)
			w.end = w.end.Add(w.length)
		}
		w.Add(v)
	}
}

func (w *WindowStat) GetLatest() (*ValueWithTime, error) {
	r, e := w.histValues.Back()
	if e != nil {
		return nil, e
	}

	return r.(*ValueWithTime), nil
}


func (w *WindowStat) GetAll() []*ValueWithTime {
	vs := w.histValues.All()
	res := []*ValueWithTime{}
	for _, v := range vs {
		res = append(res, v.(*ValueWithTime))
	}
	return res
}
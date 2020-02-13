package statistics

import (
	"time"
)

type PointStat struct {
	values *Queue
}

func NewPointStat(size int) *PointStat {
	return & PointStat{
		values: NewQueue(size),
	}
}

func (p *PointStat) GetLatest() (*ValueWithTime, error) {
	r, e := p.values.Back()
	if e != nil {
		return nil, e
	}
	return r.(*ValueWithTime), nil
}

func (p *PointStat) GetAll() []*ValueWithTime {
	vs := p.values.All()
	res := []*ValueWithTime{}
	for _, v := range vs {
		res = append(res, v.(*ValueWithTime))
	}
	return res
}

func (p *PointStat) Add(v interface{}) {
	p.values.Push(NewValueWithTime(v, time.Now()))
}
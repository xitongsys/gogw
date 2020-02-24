package monitor

import (
	"gogw/statistics"
	"time"
)

var WINLENGTH = time.Second

type SpeedMonitor struct {
	SpeedRecord *statistics.WindowStat
}

func NewSpeedMonitor() *SpeedMonitor {
	return & SpeedMonitor{
		SpeedRecord: statistics.NewWindowStat(WINLENGTH, statistics.Sum, 100),
	}
}

func (sm *SpeedMonitor) Add(size int64) {
	if size > 0 {
		sm.SpeedRecord.Add(size)
	}
}

func (sm *SpeedMonitor) GetSpeed() int {
	vs, err := sm.SpeedRecord.GetLatest()
	if err != nil || vs.Value == nil {
		return 0
	}
	
	return int(vs.Value.(int64))
}

package service

import "time"

type timer struct {
	startTime time.Time
	endTime   time.Time
}

func NewTimer() *timer {
	return &timer{}
}

func (t *timer) Start() {
	t.startTime = time.Now()
}

func (t *timer) Stop() {
	t.endTime = time.Now()
}

func (t timer) GetSeconds() float64 {
	return t.endTime.Sub(t.startTime).Seconds()
}

package main

import (
	"log"
	"time"
)

type Timer struct {
	timings []*timing
}

type timing struct {
	tag  string
	time time.Time
}

func NewTimer() *Timer {
	return &Timer{}
}

func (t *Timer) Add(tag string) {
	stamp := &timing{tag: tag, time: time.Now()}
	t.timings = append(t.timings, stamp)
}

func (t *Timer) Report() {
	for i := 1; i < len(t.timings); i++ {
		a := t.timings[i-1]
		b := t.timings[i]
		interval_ms := b.time.Sub(a.time).Nanoseconds() / 1000000
		log.Printf("Timer: [%v] to [%v] %v ms", a.tag, b.tag, interval_ms)
	}
}

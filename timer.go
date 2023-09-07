package main

import "time"

type ITimer interface {
	Start() time.Time
	Stop() time.Time
	Pause()
	Resume()

	Running() bool
	Paused() bool

	String() string
}

type Timer struct {
	start, end time.Time
	pausedAt   time.Time
	ballast    time.Duration
}

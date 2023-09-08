package main

import (
	"fmt"
	"time"
)

type ITimer interface {
	Start() time.Time
	Stop() time.Time
	Pause() time.Time
	Resume()

	Running() bool
	Paused() bool
	Stopped() bool

	String() string
}

type Timer struct {
	start, end *time.Time
	pausedAt   *time.Time
	ballast    time.Duration
}

func (t *Timer) Start() time.Time {

	now := time.Now()
	t.start = &now
	return now
}

func (t *Timer) Stop() time.Time {
	now := time.Now()
	t.end = &now
	return now
}

func (t *Timer) Pause() time.Time {
	now := time.Now()
	t.pausedAt = &now
	return now
}

func (t *Timer) Resume() {
	t.ballast += time.Since(*t.pausedAt)
	t.pausedAt = nil
}

func (t *Timer) Running() bool {
	return t.start != nil && t.end == nil
}

func (t *Timer) Paused() bool {
	return t.Running() && t.pausedAt != nil
}

func (t *Timer) String() string {
	if t.Running() {
		return fmt.Sprint(time.Since(*t.start) - t.ballast)
	}

	// TODO: handle showing end time

	// TODO: actually calculate this value
	return "00:00:00.000"
}

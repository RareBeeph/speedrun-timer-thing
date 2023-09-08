package main

import (
	"fmt"
	"time"
)

type ITimer interface {
	Start() time.Time
	Stop() time.Time
	Restart() time.Time
	Pause() time.Time
	Resume()

	Idle() bool
	Running() bool
	Paused() bool
	Stopped() bool

	String() string
}

type Timer struct {
	start, end *time.Time
	ballast    time.Duration
}

/*
	State machine transition table:

	Stopped (i.e. end != nil):
		Start() -> Running; Stop() -> Stopped; Restart() -> Idle; Pause() -> Stopped; Resume() -> Stopped
	Running (i.e. end = nil, start != nil):
		Start() -> Running; Stop() -> Stopped; Restart() -> Idle; Pause() -> Paused; Resume() -> Running
	Paused (i.e. end = nil, start = nil, ballast != 0):
		Start() -> Paused; Stop() -> Stopped; Restart() -> Idle; Pause() -> Paused; Resume() -> Running
	Idle (i.e. end = nil, start = nil, ballast = 0):
		Start() -> Running; Stop() -> Idle; Restart() -> Idle; Pause() -> Idle; Resume() -> Idle
*/

func (t *Timer) Start() time.Time {
	now := time.Now()

	// should never occur, but just in case
	if t.Paused() {
		return now
	}

	t.end = nil
	t.ballast = time.Duration(0)
	t.start = &now
	return now
}

func (t *Timer) Stop() time.Time {
	now := time.Now()

	// should never occur, but just in case
	if t.Idle() {
		return now
	}

	if t.Running() {
		t.ballast += time.Since(*t.start)
	}
	t.start = nil
	t.end = &now
	return now
}

func (t *Timer) Restart() time.Time {
	now := time.Now()
	t.start = nil
	t.end = nil
	t.ballast = time.Duration(0)
	return now
}

func (t *Timer) Pause() time.Time {
	now := time.Now()

	// should never occur, but just in case
	if !t.Running() {
		return now
	}

	t.ballast += time.Since(*t.start)
	t.start = nil
	return now
}

func (t *Timer) Resume() {
	now := time.Now()

	// should never occur, but just in case
	if !t.Paused() {
		return
	}

	t.start = &now
}

func (t *Timer) Idle() bool {
	return t.start == nil && t.end == nil && t.ballast == time.Duration(0)
}

func (t *Timer) Running() bool {
	return t.start != nil && t.end == nil
}

func (t *Timer) Paused() bool {
	return t.start == nil && t.end == nil && t.ballast != time.Duration(0)
}

func (t *Timer) Stopped() bool {
	return t.end != nil
}

func (t *Timer) String() string {
	totalTime := t.ballast.Milliseconds()
	if t.Running() {
		totalTime += time.Since(*t.start).Milliseconds()
	}

	out := fmt.Sprintf("%02d:%02d.%03d", totalTime/60000%60, totalTime/1000%60, totalTime%1000)
	if totalTime >= 3600000 {
		out = fmt.Sprintf("%02d:%s", totalTime/3600000, out)
	}
	return out
}

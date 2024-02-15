package timer

import (
	"speedruntimer/timing/formatting"
	"time"
)

type Timer interface {
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
	Milliseconds() int64
}

type timer struct {
	start, end *time.Time
	run        *Run
	ballast    time.Duration
}

func New(run *Run) Timer {
	return &timer{run: run}
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

func (t *timer) Start() time.Time {
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

func (t *timer) Stop() time.Time {
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

func (t *timer) Restart() time.Time {
	now := time.Now()
	t.start = nil
	t.end = nil
	t.ballast = time.Duration(0)
	return now
}

func (t *timer) Pause() time.Time {
	now := time.Now()

	// should never occur, but just in case
	if !t.Running() {
		return now
	}

	t.ballast += time.Since(*t.start)
	t.start = nil
	return now
}

func (t *timer) Resume() {
	now := time.Now()

	// should never occur, but just in case
	if !t.Paused() {
		return
	}

	t.start = &now
}

func (t *timer) Milliseconds() int64 {
	totalTime := t.ballast.Milliseconds()
	if t.Running() {
		totalTime += time.Since(*t.start).Milliseconds()
	}

	return totalTime
}

func (t *timer) Idle() bool {
	return t.start == nil && t.end == nil && t.ballast == time.Duration(0)
}

func (t *timer) Running() bool {
	return t.start != nil && t.end == nil
}

func (t *timer) Paused() bool {
	return t.start == nil && t.end == nil && t.ballast != time.Duration(0)
}

func (t *timer) Stopped() bool {
	return t.end != nil
}

func (t *timer) String() string {
	return formatting.TimeFormatMilliseconds(t.Milliseconds())
}

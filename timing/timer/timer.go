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
	Split() time.Time
	Resume()

	Idle() bool
	Running() bool
	Paused() bool
	Stopped() bool

	String() string
	Elapsed() time.Duration
}

type timer struct {
	start, end time.Time // end is redundant info now that we have the Run pointer
	ballast    time.Duration

	run     *Run
	segment int
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

	t.end = time.Time{}
	t.ballast = time.Duration(0)
	t.start = now
	return now
}

func (t *timer) Stop() time.Time {
	now := time.Now()

	// should never occur, but just in case
	if t.Idle() {
		return now
	}

	if t.Running() {
		t.ballast += time.Since(t.start)
	}
	t.start = time.Time{}
	t.end = now
	return now
}

func (t *timer) Restart() time.Time {
	now := time.Now()
	t.start = time.Time{}
	t.end = time.Time{}
	t.ballast = time.Duration(0)
	return now
}

func (t *timer) Pause() time.Time {
	now := time.Now()

	// should never occur, but just in case
	if !t.Running() {
		return now
	}

	t.ballast += time.Since(t.start)
	t.start = time.Time{}
	return now
}

func (t *timer) Split() time.Time {
	now := time.Now()
	sinceStart := (now.Sub(t.start) * time.Millisecond) + time.Duration(t.ballast.Milliseconds())

	if !t.Running() {
		return now
	}

	if t.run.Segments[t.segment] == nil {
		return now // probably want to do some actual error reporting here
	}

	segment := t.run.Segments[t.segment]
	prev := t.previousSegment()

	segment.Split(sinceStart, prev.ActiveRunTime)

	if t.segment == len(t.run.Segments)-1 {
		t.Stop()
	}

	return now
}

func (t *timer) Resume() {
	now := time.Now()

	// should never occur, but just in case
	if !t.Paused() {
		return
	}

	t.start = now
}

func (t *timer) Idle() bool {
	return t.start.IsZero() && t.end.IsZero() && t.ballast == time.Duration(0)
}

func (t *timer) Running() bool {
	return !t.start.IsZero() && t.end.IsZero()
}

func (t *timer) Paused() bool {
	return t.start.IsZero() && t.end.IsZero() && t.ballast != time.Duration(0)
}

func (t *timer) Stopped() bool {
	return !t.end.IsZero()
}

func (t *timer) String() string {
	return formatting.TimeFormatMilliseconds(int64(t.Elapsed()))
}

// This is suitable for display but NOT for calculation
// as the time measurement occurs inside of this function
// and is not representative of the time the keypress
// event was received
func (t *timer) Elapsed() time.Duration {
	totalTime := t.ballast.Milliseconds()
	if t.Running() {
		totalTime += time.Since(t.start).Milliseconds()
	}

	return time.Duration(totalTime)
}

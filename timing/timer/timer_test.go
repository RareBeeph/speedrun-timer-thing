package timer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

var run = &Run{}

// TODO: remove this - the switch to a constructor implies it
// func TestImplementsITimer(t *testing.T) {
// 	timer := New(run)
// 	var _ = (Timer)(timer)
// }

func TestIdle(t *testing.T) {
	timer := New(run)
	assert.True(t, timer.Idle(), "timer is idle on create")

	timer = New(run)
	timer.Start()
	assert.True(t, timer.Running(), "Idle + Start() starts the timer")

	timer = New(run)
	timer.Stop()
	assert.True(t, timer.Idle(), "Idle + Stop() remains idle")

	timer = New(run)
	timer.Restart()
	assert.True(t, timer.Idle(), "Idle + Restart() remains idle")

	timer = New(run)
	timer.Pause()
	assert.True(t, timer.Idle(), "Idle + Pause() remains idle")

	timer = New(run)
	timer.Resume()
	assert.True(t, timer.Idle(), "Idle + Resume() remains idle")
}

func TestRunning(t *testing.T) {
	timer := New(run)
	timer.Start()
	timer.Start()
	assert.True(t, timer.Running(), "Running + Start() remains running")

	timer = New(run)
	timer.Start()
	timer.Stop()
	assert.True(t, timer.Stopped(), "Running + Stop() stops the timer")

	timer = New(run)
	timer.Start()
	timer.Restart()
	assert.True(t, timer.Idle(), "Running + Restart() returns to idle")

	timer = New(run)
	timer.Start()
	timer.Pause()
	assert.True(t, timer.Paused(), "Running + Pause() pauses the timer")

	timer = New(run)
	timer.Start()
	timer.Resume()
	assert.True(t, timer.Running(), "Running + Resume() remains running")
}

func TestPaused(t *testing.T) {
	timer := New(run)
	timer.Start()
	timer.Pause()
	timer.Start()
	assert.True(t, timer.Paused(), "Paused + Start() remains paused")

	timer = New(run)
	timer.Start()
	timer.Pause()
	timer.Stop()
	assert.True(t, timer.Stopped(), "Paused + Stop() stops the timer")

	timer = New(run)
	timer.Start()
	timer.Pause()
	timer.Restart()
	assert.True(t, timer.Idle(), "Paused + Restart() returns to idle")

	timer = New(run)
	timer.Start()
	timer.Pause()
	timer.Pause()
	assert.True(t, timer.Paused(), "Paused + Pause() remains paused")

	timer = New(run)
	timer.Start()
	timer.Pause()
	timer.Resume()
	assert.True(t, timer.Running(), "Paused + Resume() resumes the timer")
}

func TestStopped(t *testing.T) {
	timer := New(run)
	timer.Start()
	timer.Stop()
	timer.Start()
	assert.True(t, timer.Running(), "Stopped + Start() starts the timer")

	timer = New(run)
	timer.Start()
	timer.Stop()
	timer.Stop()
	assert.True(t, timer.Stopped(), "Stopped + Stop() remains stopped")

	timer = New(run)
	timer.Start()
	timer.Stop()
	timer.Restart()
	assert.True(t, timer.Idle(), "Stopped + Restart() returns to idle")

	timer = New(run)
	timer.Start()
	timer.Stop()
	timer.Pause()
	assert.True(t, timer.Stopped(), "Stopped + Pause() remains stopped")

	timer = New(run)
	timer.Start()
	timer.Stop()
	timer.Resume()
	assert.True(t, timer.Stopped(), "Stopped + Resume() remains stopped")
}

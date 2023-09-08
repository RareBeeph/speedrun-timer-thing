package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStarting(t *testing.T) {
	timer := &Timer{}

	assert.False(t, timer.Running(), "timer is not running on create")
	timer.Start()
	assert.True(t, timer.Running(), "Start() starts the timer")
}

func TestPausing(t *testing.T) {
	timer := &Timer{}

	timer.Start()
	assert.False(t, timer.Paused(), "timer is not paused on start")

	timer.Pause()
	assert.True(t, timer.Paused(), "Pause() pauses the timer")

	timer.Resume()
	assert.False(t, timer.Paused(), "Resume() resumes the timer")
}

// TODO: test the rest of the exposed interface

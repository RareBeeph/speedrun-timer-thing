package timing

import (
	"speedruntimer/timing/splitter"
	"speedruntimer/timing/timer"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
	TODO: make splithandler idle check not consider index 0 to always be idle

	idle: Timer.Idle(), SplitHandler.Idle()
	active: Timer.Running(), SplitHandler.Active()
	paused: Timer.Paused(), SplitHandler.Active()
	cancelled: Timer.Stopped(), SplitHandler.Active()
	finished: Timer.Stopped(), SplitHandler.Finished()

					Split		Pause		Stop
				-----------------------------------------
	idle		|	active	|	-		|	-			|
	active		|	(a/f)	|	paused	|	cancelled	|	// this row is shorthand for splithandler substates
	paused		|	-		|	active	|	cancelled	|	// this row is shorthand for splithandler substates
	cancelled	|	-		|	-		|	idle		|
	finished	|	-		|	-		|	idle		|
				-----------------------------------------

	no other states are valid
*/

func TestSplit(t *testing.T) {
	tmch := TimeMachine{Timer: &timer.Timer{}, SplitHandler: &splitter.SplitHandler{}}

	// TODO: fake this
	tmch.SplitHandler.SetSplits([]splitter.Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	assert.True(t, tmch.Timer.Idle(),
		"Timer should be idle before first split")

	tmch.Split() // timer start
	assert.True(t, tmch.Timer.Running(),
		"Timer should be running after first split input")
	assert.True(t, tmch.SplitHandler.IsBeforeFirstSplit(),
		"SplitHandler should not split on the first split input")

	tmch.Split() // first split
	assert.True(t, tmch.Timer.Running(),
		"Timer should still be running after first real split")
	assert.False(t, tmch.SplitHandler.IsBeforeFirstSplit() || tmch.SplitHandler.IsFinished(),
		"SplitHandler should still be active after first real split")

	tmch.Split() // final split
	assert.True(t, tmch.Timer.Stopped(),
		"Timer should stop after final split")
	assert.True(t, tmch.SplitHandler.IsFinished(),
		"SplitHandler should be finished after final split")

	tmch.Split() // additional inputs after final
	assert.True(t, tmch.Timer.Stopped() && tmch.SplitHandler.IsFinished(),
		"Splits while stopped should do nothing")

	tmch.Stop()  // restart; see TestStop()
	tmch.Split() // start
	tmch.Split() // first split
	tmch.Pause() // see TestPause()
	tmch.Split() // split while paused
	assert.True(t, tmch.Timer.Paused() && !tmch.SplitHandler.IsFinished(),
		"Splits while paused should do nothing")

	tmch.Pause() // unpause; see TestPause()
	tmch.Stop()  // cancel; see TestStop()
	tmch.Split() // split while stopped
	assert.True(t, tmch.Timer.Stopped() && !tmch.SplitHandler.IsFinished(),
		"Splits while cancelled should do nothing")

	// TODO: maybe test state closure
}

func TestPause(t *testing.T) {
	tmch := TimeMachine{Timer: &timer.Timer{}, SplitHandler: &splitter.SplitHandler{}}

	// TODO: fake this
	tmch.SplitHandler.SetSplits([]splitter.Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	// note: tmch is idle by default
	tmch.Pause() // pause while idle
	assert.True(t, tmch.Timer.Idle() && tmch.SplitHandler.IsBeforeFirstSplit(),
		"Pauses while idle should do nothing")

	tmch.Split() // timer start; see TestSplit()
	tmch.Pause() // pause while active
	assert.True(t, tmch.Timer.Paused() && !tmch.SplitHandler.IsFinished(),
		"Pause inputs while active should pause")

	tmch.Pause() // pause input while paused
	assert.True(t, tmch.Timer.Running() && !tmch.SplitHandler.IsFinished(),
		"Pause inputs while paused should unpause")

	tmch.Stop()  // cancel; see TestStop()
	tmch.Pause() // pause while cancelled
	assert.True(t, tmch.Timer.Stopped() && !tmch.SplitHandler.IsFinished(),
		"Pause inputs while cancelled should do nothing")

	tmch.Stop() // restart; see TestStop()
	tmch.Split()
	tmch.Split()
	tmch.Split() // finish; see TestSplit()
	tmch.Pause() // pause while finished
	assert.True(t, tmch.Timer.Stopped() && tmch.SplitHandler.IsFinished(),
		"Pause inputs while finished should do nothing")
}

func TestStop(t *testing.T) {
	tmch := TimeMachine{Timer: &timer.Timer{}, SplitHandler: &splitter.SplitHandler{}}

	// TODO: fake this
	tmch.SplitHandler.SetSplits([]splitter.Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	// note: tmch is idle by default
	tmch.Stop() // Stop input while idle
	assert.True(t, tmch.Timer.Idle() && tmch.SplitHandler.IsBeforeFirstSplit(),
		"Stop inputs while idle should do nothing")

	tmch.Split() // start; see TestSplit()
	tmch.Stop()  // Stop input while running
	assert.True(t, tmch.Timer.Stopped() && !tmch.SplitHandler.IsFinished(),
		"Stop inputs while running should cancel the run")

	tmch.Stop() // Stop input while cancelled
	assert.True(t, tmch.Timer.Idle() && tmch.SplitHandler.IsBeforeFirstSplit(),
		"Stop inputs while cancelled should revert to idle")

	tmch.Split() // start
	tmch.Pause() // pause; see TestPause()
	tmch.Stop()  // Stop input while paused
	assert.True(t, tmch.Timer.Stopped() && !tmch.SplitHandler.IsFinished(),
		"Stop inputs while paused should still cancel the run")

	tmch.Stop()  // restart
	tmch.Split() // start
	tmch.Split()
	tmch.Split() // finish
	tmch.Stop()  // Stop input while finished
	assert.True(t, tmch.Timer.Idle() && tmch.SplitHandler.IsBeforeFirstSplit(),
		"Stop inputs while finished should revert to idle")
}

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
	machine := TimeMachine{Timer: &timer.Timer{}, SplitHandler: &splitter.SplitHandler{}}

	// TODO: fake this
	machine.SplitHandler.SetSplits([]splitter.Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	assert.True(t, machine.Timer.Idle(),
		"Timer should be idle before first split")

	machine.Split() // timer start
	assert.True(t, machine.Timer.Running(),
		"Timer should be running after first split input")
	assert.True(t, machine.SplitHandler.IsBeforeFirstSplit(),
		"SplitHandler should not split on the first split input")

	machine.Split() // first split
	assert.True(t, machine.Timer.Running(),
		"Timer should still be running after first real split")
	assert.False(t, machine.SplitHandler.IsBeforeFirstSplit() || machine.SplitHandler.IsFinished(),
		"SplitHandler should still be active after first real split")

	machine.Split() // final split
	assert.True(t, machine.Timer.Stopped(),
		"Timer should stop after final split")
	assert.True(t, machine.SplitHandler.IsFinished(),
		"SplitHandler should be finished after final split")

	machine.Split() // additional inputs after final
	assert.True(t, machine.Timer.Stopped() && machine.SplitHandler.IsFinished(),
		"Splits while stopped should do nothing")

	machine.Stop()  // restart; see TestStop()
	machine.Split() // start
	machine.Split() // first split
	machine.Pause() // see TestPause()
	machine.Split() // split while paused
	assert.True(t, machine.Timer.Paused() && !machine.SplitHandler.IsFinished(),
		"Splits while paused should do nothing")

	machine.Pause() // unpause; see TestPause()
	machine.Stop()  // cancel; see TestStop()
	machine.Split() // split while stopped
	assert.True(t, machine.Timer.Stopped() && !machine.SplitHandler.IsFinished(),
		"Splits while cancelled should do nothing")

	// TODO: maybe test state closure
}

func TestPause(t *testing.T) {
	machine := TimeMachine{Timer: &timer.Timer{}, SplitHandler: &splitter.SplitHandler{}}

	// TODO: fake this
	machine.SplitHandler.SetSplits([]splitter.Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	// note: tmch is idle by default
	machine.Pause() // pause while idle
	assert.True(t, machine.Timer.Idle() && machine.SplitHandler.IsBeforeFirstSplit(),
		"Pauses while idle should do nothing")

	machine.Split() // timer start; see TestSplit()
	machine.Pause() // pause while active
	assert.True(t, machine.Timer.Paused() && !machine.SplitHandler.IsFinished(),
		"Pause inputs while active should pause")

	machine.Pause() // pause input while paused
	assert.True(t, machine.Timer.Running() && !machine.SplitHandler.IsFinished(),
		"Pause inputs while paused should unpause")

	machine.Stop()  // cancel; see TestStop()
	machine.Pause() // pause while cancelled
	assert.True(t, machine.Timer.Stopped() && !machine.SplitHandler.IsFinished(),
		"Pause inputs while cancelled should do nothing")

	machine.Stop() // restart; see TestStop()
	machine.Split()
	machine.Split()
	machine.Split() // finish; see TestSplit()
	machine.Pause() // pause while finished
	assert.True(t, machine.Timer.Stopped() && machine.SplitHandler.IsFinished(),
		"Pause inputs while finished should do nothing")
}

func TestStop(t *testing.T) {
	machine := TimeMachine{Timer: &timer.Timer{}, SplitHandler: &splitter.SplitHandler{}}

	// TODO: fake this
	machine.SplitHandler.SetSplits([]splitter.Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	// note: tmch is idle by default
	machine.Stop() // Stop input while idle
	assert.True(t, machine.Timer.Idle() && machine.SplitHandler.IsBeforeFirstSplit(),
		"Stop inputs while idle should do nothing")

	machine.Split() // start; see TestSplit()
	machine.Stop()  // Stop input while running
	assert.True(t, machine.Timer.Stopped() && !machine.SplitHandler.IsFinished(),
		"Stop inputs while running should cancel the run")

	machine.Stop() // Stop input while cancelled
	assert.True(t, machine.Timer.Idle() && machine.SplitHandler.IsBeforeFirstSplit(),
		"Stop inputs while cancelled should revert to idle")

	machine.Split() // start
	machine.Pause() // pause; see TestPause()
	machine.Stop()  // Stop input while paused
	assert.True(t, machine.Timer.Stopped() && !machine.SplitHandler.IsFinished(),
		"Stop inputs while paused should still cancel the run")

	machine.Stop()  // restart
	machine.Split() // start
	machine.Split()
	machine.Split() // finish
	machine.Stop()  // Stop input while finished
	assert.True(t, machine.Timer.Idle() && machine.SplitHandler.IsBeforeFirstSplit(),
		"Stop inputs while finished should revert to idle")
}

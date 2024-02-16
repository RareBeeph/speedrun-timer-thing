package timing

import (
	"speedruntimer/timing/splitter"
	"speedruntimer/timing/timer"
	"time"
)

type TimeMachine struct {
	Timer        timer.Timer
	SplitHandler *splitter.SplitHandler
}

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

func (t *TimeMachine) Split() {
	if t.Timer.Idle() {
		t.Timer.Start()
		return
	}

	if t.Timer.Paused() || t.Timer.Stopped() {
		return
	}

	// TODO: magic number equal to time.Millisecond
	t.SplitHandler.Split(time.Duration(t.Timer.Milliseconds() * time.Millisecond))
	if t.SplitHandler.IsFinished() {
		t.Timer.Stop()
	}
}

func (t *TimeMachine) Pause() {
	if t.Timer.Running() {
		t.Timer.Pause()
	} else if t.Timer.Paused() {
		t.Timer.Resume()
	}
}

func (t *TimeMachine) Stop() {
	if t.Timer.Stopped() {
		t.Timer.Restart()
		t.SplitHandler.Restart()
	} else if t.Timer.Running() || t.Timer.Paused() {
		t.Timer.Stop()
	}
}

package splitter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
	Recap:

					Split()			Restart()
				  -----------------------------
	IsIdle        | IsActive_1    |	IsIdle
	...			  | ...           | ...
	IsActive_n	  |	IsActive_(n+1)|	IsIdle
	...			  | ...           | ...
	IsActive_(L-1)|	IsFinished	  |	IsIdle
	IsFinished	  |	IsFinished	  |	IsIdle


	The number of splits must only be updated while idle,
	since it determines the number of non-idle states
*/

// Split() updates the selected Split and Labels according to the current duration since the timer started,
// then increments the cursor to select the next Split and its corresponding Labels.
// If all splits have been exhausted (and so none is selected), it does nothing.

func TestSplit(t *testing.T) {
	tsh := &SplitHandler{}

	// TODO: fake this
	tsh.SetSplits([]Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	// TODO: fake cursor value
	cur := 0
	tsh.cursor = cur

	// TODO: fake the input durations
	tsh.Split(time.Duration(166000000000)) // not a best segment

	assert.True(t, tsh.cursor == cur+1, "cursor increments when the handler is not finished")
	assert.True(t, tsh.splits[cur].ActiveRunTime == time.Duration(166000000000), "Selected Split's ActiveRunTime is updated with the given time")
	assert.True(t, tsh.splits[cur].BestSegment == time.Duration(153983000000), "Selected Split's BestSegment is not updated if run is not a best segment")
	assert.False(t, tsh.IsFinished(), "IsFinished returns false when the handler is not finished")

	tsh.Split(time.Duration(299000000000)) // is a best segment

	// difference of given splits
	assert.True(t, tsh.splits[cur+1].BestSegment == time.Duration(133000000000), "Selected Split's BestSegment is updated if run is a best segment")
	assert.True(t, tsh.IsFinished(), "IsFinished returns true when the handler is finished")

	tsh.Split(time.Duration(355000000000)) // list has already been exhausted

	assert.True(t, tsh.cursor == cur+2, "cursor does not increment when the handler is finished")
}

// Restart() resets all Splits to their default state,
// updating all of their PBTime fields if the final split is better than its stored pb time.
// It then updates all Labels to correspond to match the new states of their respective Splits,
// and resets the cursor to select the first Split.

func TestRestart(t *testing.T) {
	tsh := &SplitHandler{}

	// TODO: fake this
	tsh.SetSplits([]Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	// TODO: fake these
	tsh.Split(time.Duration(166000000000))
	tsh.Split(time.Duration(299000000000)) // run is complete, is a pb
	tsh.Restart()

	// TODO: make sure this applies to arbitrary index
	assert.True(t, tsh.splits[0].PBTime == time.Duration(166000000000), "PBTimes are updated on PB")
	assert.True(t, tsh.cursor == 0, "Cursor is reset")

	tsh.Split(time.Duration(155000000000)) // is green
	tsh.Split(time.Duration(311000000000)) // but is not pb
	tsh.Restart()

	assert.True(t, tsh.splits[0].PBTime == time.Duration(166000000000), "PBTimes are not updated on non-PB, even if a non-final split is green")

	tsh.Split(time.Duration(155000000000))
	tsh.Restart() // is incomplete

	assert.True(t, tsh.splits[0].PBTime == time.Duration(166000000000), "PBTimes are not updated on incomplete run, even if a split is green")
}

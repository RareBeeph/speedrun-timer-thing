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

func TestSetSplits(t *testing.T) {
	tsh := &SplitHandler{}

	// TODO: maybe fake input list length
	fakesplits := []Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	}
	tsh.SetSplits([]Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	assert.ElementsMatch(t, tsh.splits, fakesplits,
		"Splits should be set as given")
	assert.Len(t, tsh.SplitLabels, len(tsh.splits),
		"Split label list should be as long as split list")
	assert.Equal(t, tsh.SplitLabels[0].Text, tsh.splits[0].String(),
		"Split label should contain string returned by split String()")
	assert.Len(t, tsh.DeltaLabels, len(tsh.splits),
		"Delta label list should be as long as split list")
	assert.Equal(t, tsh.DeltaLabels[0].Text, tsh.splits[0].Delta(),
		"Delta label should contain string returned by split Delta()")

	fakesplits = []Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
	} // shorter than before
	tsh.SetSplits(fakesplits)

	assert.ElementsMatch(t, tsh.splits, fakesplits,
		"Splits shouldn't leave residue when overwriting")
	assert.Len(t, tsh.SplitLabels, len(tsh.splits),
		"Label lists should be able to shorten to be only as long as split list")
}

// Split() updates the selected Split and Labels according to the current duration since the timer started,
// then increments the cursor to select the next Split and its corresponding Labels.
// If all splits have been exhausted (and so none is selected), it does nothing.

func TestHandlerSplit(t *testing.T) {
	tsh := &SplitHandler{}

	// TODO: fake the input durations
	initialBestSegment := time.Duration(153983000000) // initial best segment for split idx 0
	badFirstSplit := time.Duration(166000000000)      // greater than initialBestSegment
	goodSecondSplit := time.Duration(299000000000)    // greater than badFirstSplit by less than 398000000000 (arbitrary given BestSegment for split idx 1)

	tsh.SetSplits([]Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: initialBestSegment},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	// TODO: fake cursor value
	cur := 0
	tsh.cursor = cur

	tsh.Split(badFirstSplit) // not a best segment

	assert.Equal(t, tsh.cursor, cur+1,
		"Cursor should increment when the handler is not finished")
	assert.Equal(t, tsh.splits[cur].ActiveRunTime, badFirstSplit,
		"Selected Split's ActiveRunTime should be updated with the given time")
	assert.Equal(t, tsh.splits[cur].BestSegment, initialBestSegment,
		"Selected Split's BestSegment should not be updated if segment in run is not a best segment")
	assert.False(t, tsh.IsFinished(),
		"IsFinished should return false when the handler is not finished")
	assert.Equal(t, tsh.SplitLabels[cur].Text, tsh.splits[cur].String(),
		"Split labels should be updated on split")
	assert.Equal(t, tsh.DeltaLabels[cur].Text, tsh.splits[cur].Delta(),
		"Delta labels should be updated on split")

	cur = tsh.cursor           // update knowledge
	tsh.Split(goodSecondSplit) // is a best segment

	// difference of given splits
	assert.Equal(t, tsh.splits[cur].BestSegment, goodSecondSplit-badFirstSplit,
		"Selected Split's BestSegment should be updated if run is a best segment")
	assert.True(t, tsh.IsFinished(),
		"IsFinished should return true when the handler is finished")

	cur = tsh.cursor                       // update knowledge
	tsh.Split(time.Duration(355000000000)) // list has already been exhausted

	assert.Equal(t, tsh.cursor, cur,
		"Cursor should not increment when the handler is finished")
}

// Restart() resets all Splits to their default state,
// updating all of their PBTime fields if the final split is better than its stored pb time.
// It then updates all Labels to correspond to match the new states of their respective Splits,
// and resets the cursor to select the first Split.

func TestHandlerRestart(t *testing.T) {
	tsh := &SplitHandler{}

	// TODO: fake these
	badFirstSplit := time.Duration(166000000000)   // greater than 154500000000 (arbitrary given first split)
	goodSecondSplit := time.Duration(299000000000) // less than 400000000000 (arbitrary given last split)

	goodFirstSplit := time.Duration(155000000000) // less than badFirstSplit
	badSecondSplit := time.Duration(311000000000) // greater than goodSecondSplit

	tsh.SetSplits([]Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	tsh.Split(badFirstSplit)
	tsh.Split(goodSecondSplit) // run is complete, is a pb
	tsh.Restart()

	// TODO: make sure this applies to arbitrary index
	assert.Equal(t, tsh.splits[0].PBTime, badFirstSplit,
		"PBTimes should be updated on PB, even if a non-final split isn't green")
	assert.Equal(t, tsh.SplitLabels[0].Text, tsh.splits[0].String(),
		"Split labels should be updated on PB")
	assert.Equal(t, tsh.DeltaLabels[0].Text, tsh.splits[0].Delta(),
		"Delta labels should be updated on PB")
	assert.Zero(t, tsh.cursor, "Cursor should be reset")

	tsh.Split(goodFirstSplit) // is green
	tsh.Split(badSecondSplit) // but is not pb
	tsh.Restart()

	assert.Equal(t, tsh.splits[0].PBTime, badFirstSplit,
		"PBTimes should not be updated on non-PB, even if a non-final split is green")
	assert.Equal(t, tsh.SplitLabels[0].Text, tsh.splits[0].String(),
		"Split labels should still be updated on non-PB (return from active run to previous PB text)")
	assert.Equal(t, tsh.DeltaLabels[0].Text, tsh.splits[0].Delta(),
		"Delta labels should be updated on non-PB")

	tsh.Split(goodFirstSplit)
	tsh.Restart() // is incomplete
	assert.Equal(t, tsh.splits[0].PBTime, badFirstSplit,
		"PBTimes should not be updated on incomplete run, even if a split is green")
	// labels should still update, but that's redundant

	tsh.Restart() // from idle
	// not sure what to test about this besides that it doesn't crash or anything.
}

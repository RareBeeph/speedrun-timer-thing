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
	handler := &SplitHandler{}

	// TODO: maybe fake input list length
	splits := []Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	}
	handler.SetSplits([]Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	assert.ElementsMatch(t, handler.splits, splits,
		"Splits should be set as given")
	assert.Len(t, handler.SplitLabels, len(handler.splits),
		"Split label list should be as long as split list")
	assert.Equal(t, handler.SplitLabels[0].Text, handler.splits[0].String(),
		"Split label should contain string returned by split String()")
	assert.Len(t, handler.DeltaLabels, len(handler.splits),
		"Delta label list should be as long as split list")
	assert.Equal(t, handler.DeltaLabels[0].Text, handler.splits[0].Delta(),
		"Delta label should contain string returned by split Delta()")

	splits = []Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
	} // shorter than before
	handler.SetSplits(splits)

	assert.ElementsMatch(t, handler.splits, splits,
		"Splits shouldn't leave residue when overwriting")
	assert.Len(t, handler.SplitLabels, len(handler.splits),
		"Label lists should be able to shorten to be only as long as split list")
}

// Split() updates the selected Split and Labels according to the current duration since the timer started,
// then increments the cursor to select the next Split and its corresponding Labels.
// If all splits have been exhausted (and so none is selected), it does nothing.

func TestHandlerSplit(t *testing.T) {
	handler := &SplitHandler{}

	// TODO: fake the input durations
	initialBestSegment := time.Duration(153983000000) // initial best segment for split idx 0
	badFirstSplit := time.Duration(166000000000)      // greater than initialBestSegment
	goodSecondSplit := time.Duration(299000000000)    // greater than badFirstSplit by less than 398000000000 (arbitrary given BestSegment for split idx 1)

	handler.SetSplits([]Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: initialBestSegment},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	// TODO: fake cursor value
	cursor := 0

	handler.Split(badFirstSplit) // not a best segment

	assert.Equal(t, handler.cursor, cursor+1,
		"Cursor should increment when the handler is not finished")
	assert.Equal(t, handler.splits[cursor].ActiveRunTime, badFirstSplit,
		"Selected Split's ActiveRunTime should be updated with the given time")
	assert.Equal(t, handler.splits[cursor].BestSegment, initialBestSegment,
		"Selected Split's BestSegment should not be updated if segment in run is not a best segment")
	assert.False(t, handler.IsFinished(),
		"IsFinished should return false when the handler is not finished")
	assert.Equal(t, handler.SplitLabels[cursor].Text, handler.splits[cursor].String(),
		"Split labels should be updated on split")
	assert.Equal(t, handler.DeltaLabels[cursor].Text, handler.splits[cursor].Delta(),
		"Delta labels should be updated on split")

	cursor = handler.cursor        // update knowledge
	handler.Split(goodSecondSplit) // is a best segment

	// difference of given splits
	assert.Equal(t, handler.splits[cursor].BestSegment, goodSecondSplit-badFirstSplit,
		"Selected Split's BestSegment should be updated if run is a best segment")
	assert.True(t, handler.IsFinished(),
		"IsFinished should return true when the handler is finished")

	cursor = handler.cursor                    // update knowledge
	handler.Split(time.Duration(355000000000)) // list has already been exhausted

	assert.Equal(t, handler.cursor, cursor,
		"Cursor should not increment when the handler is finished")
}

// Restart() resets all Splits to their default state,
// updating all of their PBTime fields if the final split is better than its stored pb time.
// It then updates all Labels to correspond to match the new states of their respective Splits,
// and resets the cursor to select the first Split.

func TestHandlerRestart(t *testing.T) {
	handler := &SplitHandler{}

	// TODO: fake these
	badFirstSplit := time.Duration(166000000000)   // greater than 154500000000 (arbitrary given first split)
	goodSecondSplit := time.Duration(299000000000) // less than 400000000000 (arbitrary given last split)

	goodFirstSplit := time.Duration(155000000000) // less than badFirstSplit
	badSecondSplit := time.Duration(311000000000) // greater than goodSecondSplit

	handler.SetSplits([]Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	handler.Split(badFirstSplit)
	handler.Split(goodSecondSplit) // run is complete, is a pb
	handler.Restart()

	// TODO: make sure this applies to arbitrary index
	assert.Equal(t, handler.splits[0].PBTime, badFirstSplit,
		"PBTimes should be updated on PB, even if a non-final split isn't green")
	assert.Equal(t, handler.SplitLabels[0].Text, handler.splits[0].String(),
		"Split labels should be updated on PB")
	assert.Equal(t, handler.DeltaLabels[0].Text, handler.splits[0].Delta(),
		"Delta labels should be updated on PB")
	assert.Zero(t, handler.cursor, "Cursor should be reset")

	handler.Split(goodFirstSplit) // is green
	handler.Split(badSecondSplit) // but is not pb
	handler.Restart()

	assert.Equal(t, handler.splits[0].PBTime, badFirstSplit,
		"PBTimes should not be updated on non-PB, even if a non-final split is green")
	assert.Equal(t, handler.SplitLabels[0].Text, handler.splits[0].String(),
		"Split labels should still be updated on non-PB (return from active run to previous PB text)")
	assert.Equal(t, handler.DeltaLabels[0].Text, handler.splits[0].Delta(),
		"Delta labels should be updated on non-PB")

	handler.Split(goodFirstSplit)
	handler.Restart() // is incomplete
	assert.Equal(t, handler.splits[0].PBTime, badFirstSplit,
		"PBTimes should not be updated on incomplete run, even if a split is green")
	// labels should still update, but that's redundant

	handler.Restart() // from idle
	// not sure what to test about this besides that it doesn't crash or anything.
}

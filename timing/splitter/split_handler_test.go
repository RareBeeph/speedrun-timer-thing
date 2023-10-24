package splitter

import (
	"log"
	"math/rand"
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

func fakeSplits() (out []Split) {
	// not sure if i should do this seeding
	rand.Seed(time.Now().Unix())

	l := rand.Intn(10) // l == the last index of the output. so len(out) will be l+1, ranging from 1 to 11.
	for l >= 0 {
		t := time.Duration(0)
		if len(out) > 0 {
			t = out[len(out)-1].PBTime
		}
		s := Split{
			// durations capped to 24 hours to prevent integer overflow (would only occur in practice if program runs for ~300 years)
			PBTime:      t + time.Duration(rand.Int63n(int64(time.Hour*24))), // monotonically increasing
			BestSegment: time.Duration(rand.Int63n(int64(time.Hour * 24))),
		}
		out = append(out, s)
		l--
	}
	return out
}

func TestSetSplits(t *testing.T) {
	handler := &SplitHandler{}

	fakesplits := fakeSplits()
	handler.SetSplits(fakesplits)

	assert.ElementsMatch(t, handler.splits, fakesplits,
		"Splits should be set as given")
	assert.Len(t, handler.SplitLabels, len(handler.splits),
		"Split label list should be as long as split list")
	assert.Len(t, handler.DeltaLabels, len(handler.splits),
		"Delta label list should be as long as split list")
	for i := range handler.splits {
		assert.Equal(t, handler.SplitLabels[i].Text, handler.splits[i].String(),
			"Split label should contain string returned by split String()")
		assert.Equal(t, handler.DeltaLabels[i].Text, handler.splits[i].Delta(),
			"Delta label should contain string returned by split Delta()")
	}

	log.Println(len(fakesplits))
	fakesplits = fakesplits[1:] // shorter than before; possibly len == 0
	log.Println(len(fakesplits))
	handler.SetSplits(fakesplits)

	assert.ElementsMatch(t, handler.splits, fakesplits,
		"Splits shouldn't leave residue when overwriting")
	assert.Len(t, handler.SplitLabels, len(handler.splits),
		"Label lists should be able to shorten to be only as long as split list")
}

// Split() updates the selected Split and Labels according to the current duration since the timer started,
// then increments the cursor to select the next Split and its corresponding Labels.
// If all splits have been exhausted (and so none is selected), it does nothing.
func TestHandlerSplit(t *testing.T) {
	handler := &SplitHandler{}

	fakeStoredSplits := fakeSplits()
	// hacky run faking method; encodes active run time as pb time
	fakeRunTimes := fakeSplits()

	handler.SetSplits(fakeStoredSplits)

	for i := range fakeRunTimes {
		var initialBestSegment time.Duration
		if i < len(handler.splits) {
			initialBestSegment = handler.splits[i].BestSegment
		} else {
			initialBestSegment = handler.splits[len(handler.splits)-1].BestSegment
		}

		var newSegment time.Duration
		if i == 0 {
			newSegment = fakeRunTimes[0].PBTime
		} else {
			newSegment = fakeRunTimes[i].PBTime - fakeRunTimes[i-1].PBTime
		}

		handler.Split(fakeRunTimes[i].PBTime)

		if i < len(handler.splits) {
			assert.Equal(t, handler.splits[i].ActiveRunTime, fakeRunTimes[i].PBTime,
				"Selected Split's ActiveRunTime should be updated with the given time")
			assert.Equal(t, handler.SplitLabels[i].Text, handler.splits[i].String(),
				"Split labels should be updated on split")
			assert.Equal(t, handler.DeltaLabels[i].Text, handler.splits[i].Delta(),
				"Delta labels should be updated on split")

			if newSegment < initialBestSegment {
				assert.Equal(t, handler.splits[i].BestSegment, newSegment,
					"Selected Split's BestSegment should be updated if run is a best segment")
			} else {
				assert.Equal(t, handler.splits[i].BestSegment, initialBestSegment,
					"Selected Split's BestSegment should not be updated if segment in run is not a best segment")
			}
		} else {
			assert.Equal(t, initialBestSegment, handler.splits[len(handler.splits)-1].BestSegment,
				"Final split's BestSegment should not be updated on splits performed after the end of the split list")
		}

		if handler.cursor < len(handler.splits) {
			assert.False(t, handler.IsFinished(),
				"IsFinished should return false when the handler is not finished")
			assert.Equal(t, handler.cursor, i+1,
				"Cursor should increment when the handler is not finished")
		} else {
			assert.True(t, handler.IsFinished(),
				"IsFinished should return true when the handler is finished")
			assert.Equal(t, len(handler.splits), handler.cursor,
				"Cursor should not increment when the handler is finished")
		}
	}
}

// Restart() resets all Splits to their default state,
// updating all of their PBTime fields if the final split is better than its stored pb time.
// It then updates all Labels to correspond to match the new states of their respective Splits,
// and resets the cursor to select the first Split.
func TestHandlerRestart(t *testing.T) {
	handler := &SplitHandler{}

	fakeStoredSplits := fakeSplits()
	// a little bit of a hacky run faking method:
	fakeRunTimes := fakeSplits() // can represent an incomplete run or an over-complete run, which covers all my bases

	handler.SetSplits(fakeStoredSplits)

	for i := range fakeRunTimes {
		handler.Split(fakeRunTimes[i].PBTime)
	}
	handler.Restart()

	if len(fakeRunTimes) < len(handler.splits) {
		for i := range handler.splits {
			assert.Equal(t, handler.splits[i].PBTime, fakeStoredSplits[i].PBTime,
				"PBTimes should not be updated on incomplete run")
		}
	} else if fakeRunTimes[len(handler.splits)-1].PBTime <= handler.splits[len(handler.splits)-1].PBTime {
		for i := range handler.splits {
			assert.Equal(t, handler.splits[i].PBTime, fakeRunTimes[i].PBTime,
				"PBTimes should be updated on PB")
		}
	} else {
		for i := range handler.splits {
			assert.Equal(t, handler.splits[i].PBTime, fakeStoredSplits[i].PBTime,
				"PBTimes should not be updated on non-PB")
		}
	}

	for i := range handler.splits {
		assert.Equal(t, handler.SplitLabels[i].Text, handler.splits[i].String(),
			"Split labels should be updated")
		assert.Equal(t, handler.DeltaLabels[i].Text, handler.splits[i].Delta(),
			"Delta labels should be updated")
	}
	assert.Zero(t, handler.cursor, "Cursor should be reset")

	handler.Restart() // from idle
	// not sure what to test about this besides that it doesn't crash or anything.
}

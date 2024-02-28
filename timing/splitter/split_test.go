package splitter

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// returns a random duration between 0 and max
func randDurationWithMax(max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max)))
}

const day = time.Hour * 24

func TestSplit(t *testing.T) {
	split := Split{Name: "Fake Split 1", bestSegment: randDurationWithMax(day)}

	baseTime := split.bestSegment + randDurationWithMax(day)
	bestSegmentEndTime := randDurationWithMax(day)
	bestSegmentStartTime := bestSegmentEndTime - randDurationWithMax(split.bestSegment)

	split.Split(baseTime, time.Duration(0)) // may be green, but not a best segment, and not resetting to update PBTime
	assert.Equal(t, split.ActiveRunTime, baseTime,
		"Split() should set ActiveRunTime")

	split.ActiveRunTime = time.Duration(0)
	split.Split(bestSegmentEndTime, bestSegmentStartTime) // Best segment, but may be not green
	assert.Equal(t, split.bestSegment, bestSegmentEndTime-bestSegmentStartTime,
		"Split() should set BestSegment on a best segment, even if not green")
}

func TestRestart(t *testing.T) {
	initialPBTime := randDurationWithMax(day)
	// TODO: perhaps just have a single split time, and let running the test a lot hit these two cases roughly equally often
	goodTime := randDurationWithMax(initialPBTime)
	badTime := initialPBTime + randDurationWithMax(day)

	split := Split{Name: "Fake Split 1", pbTime: initialPBTime}

	split.ActiveRunTime = goodTime // is green
	split.Restart(false)           // but not a PB run
	assert.Zero(t, split.ActiveRunTime,
		"Restart() should reset ActiveRunTime to 0")
	assert.Equal(t, split.pbTime, initialPBTime,
		"Restart() should not set PBTime if isPB is false, even if green")

	split.ActiveRunTime = badTime // is not green
	split.Restart(true)           // but is a PB run
	assert.Equal(t, split.pbTime, badTime,
		"Restart() should set PBTime if isPB is true, even if not green")
}

func TestIsGreen(t *testing.T) {
	initialPBTime := randDurationWithMax(day)

	split := Split{Name: "Fake Split 1", pbTime: initialPBTime}

	// TODO: perhaps let running the test a lot hit these two cases roughly equally often

	split.ActiveRunTime = initialPBTime + randDurationWithMax(day) // is not green
	assert.False(t, split.IsGreen(),
		"IsGreen() should return false when active run is not ahead of PB")

	split.ActiveRunTime = randDurationWithMax(initialPBTime) // is green
	assert.True(t, split.IsGreen(),
		"IsGreen() should return true when active run is ahead of PB")
}

func TestDisplayTime(t *testing.T) {
	initialPBTime := randDurationWithMax(day)
	runTime := randDurationWithMax(day)

	split := Split{Name: "Fake Split 1", pbTime: initialPBTime}

	assert.Equal(t, split.DisplayTime(), initialPBTime,
		"DisplayTime() should return PBTime when no run is active")

	split.ActiveRunTime = runTime // is active but not green
	assert.Equal(t, split.DisplayTime(), runTime,
		"DisplayTime() should return ActiveRunTime when a run is active, even if slower than PBTime")
}

func TestDelta(t *testing.T) {
	split := Split{Name: "Fake Split 1", pbTime: randDurationWithMax(day)}

	assert.Zero(t, split.Delta(),
		"Delta() should return the empty string when no run is active")

	// TODO: where is the second half of this
}

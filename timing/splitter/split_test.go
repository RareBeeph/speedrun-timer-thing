package splitter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	split := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)}

	// TODO: fake this
	baseTime := time.Duration(154400000000)
	bestSegmentEndTime := time.Duration(180000000000)  // more than goodNotBestTime
	bestSegmentStartTime := time.Duration(40000000000) // less than end time by less than 153983000000 (arbitrary given BestSegment)

	split.Split(baseTime, time.Duration(0)) // Green, but not a best segment, and not resetting to update PBTime
	assert.Equal(t, split.ActiveRunTime, baseTime,
		"Split() should set ActiveRunTime")

	split.ActiveRunTime = time.Duration(0)
	split.Split(bestSegmentEndTime, bestSegmentStartTime) // Best segment, but not green
	assert.Equal(t, split.BestSegment, bestSegmentEndTime-bestSegmentStartTime,
		"Split() should set BestSegment on a best segment, even if not green")
}

func TestRestart(t *testing.T) {
	// TODO: fake this
	initialPBTime := time.Duration(154500000000)
	goodTime := time.Duration(154400000000) // Less than initialPBTime
	badTime := time.Duration(180000000000)  // More than initialPBTime

	split := Split{Name: "Fake Split 1", PBTime: initialPBTime}

	split.ActiveRunTime = goodTime // is green
	split.Restart(false)           // but not a PB run
	assert.Zero(t, split.ActiveRunTime,
		"Restart() should reset ActiveRunTime to 0")
	assert.Equal(t, split.PBTime, initialPBTime,
		"Restart() should not set PBTime if isPB is false, even if green")

	split.ActiveRunTime = badTime // is not green
	split.Restart(true)           // but is a PB run
	assert.Equal(t, split.PBTime, badTime,
		"Restart() should set PBTime if isPB is true, even if not green")
}

func TestIsGreen(t *testing.T) {
	// TODO: fake this
	split := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)}

	split.ActiveRunTime = time.Duration(180000000000) // is not green
	assert.False(t, split.IsGreen(),
		"IsGreen() should return false when active run is not ahead of PB")

	split.ActiveRunTime = time.Duration(154400000000) // is green
	assert.True(t, split.IsGreen(),
		"IsGreen() should return true when active run is ahead of PB")
}

func TestDisplayTime(t *testing.T) {
	// TODO: fake this
	initialPBTime := time.Duration(154500000000)
	runTime := time.Duration(180000000000) // greater than initialPBTime

	split := Split{Name: "Fake Split 1", PBTime: initialPBTime}

	assert.Equal(t, split.DisplayTime(), initialPBTime,
		"DisplayTime() should return PBTime when no run is active")

	split.ActiveRunTime = runTime // is active but not green
	assert.Equal(t, split.DisplayTime(), runTime,
		"DisplayTime() should return ActiveRunTime when a run is active, even if slower than PBTime")
}

func TestDelta(t *testing.T) {
	// TODO: fake this
	split := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000)}

	assert.Zero(t, split.Delta(),
		"Delta() should return the empty string when no run is active")
}

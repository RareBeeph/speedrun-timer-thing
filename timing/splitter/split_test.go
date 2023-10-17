package splitter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	spl := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)}

	// TODO: fake this
	goodNotBestTime := time.Duration(154400000000)
	bestSegmentEndTime := time.Duration(180000000000)  // more than goodNotBestTime
	bestSegmentStartTime := time.Duration(40000000000) // less than end time by less than 153983000000 (arbitrary given BestSegment)

	spl.Split(goodNotBestTime, time.Duration(0)) // Green, but not a best segment, and not resetting to update PBTime
	assert.Equal(t, spl.ActiveRunTime, goodNotBestTime,
		"Split() should set ActiveRunTime")

	spl.ActiveRunTime = time.Duration(0)
	spl.Split(bestSegmentEndTime, bestSegmentStartTime) // Best segment, but not green
	assert.Equal(t, spl.BestSegment, bestSegmentEndTime-bestSegmentStartTime,
		"Split() should set BestSegment on a best segment, even if not green")
}

func TestRestart(t *testing.T) {
	// TODO: fake this
	initialPBTime := time.Duration(154500000000)
	goodTime := time.Duration(154400000000) // Less than initialPBTime
	badTime := time.Duration(180000000000)  // More than initialPBTime

	spl := Split{Name: "Fake Split 1", PBTime: initialPBTime}

	spl.ActiveRunTime = goodTime // is green
	spl.Restart(false)           // but not a PB run
	assert.Zero(t, spl.ActiveRunTime,
		"Restart() should reset ActiveRunTime to 0")
	assert.Equal(t, spl.PBTime, initialPBTime,
		"Restart() should not set PBTime if isPB is false, even if green")

	spl.ActiveRunTime = badTime // is not green
	spl.Restart(true)           // but is a PB run
	assert.Equal(t, spl.PBTime, badTime,
		"Restart() should set PBTime if isPB is true, even if not green")
}

func TestIsGreen(t *testing.T) {
	// TODO: fake this
	spl := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)}

	spl.ActiveRunTime = time.Duration(180000000000) // is not green
	assert.False(t, spl.IsGreen(),
		"IsGreen() should return false when active run is not ahead of PB")

	spl.ActiveRunTime = time.Duration(154400000000) // is green
	assert.True(t, spl.IsGreen(),
		"IsGreen() should return true when active run is ahead of PB")
}

func TestDisplayTime(t *testing.T) {
	// TODO: fake this
	initialPBTime := time.Duration(154500000000)
	runTime := time.Duration(180000000000) // greater than initialPBTime

	spl := Split{Name: "Fake Split 1", PBTime: initialPBTime}

	assert.Equal(t, spl.DisplayTime(), initialPBTime,
		"DisplayTime() should return PBTime when no run is active")

	spl.ActiveRunTime = runTime // is active but not green
	assert.Equal(t, spl.DisplayTime(), runTime,
		"DisplayTime() should return ActiveRunTime when a run is active, even if slower than PBTime")
}

func TestDelta(t *testing.T) {
	// TODO: fake this
	spl := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000)}

	assert.Zero(t, spl.Delta(),
		"Delta() should return the empty string when no run is active")
}

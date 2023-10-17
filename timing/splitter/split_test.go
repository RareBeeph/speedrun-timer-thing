package splitter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	spl := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)}

	// TODO: fake this
	t1 := time.Duration(154400000000)
	t2 := time.Duration(180000000000) // more than t1
	t3 := time.Duration(40000000000)  // less than t2 by less than 153983000000 (arbitrary given BestSegment)

	spl.Split(t1, time.Duration(0)) // Green, but not a best segment, and not resetting to update PBTime
	assert.Equal(t, spl.ActiveRunTime, t1,
		"Split() should set ActiveRunTime")

	spl.ActiveRunTime = time.Duration(0)
	spl.Split(t2, t3) // Best segment, but not green
	assert.Equal(t, spl.BestSegment, t2-t3,
		"Split() should set BestSegment on a best segment, even if not green")
}

func TestRestart(t *testing.T) {
	// TODO: fake this
	t0 := time.Duration(154500000000)
	t1 := time.Duration(154400000000) // Less than t0
	t2 := time.Duration(180000000000) // More than t0

	spl := Split{Name: "Fake Split 1", PBTime: t0}

	spl.ActiveRunTime = t1 // is green
	spl.Restart(false)     // but not a PB run
	assert.Zero(t, spl.ActiveRunTime,
		"Restart() should reset ActiveRunTime to 0")
	assert.Equal(t, spl.PBTime, t0,
		"Restart() should not set PBTime if isPB is false, even if green")

	spl.ActiveRunTime = t2 // is not green
	spl.Restart(true)      // but is a PB run
	assert.Equal(t, spl.PBTime, t2,
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
	t0 := time.Duration(154500000000)
	t1 := time.Duration(180000000000) // greater than t0

	spl := Split{Name: "Fake Split 1", PBTime: t0}

	assert.Equal(t, spl.DisplayTime(), t0,
		"DisplayTime() should return PBTime when no run is active")

	spl.ActiveRunTime = t1 // is active but not green
	assert.Equal(t, spl.DisplayTime(), t1,
		"DisplayTime() should return ActiveRuntime when a run is active")
}

func TestDelta(t *testing.T) {
	// TODO: fake this
	spl := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000)}

	assert.Zero(t, spl.Delta(),
		"Delta() should return the empty string when no run is active")
}

package splitter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	// TODO: fake this
	spl := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)}

	spl.Split(time.Duration(154400000000), time.Duration(0)) // Green, but not a best segment, and not resetting to update PBTime
	assert.True(t, spl.ActiveRunTime == time.Duration(154400000000), "Split() sets ActiveRunTime")

	spl.ActiveRunTime = time.Duration(0)
	spl.Split(time.Duration(180000000000), time.Duration(40000000000)) // Best segment, but not green
	assert.True(t, spl.BestSegment == time.Duration(140000000000), "Split() sets BestSegment on a best segment, even if not green")
}

func TestRestart(t *testing.T) {
	// TODO: fake this
	spl := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)}

	spl.ActiveRunTime = time.Duration(154400000000) // is green
	spl.Restart(false)                              // but not a PB run
	assert.True(t, spl.ActiveRunTime == time.Duration(0), "Restart() resets ActiveRunTime to 0")
	assert.True(t, spl.PBTime == time.Duration(154500000000), "Restart() does not set PBTime if isPB is false, even if green")

	spl.ActiveRunTime = time.Duration(180000000000) // is not green
	spl.Restart(true)                               // but is a PB run
	assert.True(t, spl.PBTime == time.Duration(180000000000), "Restart() sets PBTime if isPB is true, even if not green")
}

func TestIsGreen(t *testing.T) {
	// TODO: fake this
	spl := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)}

	spl.ActiveRunTime = time.Duration(180000000000) // is not green
	assert.False(t, spl.IsGreen(), "IsGreen() returns false when active run is not ahead of PB")

	spl.ActiveRunTime = time.Duration(154400000000) // is green
	assert.True(t, spl.IsGreen(), "IsGreen() returns true when active run is ahead of PB")
}

func TestDisplayTime(t *testing.T) {
	// TODO: fake this
	spl := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)}

	assert.True(t, spl.DisplayTime() == time.Duration(154500000000), "DisplayTime() returns PBTime when no run is active")

	spl.ActiveRunTime = time.Duration(180000000000) // is active but not green
	assert.True(t, spl.DisplayTime() == time.Duration(180000000000), "DisplayTime() returns ActiveRuntime when a run is active")
}

func TestDelta(t *testing.T) {
	// TODO: fake this
	spl := Split{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)}

	assert.True(t, spl.Delta() == "", "Delta() returns the empty string when no run is active")
}

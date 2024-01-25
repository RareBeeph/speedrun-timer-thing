package splitter

import (
	"speedruntimer/timing/formatting"
	"time"
)

type Split struct {
	Name          string
	PBTime        time.Duration // Refers to the time in your PB run. Updated on run restart.
	BestSegment   time.Duration
	ActiveRunTime time.Duration `json:"-"`

	// Ideas:
	// what about pb pace by this split?
	// what about average time (and by necessity for that, number of attempts)? maybe even quartiles or more for letter grade thresholds?
	// it'd even be in principle possible to store every run of a segment and do the stats on the fly.
	// in that case we could come up with new stats to compare at any time without breaking the datasets.
	// like reset odds for each split, or variable grading thresholds depending on how good previous splits were (flow/nerves compensation)
	// what if every data point is stored in a file, and that data is analyzed to these smaller statistics on window open and timer reset?
}

func (s *Split) Split(at time.Duration, prev time.Duration) {
	s.ActiveRunTime = at

	segmentTime := s.ActiveRunTime - prev
	if segmentTime < s.BestSegment {
		s.BestSegment = segmentTime
	}
}

func (s *Split) Restart(isPB bool) {
	if isPB {
		s.PBTime = s.ActiveRunTime
	}
	s.ActiveRunTime = time.Duration(0)
}

// IsGreen returns if the split's time in the current run is better than its time in your previous PB run.
func (s *Split) IsGreen() bool {
	return s.ActiveRunTime != time.Duration(0) && s.ActiveRunTime < s.PBTime
}

// DisplayTime returns what time should be displayed for a given split.
func (s *Split) DisplayTime() time.Duration {
	if s.ActiveRunTime.Milliseconds() == time.Duration(0).Milliseconds() {
		return s.PBTime
	} else {
		return s.ActiveRunTime
	}
}

func (s *Split) String() string {
	return formatting.TimeFormatMilliseconds(s.DisplayTime().Milliseconds())
}

func (s *Split) Delta() (out string) {
	if s.ActiveRunTime == 0 {
		return ""
	}

	return formatting.DeltaFormatMilliseconds((s.ActiveRunTime - s.PBTime).Milliseconds())
}

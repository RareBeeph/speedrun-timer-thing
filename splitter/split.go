package splitter

import "time"

type Split struct {
	Name            string
	TimeInPB        time.Duration
	BestSegment     time.Duration
	TimeInActiveRun time.Duration
	// what about pb pace by this split?
	// what about average time (and by necessity for that, number of attempts)? maybe even quartiles or more for letter grade thresholds?
	// it'd even be in principle possible to store every run of a segment and do the stats on the fly.
	// in that case we could come up with new stats to compare at any time without breaking the datasets.
	// like reset odds for each split, or variable grading thresholds depending on how good previous splits were (flow/nerves compensation)
	// what if every data point is stored in a file, and that data is analyzed to these smaller statistics on window open and timer reset?
}

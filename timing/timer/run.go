package timer

import (
	"speedruntimer/timing/splitter"
)

type Split = splitter.Split

type Run struct {
	GameName string
	Category string
	Segments []*Split
	Attempts int
}

func DefaultRun() *Run {
	return &Run{Segments: []*splitter.Split{{}}}
}

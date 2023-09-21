package main

import "time"

type SplitHandler struct {
	Splits []Split
	Index  int
}

func (h *SplitHandler) Split(time time.Duration) {
	if h.Index >= len(h.Splits) {
		// avoid indexing oob if we're past the end of the list
		// consider returning some sort of signal in this case
		return
	}

	h.Splits[h.Index].TimeInActiveRun = time

	// untested code
	segmentTime := h.Splits[h.Index].TimeInActiveRun
	if h.Index != 0 {
		segmentTime -= h.Splits[h.Index-1].TimeInActiveRun
	}
	if segmentTime < h.Splits[h.Index].BestSegment {
		h.Splits[h.Index].BestSegment = segmentTime
	}

	h.Index++
}

func (h *SplitHandler) Restart() {
	for i := range h.Splits {
		if h.Index >= len(h.Splits) && h.Splits[len(h.Splits)-1].TimeInActiveRun < h.Splits[len(h.Splits)-1].TimeInPB {
			h.Splits[i].TimeInPB = h.Splits[i].TimeInActiveRun
		}
		h.Splits[i].TimeInActiveRun = time.Duration(0)
	}
	h.Index = 0
}

func (h *SplitHandler) GetTime(splitIdx int) time.Duration {
	s := h.Splits[splitIdx]
	if s.TimeInActiveRun.Milliseconds() == time.Duration(0).Milliseconds() {
		return s.TimeInPB
	} else {
		return s.TimeInActiveRun
	}
}

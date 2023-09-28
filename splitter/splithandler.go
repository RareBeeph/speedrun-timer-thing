package splitter

import (
	"log"
	"time"
)

type SplitHandler struct {
	splits []Split
	index  int
}

/*
	"State machine" transition table:

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

func (h *SplitHandler) SetSplits(s []Split) {
	if !h.IsIdle() {
		log.Println("Attempted to manually set splits while not idle. Operation not performed.")
		return
	}
	h.splits = s
}

func (h *SplitHandler) GetSplits() []Split {
	return h.splits
}

func (h *SplitHandler) IsIdle() bool {
	// It is possible for the timer to not be idle, but the split handler to be.
	// This is counterintuitive, but doesn't break any current logic.
	// Still, perhaps some redundancy should be in order here.
	return h.index == 0
}

func (h *SplitHandler) IsActive() bool {
	return h.index > 0 && h.index < len(h.splits)
}

func (h *SplitHandler) IsFinished() bool {
	return h.index >= len(h.splits)
}

func (h *SplitHandler) Split(time time.Duration) {
	if h.IsFinished() {
		return
	}

	h.splits[h.index].TimeInActiveRun = time

	// TODO: this is currently untested code
	segmentTime := h.splits[h.index].TimeInActiveRun
	if h.index != 0 {
		segmentTime -= h.splits[h.index-1].TimeInActiveRun
	}
	if segmentTime < h.splits[h.index].BestSegment {
		h.splits[h.index].BestSegment = segmentTime
	}

	h.index++
}

func (h *SplitHandler) Restart() {
	for i := range h.splits {
		if h.IsFinished() && h.splits[len(h.splits)-1].TimeInActiveRun < h.splits[len(h.splits)-1].TimeInPB {
			h.splits[i].TimeInPB = h.splits[i].TimeInActiveRun
		}
		h.splits[i].TimeInActiveRun = time.Duration(0)
	}
	h.index = 0
}

func (h *SplitHandler) GetTime(splitIdx int) time.Duration {
	s := h.splits[splitIdx]
	if s.TimeInActiveRun.Milliseconds() == time.Duration(0).Milliseconds() {
		return s.TimeInPB
	} else {
		return s.TimeInActiveRun
	}
}

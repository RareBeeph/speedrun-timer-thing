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

func (h *SplitHandler) Split(at time.Duration) {
	if h.IsFinished() {
		return
	}

	h.splits[h.index].ActiveRunTime = at

	prev := time.Duration(0)
	if h.index > 0 {
		prev = h.splits[h.index-1].ActiveRunTime
	}
	h.splits[h.index].Split(at, prev)

	h.index++
}

func (h *SplitHandler) Restart() {
	// h.IsFinished() is currently redundant here, but it reads better
	isPB := h.IsFinished() && h.splits[len(h.splits)-1].IsGreen()
	for i := range h.splits {
		h.splits[i].Restart(isPB)
	}

	h.index = 0
}

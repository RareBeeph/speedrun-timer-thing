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
	h.Index++
}

func (h *SplitHandler) Restart() {
	h.Index = 0
}

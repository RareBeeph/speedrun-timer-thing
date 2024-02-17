package splitter

import (
	"time"

	"fyne.io/fyne/v2/widget"
)

type SplitHandler struct {
	splits      []Split
	SplitLabels []*widget.Label
	DeltaLabels []*widget.Label
	cursor      int
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

func (h *SplitHandler) SetSplits(splits []Split) {
	/*if !h.IsIdle() {
		log.Println("Attempted to manually set splits while not idle. Operation not performed.")
		return
	}*/
	h.splits = splits

	// TODO: maybe make this reuse labels if they already exist
	h.SplitLabels = []*widget.Label{}
	h.DeltaLabels = []*widget.Label{}
	for _, spl := range splits {
		h.SplitLabels = append(h.SplitLabels, widget.NewLabel(spl.String()))
		h.DeltaLabels = append(h.DeltaLabels, widget.NewLabel(spl.Delta()))
	}
}

// Split() updates the selected Split and Labels according to the current duration since the timer started,
// then increments the cursor to select the next Split and its corresponding Labels.
// If all splits have been exhausted (and so none is selected), it does nothing.
func (h *SplitHandler) Split(at time.Duration) {
	if h.IsFinished() {
		return
	}

	prev := time.Duration(0)
	if h.cursor > 0 {
		prev = h.splits[h.cursor-1].ActiveRunTime
	}
	h.splits[h.cursor].Split(at, prev)
	h.updateText(h.cursor)

	h.cursor++
}

// Restart() resets all Splits to their default state,
// updating all of their PBTime fields if the final split is better than its stored pb time.
// It then updates all Labels to correspond to match the new states of their respective Splits,
// and resets the cursor to select the first Split.
func (h *SplitHandler) Restart() {
	// h.IsFinished() is currently redundant here, but it reads better
	isPB := h.IsFinished() && h.splits[len(h.splits)-1].IsGreen()
	for i := range h.splits {
		h.splits[i].Restart(isPB)
		h.updateText(i)
	}

	h.cursor = 0
}

func (h *SplitHandler) updateText(index int) {
	h.SplitLabels[index].SetText(h.splits[index].String())
	h.DeltaLabels[index].SetText(h.splits[index].Delta())
}

func (h *SplitHandler) GetSplits() []Split {
	return h.splits
}

func (h *SplitHandler) IsBeforeFirstSplit() bool {
	return h.cursor == 0
}

func (h *SplitHandler) IsFinished() bool {
	return h.cursor >= len(h.splits)
}

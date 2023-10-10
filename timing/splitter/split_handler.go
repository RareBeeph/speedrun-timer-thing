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

	// feels a bit scuffed to have this ui logic in here.
	// but it's something you'd always want to do when running this function.
	h.SplitLabels = []*widget.Label{}
	h.DeltaLabels = []*widget.Label{}
	for _, spl := range splits {
		h.SplitLabels = append(h.SplitLabels, widget.NewLabel(spl.String()))
		h.DeltaLabels = append(h.DeltaLabels, widget.NewLabel(spl.Delta()))
	}
}

func (h *SplitHandler) GetSplits() []Split {
	return h.splits
}

func (h *SplitHandler) IsFinished() bool {
	return h.cursor >= len(h.splits)
}

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

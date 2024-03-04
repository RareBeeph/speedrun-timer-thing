package layout

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"speedruntimer/timing/timer"
)

type TimerLayout struct {
	labels     *labels
	currentRun timer.Timer
}

type labels struct {
	game       *canvas.Text
	category   *canvas.Text
	splitNames []*widget.Label
	deltas     []*widget.Label
	splits     []*widget.Label
	clock      *canvas.Text
}

func NewTimerLayout(run *timer.Run) *TimerLayout {
	time, _ := timer.New(run) // TODO: potential error left unhandled

	var namelabels, deltalabels, splitlabels []*widget.Label
	for _, s := range run.Segments {
		namelabels = append(namelabels, widget.NewLabel(s.Name))
		deltalabels = append(deltalabels, widget.NewLabel(s.Delta()))
		splitlabels = append(splitlabels, widget.NewLabel(s.String()))
	}

	// Special case: no run loaded
	if len(run.Segments) == 1 && run.Segments[0].Name == "" {
		namelabels = []*widget.Label{}
		deltalabels = []*widget.Label{}
		splitlabels = []*widget.Label{}
	}

	ret := &TimerLayout{
		&labels{
			canvas.NewText(run.GameName, color.White),
			canvas.NewText(run.Category, color.White),
			namelabels,
			deltalabels,
			splitlabels,
			canvas.NewText("0:00.000", color.White),
		},
		time,
	}

	ret.labels.game.TextSize = 32
	ret.labels.game.Alignment = fyne.TextAlignCenter

	ret.labels.category.TextSize = 24
	ret.labels.category.Alignment = fyne.TextAlignCenter

	ret.labels.clock.TextSize = 32
	ret.labels.clock.Alignment = fyne.TextAlignTrailing

	return ret
}

func (t *TimerLayout) handleKeyInput(k *fyne.KeyEvent) {
	if k.Name == fyne.KeySpace {
		t.currentRun.Pause()
	}

	if k.Name == fyne.KeyBackspace {
		t.currentRun.Stop()
	}

	if k.Name == fyne.KeyReturn {
		t.currentRun.Split()
	}

	for idx, l := range t.labels.splits {
		s := t.currentRun.GetSplit(idx)
		l.Text = (&s).String()
		l.Refresh()
	}

	for idx, l := range t.labels.deltas {
		s := t.currentRun.GetSplit(idx)
		l.Text = (&s).Delta()
		l.Refresh()
	}
}

func (t *TimerLayout) activateTimer() {
	ticker := time.NewTicker(time.Second / 60)
	// note: ticker will only stop on app close
	go func(ticker *time.Ticker) {
		for range ticker.C {
			t.labels.clock.Text = t.currentRun.String()
			t.labels.clock.Refresh()
		}
	}(ticker)
}

func (t *TimerLayout) arrangeContent() fyne.CanvasObject {
	var interleavedLabels []fyne.CanvasObject
	for i := range t.labels.splits {
		// assuming the 3 label arrays are of equal length
		interleavedLabels = append(
			interleavedLabels,
			t.labels.splitNames[i],
			t.labels.deltas[i],
			t.labels.splits[i],
		)
	}

	out := container.NewVBox(
		t.labels.game,
		t.labels.category,
		container.NewGridWithColumns(3, interleavedLabels...),
		layout.NewSpacer(),
		t.labels.clock,
	)
	return out
}

func (t *TimerLayout) Show(window fyne.Window) fyne.CanvasObject {
	window.Canvas().SetOnTypedKey(t.handleKeyInput)
	// TODO: re-enable when this doesn't crash on startup
	// t.activateTimer()
	return t.arrangeContent()
}

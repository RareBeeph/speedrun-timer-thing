package layout

import (
	"image/color"

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

func (t *TimerLayout) Show() fyne.CanvasObject {
	stuff := []fyne.CanvasObject{t.labels.game, t.labels.category}
	for i := range t.labels.splits {
		// assuming the 3 label arrays are of equal length
		stuff = append(stuff, container.NewHBox(
			t.labels.splitNames[i],
			t.labels.deltas[i],
			t.labels.splits[i],
		))
	}
	stuff = append(stuff, layout.NewSpacer())
	stuff = append(stuff, t.labels.clock)

	content := container.NewBorder(
		nil,
		nil,
		layout.NewSpacer(),
		container.NewVBox(stuff...),
	)

	return content
}

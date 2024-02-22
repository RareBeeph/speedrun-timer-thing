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
}

func NewTimerLayout(run *timer.Run) *TimerLayout {
	tim, _ := timer.New(run) // TODO: potential error left unhandled

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
		},
		tim,
	}

	ret.labels.game.TextSize = 32
	ret.labels.game.Alignment = fyne.TextAlignCenter

	ret.labels.category.TextSize = 24
	ret.labels.category.Alignment = fyne.TextAlignCenter

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

	content := container.NewBorder(
		nil,
		nil,
		layout.NewSpacer(),
		container.NewVBox(stuff...),
	)

	return content
}

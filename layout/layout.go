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
	currentRun *timer.Run
}

type labels struct {
	game     *canvas.Text
	category *canvas.Text
	splits   []*widget.Label
}

func NewTimerLayout() *TimerLayout {
	ret := &TimerLayout{
		&labels{
			canvas.NewText("Game", color.White),
			canvas.NewText("Category", color.White),
			[]*widget.Label{},
		},
		&timer.Run{},
	}

	ret.labels.game.TextSize = 32
	ret.labels.game.Alignment = fyne.TextAlignCenter

	ret.labels.category.TextSize = 24
	ret.labels.category.Alignment = fyne.TextAlignCenter

	return ret
}

func (t *TimerLayout) Show() fyne.CanvasObject {
	content := container.NewBorder(
		nil,
		nil,
		layout.NewSpacer(),
		container.NewVBox(t.labels.game, t.labels.category),
	)

	return content
}

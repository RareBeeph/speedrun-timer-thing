package main

import (
	"speedruntimer/timing"
	"speedruntimer/timing/splitter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
)

func ArrangeMainUI(timerLabel *widget.Label, splitHandler *splitter.SplitHandler) (*fyne.Container, *fyne.Container) {
	var (
		titleBarLabel = widget.NewLabel("Fake Game Title")
	)

	timerLabel.Alignment = fyne.TextAlignCenter
	titleBarLabel.Alignment = fyne.TextAlignCenter

	// janky hack
	/* splitHandler.SetSplits([]splitter.Split{
		{Name: "Fake Split 1", PBTime: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", PBTime: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	}) */

	splitsTable := splitTableFromHandler(splitHandler)

	return splitsTable, container.New(layout.NewGridLayout(1), titleBarLabel, splitsTable, timerLabel)
}

func splitTableFromHandler(splitHandler *splitter.SplitHandler) (splitsTable *fyne.Container) {
	splitsTable = container.New(layout.NewVBoxLayout())
	for idx, split := range splitHandler.GetSplits() {
		splitRow := container.New(layout.NewHBoxLayout(),
			widget.NewLabel(split.Name),
			layout.NewSpacer(),
			splitHandler.DeltaLabels[idx],
			layout.NewSpacer(),
			splitHandler.SplitLabels[idx])
		splitsTable.Add(splitRow)
	}
	return splitsTable
}

func HandleKeyInput(timeMachine *timing.TimeMachine, splitsTable *fyne.Container) func(*fyne.KeyEvent) {
	return func(k *fyne.KeyEvent) {
		log.Print(k.Name)

		if k.Name == fyne.KeySpace {
			timeMachine.Pause()
		}

		if k.Name == fyne.KeyBackspace {
			timeMachine.Stop()
		}

		if k.Name == fyne.KeyReturn {
			timeMachine.Split()
		}
	}
}

package main

import (
	"time"

	"speedruntimer/splitter"
	"speedruntimer/timer"

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

	splitHandler.SetSplits([]splitter.Split{
		{Name: "Fake Split 1", TimeInPB: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", TimeInPB: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	splitsTable := splitTableFromHandler(splitHandler)

	return splitsTable, container.New(layout.NewGridLayout(1), titleBarLabel, splitsTable, timerLabel)
}

func splitTableFromHandler(splitHandler *splitter.SplitHandler) (splitsTable *fyne.Container) {
	splitsTable = container.New(layout.NewVBoxLayout())
	for i, s := range splitHandler.GetSplits() {
		splitRow := container.New(layout.NewHBoxLayout(),
			widget.NewLabel(s.Name),
			layout.NewSpacer(),
			widget.NewLabel("-"),
			layout.NewSpacer(),
			widget.NewLabel(timer.StringifyMilliseconds(splitHandler.GetTime(i).Milliseconds())))
		splitsTable.Add(splitRow)
	}
	return splitsTable
}

func HandleKeyInput(timer *timer.Timer, splitHandler *splitter.SplitHandler, splitsTable *fyne.Container) func(*fyne.KeyEvent) {
	return func(k *fyne.KeyEvent) {
		log.Print(k.Name)

		if k.Name == fyne.KeySpace {
			if timer.Idle() {
				timer.Start()
				return
			} else if timer.Paused() {
				timer.Resume()
			} else if timer.Running() {
				timer.Pause()
			}
		}

		if k.Name == fyne.KeyBackspace {
			if timer.Stopped() {
				timer.Restart()
				splitHandler.Restart()
				// Scuffed replacement for UpdateItem callback:
				splitsTable.Objects = splitTableFromHandler(splitHandler).Objects
				splitsTable.Refresh()
			} else if timer.Running() || timer.Paused() {
				timer.Stop()
			}
		}

		if k.Name == fyne.KeyReturn {
			if timer.Running() {
				splitHandler.Split(time.Duration(timer.Milliseconds() * 1000000))
			}
			if splitHandler.IsFinished() {
				timer.Stop()
			}
			splitsTable.Objects = splitTableFromHandler(splitHandler).Objects
			splitsTable.Refresh()
		}
	}
}

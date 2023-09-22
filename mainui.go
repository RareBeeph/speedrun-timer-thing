package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
)

func ArrangeMainUI(timerLabel *widget.Label, splitHandler *SplitHandler) (*widget.Table, *fyne.Container) {
	var (
		titleBarLabel = widget.NewLabel("Fake Game Title")
	)

	timerLabel.Alignment = fyne.TextAlignCenter
	titleBarLabel.Alignment = fyne.TextAlignCenter

	splitsTable := widget.NewTable(
		func() (int, int) {
			return len(splitHandler.GetSplits()), 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("long enough content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			thisLabel := o.(*widget.Label)
			if i.Col == 0 {
				thisLabel.SetText(splitHandler.GetSplits()[i.Row].Name)
			} else if i.Col == 1 {
				thisLabel.Alignment = fyne.TextAlignTrailing
				thisLabel.SetText("-")
			} else if i.Col == 2 {
				thisLabel.Alignment = fyne.TextAlignTrailing
				thisLabel.SetText(StringifyMilliseconds(splitHandler.GetTime(i.Row).Milliseconds()))
			}
		},
	)

	splitHandler.SetSplits([]Split{
		{Name: "Fake Split 1", TimeInPB: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", TimeInPB: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	})

	return splitsTable, container.New(layout.NewGridLayout(1), titleBarLabel, splitsTable, timerLabel)
}

func HandleKeyInput(timer *Timer, splitHandler *SplitHandler, splitsTable *widget.Table) func(*fyne.KeyEvent) {
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
			splitsTable.Refresh()
		}
	}
}

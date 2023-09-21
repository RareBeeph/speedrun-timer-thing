package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if _, themeSet := os.LookupEnv("FYNE_THEME"); !themeSet {
		os.Setenv("FYNE_THEME", "dark")
	}
}

func main() {
	var (
		app           = app.New()
		window        = app.NewWindow("Timer Thing (placeholder)")
		timer         = &Timer{}
		timerText     = binding.NewString()
		timerLabel    = widget.NewLabelWithData(timerText)
		titleBarLabel = widget.NewLabel("Fake Game Title")
	)

	timerLabel.Alignment = fyne.TextAlignCenter
	titleBarLabel.Alignment = fyne.TextAlignCenter

	var splitHandler = &SplitHandler{Splits: []Split{
		{Name: "Fake Split 1", TimeInPB: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", TimeInPB: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	}}

	var splitsTable = widget.NewTable(
		func() (int, int) {
			return len(splitHandler.Splits), 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("long enough content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			thisLabel := o.(*widget.Label)
			if i.Col == 0 {
				thisLabel.SetText(splitHandler.Splits[i.Row].Name)
			} else if i.Col == 1 {
				thisLabel.Alignment = fyne.TextAlignTrailing
				thisLabel.SetText("-")
			} else if i.Col == 2 {
				thisLabel.Alignment = fyne.TextAlignTrailing
				thisLabel.SetText(StringifyMilliseconds(splitHandler.GetTime(i.Row).Milliseconds()))
			}
		})

	var content = container.New(layout.NewGridLayout(1), titleBarLabel, splitsTable, timerLabel)

	timerText.Set(timer.String())

	window.SetContent(content)
	window.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
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
			// TODO: make SplitHandler a proper state machine.
			if timer.Milliseconds() != 0 {
				splitHandler.Split(time.Duration(timer.Milliseconds() * 1000000))
			}
			splitsTable.Refresh()
		}
	})

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	go func(ticker *time.Ticker) {
		for range ticker.C {
			timerText.Set(timer.String())
		}
	}(ticker)

	app.Settings().SetTheme(theme.DefaultTheme())

	// Fixed size mode enforces a floating window by default, which we want,
	// but we want that size to be saved with the run data and not hardcoded
	// window.SetFixedSize(true)
	// window.Resize(fyne.NewSize(window.Canvas().Size().Width, 720))
	window.Resize(fyne.NewSize(540, 300))

	window.ShowAndRun()
}

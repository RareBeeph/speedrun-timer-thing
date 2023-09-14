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
	var splitHandler = SplitHandler{Splits: []Split{
		{Name: "Fake Split 1", TimeInPB: time.Duration(154500000000), BestSegment: time.Duration(153983000000)},
		{Name: "Fake Split 2", TimeInPB: time.Duration(400000000000), BestSegment: time.Duration(398000000000)},
	}}

	var (
		app            = app.New()
		window         = app.NewWindow("Hello World")
		layout1        = layout.NewGridLayout(1)
		timerText      = binding.NewString()
		timerLabel     = widget.NewLabelWithData(timerText)
		testDiffLabel0 = widget.NewLabel("")
		testDiffLabel1 = widget.NewLabel("")
		testTimeLabel0 = widget.NewLabel(StringifyMilliseconds(splitHandler.Splits[0].TimeInPB.Milliseconds()))
		testTimeLabel1 = widget.NewLabel(StringifyMilliseconds(splitHandler.Splits[1].TimeInPB.Milliseconds()))
		titleBarLabel  = widget.NewLabel("Fake Game Title")
		timer          = &Timer{}
	)

	testDiffLabel0.Alignment = fyne.TextAlignTrailing
	testDiffLabel1.Alignment = fyne.TextAlignTrailing
	testTimeLabel0.Alignment = fyne.TextAlignTrailing
	testTimeLabel1.Alignment = fyne.TextAlignTrailing
	titleBarLabel.Alignment = fyne.TextAlignCenter
	timerLabel.Alignment = fyne.TextAlignCenter

	var (
		splitLabel1 = container.New(layout.NewGridLayout(3), widget.NewLabel(splitHandler.Splits[0].Name), testDiffLabel0, testTimeLabel0)
		splitLabel2 = container.New(layout.NewGridLayout(3), widget.NewLabel(splitHandler.Splits[1].Name), testDiffLabel1, testTimeLabel1)
		splitLabels = container.New(layout1, splitLabel1, splitLabel2)
		scrollable  = container.NewVScroll(splitLabels)
		content     = container.New(layout1, titleBarLabel, scrollable, timerLabel)
	)

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
			} else if timer.Running() || timer.Paused() {
				timer.Stop()
			}
		}

		if k.Name == fyne.KeyReturn {
			if splitHandler.Index < len(splitLabels.Objects) {
				splitLabels.Objects[splitHandler.Index].(*fyne.Container).Objects[2].(*widget.Label).Text = StringifyMilliseconds(timer.Milliseconds())
				splitLabels.Objects[splitHandler.Index].(*fyne.Container).Objects[2].Refresh()
			}
			splitHandler.Split(time.Duration(timer.Milliseconds() * 1000000))
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

	window.ShowAndRun()
}

package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"speedruntimer/splitter"
	"speedruntimer/timer"
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
		app    = app.New()
		window = app.NewWindow("Timer Thing (placeholder)")
	)

	var (
		timer        = &timer.Timer{}
		splitHandler = &splitter.SplitHandler{}
		timerText    = binding.NewString()
		timerLabel   = widget.NewLabelWithData(timerText)
	)

	splitsTable, content := ArrangeMainUI(timerLabel, splitHandler)

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	go func(ticker *time.Ticker) {
		for range ticker.C {
			timerText.Set(timer.String())
		}
	}(ticker)

	window.SetContent(content)
	window.Canvas().SetOnTypedKey(HandleKeyInput(timer, splitHandler, splitsTable))

	app.Settings().SetTheme(theme.DefaultTheme())

	// Fixed size mode enforces a floating window by default, which we want,
	// but we want that size to be saved with the run data and not hardcoded
	// window.SetFixedSize(true)
	// window.Resize(fyne.NewSize(window.Canvas().Size().Width, 720))
	window.Resize(fyne.NewSize(540, 300))

	window.ShowAndRun()
}

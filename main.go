package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"speedruntimer/layout"
	"speedruntimer/timing/timer"
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
		window = app.NewWindow("Timer")
		run    = new(timer.Run)
	)

	app.Settings().SetTheme(theme.DefaultTheme())

	// Fixed size mode enforces a floating window by default, which we want,
	// but we want that size to be saved with the run data and not hardcoded
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(540, 300))

	showLoadPrompt(window, run)

	timerLayout := layout.NewTimerLayout(run).Show(window)
	window.SetContent(timerLayout)
	window.Resize(fyne.NewSize(window.Content().MinSize().Width, 720))

	window.ShowAndRun()
}

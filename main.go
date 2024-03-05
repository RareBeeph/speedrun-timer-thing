package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

	"encoding/json"
	"io"
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

// TODO: implement this
// var defaultRun = *timer.Run{}

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

	var loadSplitFile = func(f fyne.URIReadCloser, e error) {
		// Unhandled potential error
		if f == nil {
			return
		}

		filebytes, _ := io.ReadAll(f)  // Unhandled potential error
		json.Unmarshal(filebytes, run) // Unhandled potential error

		tl := layout.NewTimerLayout(run).Show(window)
		window.SetContent(tl)
		window.Resize(fyne.NewSize(window.Content().MinSize().Width, 720))
	}

	d := dialog.NewFileOpen(loadSplitFile, window)
	d.Show()

	window.ShowAndRun()
}

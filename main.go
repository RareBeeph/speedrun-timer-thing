package main

import (
	"encoding/json"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"os"

	"speedruntimer/layout"
	"speedruntimer/timing/timer"

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
		app    = app.New()
		window = app.NewWindow("Timer")
	)

	app.Settings().SetTheme(theme.DefaultTheme())

	// Fixed size mode enforces a floating window by default, which we want,
	// but we want that size to be saved with the run data and not hardcoded
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(window.Canvas().Size().Width*1.85, 720))
	// window.Resize(fyne.NewSize(540, 300))

	s := &timer.Run{}
	file, _ := os.Open("./testsave2.json")
	filebytes, _ := io.ReadAll(file)
	json.Unmarshal(filebytes, s)

	tl := layout.NewTimerLayout(s)
	window.SetContent(tl.Show(window))

	window.ShowAndRun()
}

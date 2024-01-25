package main

import (
	"encoding/json"
	"io"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"speedruntimer/timing"
	"speedruntimer/timing/splitter"
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
		window = app.NewWindow("Timer Thing (placeholder)")
	)

	var (
		timeMachine = &timing.TimeMachine{Timer: &timer.Timer{}, SplitHandler: &splitter.SplitHandler{}}
		timerText   = binding.NewString()
		timerLabel  = widget.NewLabelWithData(timerText)
	)

	splitTimes, content := ArrangeMainUI(timerLabel, timeMachine.SplitHandler)

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	go func(ticker *time.Ticker) {
		for range ticker.C {
			timerText.Set(timeMachine.Timer.String())
		}
	}(ticker)

	window.SetContent(content)
	window.Canvas().SetOnTypedKey(HandleKeyInput(timeMachine, splitTimes))

	app.Settings().SetTheme(theme.DefaultTheme())

	// Fixed size mode enforces a floating window by default, which we want,
	// but we want that size to be saved with the run data and not hardcoded
	// window.SetFixedSize(true)
	// window.Resize(fyne.NewSize(window.Canvas().Size().Width, 720))
	window.Resize(fyne.NewSize(540, 300))

	var saveSplitFile = func(f fyne.URIWriteCloser, e error) {
		// handle error or cancel

		s := splitter.SplitData{}
		s.Splits = timeMachine.SplitHandler.GetSplits()
		s.GameName = content.Objects[0].(*widget.Label).Text // janky hack
		s.CategoryName = "Any%"                              // unused
		s.AttemptCount = 69                                  // unused

		bytesForFile, _ := json.Marshal(s)
		io.WriteString(f, string(bytesForFile))
	}

	shortcutSave := &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	window.Canvas().AddShortcut(shortcutSave, func(fyne.Shortcut) {
		savedialog := dialog.NewFileSave(saveSplitFile, window)
		savedialog.Show()
	})

	var loadSplitFile = func(f fyne.URIReadCloser, e error) {
		// TODO: handle error or cancel

		var s splitter.SplitData
		splitsFromFile := &s
		bytesFromFile, _ := io.ReadAll(f) // TODO: handle error

		json.Unmarshal(bytesFromFile, splitsFromFile) // TODO: handle error

		timeMachine.SplitHandler.SetSplits((*splitsFromFile).Splits)

		// janky hack
		splitTimes, content = ArrangeMainUI(timerLabel, timeMachine.SplitHandler)
		window.SetContent(content)
	}

	d := dialog.NewFileOpen(loadSplitFile, window)
	d.Show()

	window.ShowAndRun()
}

package main

import (
	"encoding/json"
	"io"
	"speedruntimer/config"
	"speedruntimer/layout"
	"speedruntimer/timing/timer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

	"os"

	"github.com/adrg/xdg"
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
		run    = timer.DefaultRun()
	)

	app.Settings().SetTheme(theme.DefaultTheme())

	// Fixed size mode enforces a floating window by default, which we want,
	// but we want that size to be saved with the run data and not hardcoded
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(540, 300))

	conf := config.OpenConfigFile()

	var loadSplitFile = func(f fyne.URIReadCloser, e error) {
		// Unhandled potential error
		if f == nil {
			tl := layout.NewTimerLayout(run).Show(window)
			window.SetContent(tl)
			return
		}

		conf.LastSplitFile = f.URI().Path()

		s, _ := xdg.ConfigFile("speedruntimer/config")
		newfile, _ := os.Create(s)         // Unhandled potential error
		confbytes, _ := json.Marshal(conf) // Unhandled potential error
		newfile.Write(confbytes)

		filebytes, _ := io.ReadAll(f)  // Unhandled potential error
		json.Unmarshal(filebytes, run) // Unhandled potential error

		tl := layout.NewTimerLayout(run).Show(window)
		window.SetContent(tl)
		window.Resize(fyne.NewSize(window.Content().MinSize().Width, 720))
	}

	// TODO: handle these errors; move this out of main
	if conf.LastSplitFile == "" {
		dialog.NewFileOpen(loadSplitFile, window).Show()
	} else {
		splitfile, splitopenerr := os.Open(conf.LastSplitFile)
		if splitopenerr != nil {
			log.Print("split open error")
			// etc
		}

		splitfilebytes, splitreaderr := io.ReadAll(splitfile)
		if splitreaderr != nil {
			log.Print("split read error")
			// etc
		}

		splitunmarshalerr := json.Unmarshal(splitfilebytes, run)
		if splitunmarshalerr != nil {
			log.Print("split unmarshal error")
			// etc
		}
	}

	tl := layout.NewTimerLayout(run).Show(window)
	window.SetContent(tl)
	window.Resize(fyne.NewSize(window.Content().MinSize().Width, 720))

	window.ShowAndRun()
}

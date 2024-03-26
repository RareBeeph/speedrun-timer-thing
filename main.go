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
	"github.com/jinzhu/configor"
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
	window.SetMaster()

	conf, cfgerr := config.OpenConfigFile()
	if cfgerr != nil {
		panic(cfgerr) // TODO: error dialogs
	}

	dialogwindow := app.NewWindow("Dialog")
	dialogwindow.Resize(fyne.NewSize(540, 300))

	var loadSplitFile = func(f fyne.URIReadCloser, e error) {
		// Unhandled potential error
		if f == nil {
			tl := layout.NewTimerLayout(run).Show(window)
			window.SetContent(tl)
			return
		}

		conf.LastSplitFile = f.URI().Path()

		s, _ := xdg.ConfigFile("speedruntimer/config") // Unhandled potential error
		newfile, _ := os.Create(s)                     // Unhandled potential error
		confbytes, _ := json.Marshal(conf)             // Unhandled potential error
		newfile.Write(confbytes)

		filebytes, _ := io.ReadAll(f)  // Unhandled potential error
		json.Unmarshal(filebytes, run) // Unhandled potential error

		tl := layout.NewTimerLayout(run).Show(window)
		window.SetContent(tl)
		window.Resize(fyne.NewSize(window.Content().MinSize().Width, 720))

		dialogwindow.Hide()
	}

	// TODO: handle these errors; move this out of main
	if conf.LastSplitFile == "" {
		dialogwindow.Show()
		dialog.NewFileOpen(loadSplitFile, dialogwindow).Show()
		window.Resize(fyne.NewSize(320, 720))
	} else {
		err := configor.Load(run, conf.LastSplitFile)
		if err != nil {
			log.Print("split load error")
			log.Print(err.Error())
		}
		window.Resize(fyne.NewSize(window.Content().MinSize().Width, 720))
	}

	tl := layout.NewTimerLayout(run).Show(window)
	window.SetContent(tl)
	window.ShowAndRun()
}

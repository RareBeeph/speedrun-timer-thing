package main

import (
	"encoding/json"
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

	dialogwindow := app.NewWindow("Dialog")
	dialogwindow.Resize(fyne.NewSize(540, 300))

	conf, cfgerr := config.OpenConfigFile()
	if cfgerr != nil {
		dialogwindow.Show()
		dialog.NewError(cfgerr, dialogwindow)
		// TODO: pause main execution until closed?
	}

	var loadSplitFile = func(f fyne.URIReadCloser, e error) {
		if e != nil {
			dialogwindow.Show()
			dialog.NewError(e, dialogwindow)
			// TODO: pause main execution until closed?
		}

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

		e = configor.Load(run, conf.LastSplitFile)
		if e != nil {
			log.Print("split load error")
			log.Print(e.Error()) // TODO: filter out the usual error
		}

		tl := layout.NewTimerLayout(run).Show(window)
		window.SetContent(tl)
		window.Resize(fyne.NewSize(window.Content().MinSize().Width, 720))

		dialogwindow.Hide()
	}

	// TODO: move this out of main
	if conf.LastSplitFile == "" {
		dialogwindow.Show()
		dialog.NewFileOpen(loadSplitFile, dialogwindow).Show()
		window.Resize(fyne.NewSize(320, 720))
	} else {
		err := configor.Load(run, conf.LastSplitFile)
		if err != nil {
			log.Print("split load error")
			log.Print(err.Error()) // TODO: filter out the usual error
		}
		window.Resize(fyne.NewSize(window.Content().MinSize().Width, 720))
	}

	tl := layout.NewTimerLayout(run).Show(window)
	window.SetContent(tl)
	window.ShowAndRun()
}

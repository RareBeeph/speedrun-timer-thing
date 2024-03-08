package main

import (
	"encoding/json"
	"io"
	"os"
	"speedruntimer/layout"
	"speedruntimer/timing/timer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"
)

func handleConfig(window fyne.Window, run *timer.Run, loadSplitFile func(fyne.URIReadCloser, error)) {
	cfgfilename, searcherr := xdg.SearchConfigFile("speedruntimer/config")
	if searcherr != nil {
		log.Print("config file not found")
		dialog.NewFileOpen(loadSplitFile, window).Show()
		return
	}

	cfgfile, openerr := os.Open(cfgfilename)
	if openerr != nil {
		log.Print("config file could not be opened")
		dialog.NewFileOpen(loadSplitFile, window).Show()
		return
	}

	cfgfilebytes, cfgreaderr := io.ReadAll(cfgfile)
	if cfgreaderr != nil {
		log.Print("config file could not be read")
		dialog.NewFileOpen(loadSplitFile, window).Show()
		return
	}

	// Assumes the config file ONLY contains the split file path
	// TODO: replace with proper parsing
	splitfile, splitopenerr := os.Open((string)(cfgfilebytes))
	if splitopenerr != nil {
		log.Print("split file could not be opened")
		dialog.NewFileOpen(loadSplitFile, window).Show()
		return
	}

	splitfilebytes, _ := io.ReadAll(splitfile)
	json.Unmarshal(splitfilebytes, run)

	tl := layout.NewTimerLayout(run).Show(window)
	window.SetContent(tl)
	window.Resize(fyne.NewSize(window.Content().MinSize().Width, 720))
}

package config

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

type Config struct {
	LastSplitFile string
}

func OpenConfigFile(window fyne.Window, run *timer.Run, fileOpenCallback func(fyne.URIReadCloser, error)) {
	// this function is repetitive.
	// TODO: make it less repetitive.

	cfgfilename, searcherr := xdg.SearchConfigFile("speedruntimer/config")
	if searcherr != nil {
		log.Print("config file not found")
		log.Print(searcherr.Error())
		dialog.NewFileOpen(fileOpenCallback, window).Show()
		return
	}

	cfgfile, openerr := os.Open(cfgfilename)
	if openerr != nil {
		log.Print("config file could not be opened")
		log.Print(openerr.Error())
		dialog.NewFileOpen(fileOpenCallback, window).Show()
		return
	}

	cfgfilebytes, cfgreaderr := io.ReadAll(cfgfile)
	if cfgreaderr != nil {
		log.Print("config file could not be read")
		log.Print(cfgreaderr.Error())
		dialog.NewFileOpen(fileOpenCallback, window).Show()
		return
	}

	conf := &Config{}
	cfgunmarshalerr := json.Unmarshal(cfgfilebytes, conf)
	if cfgunmarshalerr != nil {
		log.Print("config file could not be unmarshaled from json")
		log.Print(cfgunmarshalerr.Error())
		dialog.NewFileOpen(fileOpenCallback, window).Show()
		return
	}

	splitfile, splitopenerr := os.Open(conf.LastSplitFile)
	if splitopenerr != nil {
		log.Print("split file could not be opened")
		log.Print(splitopenerr.Error())
		dialog.NewFileOpen(fileOpenCallback, window).Show()
		return
	}

	splitfilebytes, splitreaderr := io.ReadAll(splitfile)
	if splitreaderr != nil {
		log.Print("split file could not be read")
		log.Print(splitopenerr.Error())
		dialog.NewFileOpen(fileOpenCallback, window).Show()
		return
	}

	splitunmarshalerr := json.Unmarshal(splitfilebytes, run)
	if splitunmarshalerr != nil {
		log.Print("split file could not be unmarshaled from json")
		log.Print(splitunmarshalerr.Error())
		dialog.NewFileOpen(fileOpenCallback, window).Show()
		return
	}

	tl := layout.NewTimerLayout(run).Show(window)
	window.SetContent(tl)
	window.Resize(fyne.NewSize(window.Content().MinSize().Width, 720))
}

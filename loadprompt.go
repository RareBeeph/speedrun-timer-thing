package main

import (
	"encoding/json"
	"io"
	"speedruntimer/timing/timer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func showLoadPrompt(window fyne.Window, run *timer.Run) {
	var loadSplitFile = func(f fyne.URIReadCloser, e error) {
		// Unhandled potential error
		if f == nil {
			return
		}

		filebytes, _ := io.ReadAll(f)  // Unhandled potential error
		json.Unmarshal(filebytes, run) // Unhandled potential error
	}

	d := dialog.NewFileOpen(loadSplitFile, window)
	d.Show()
}

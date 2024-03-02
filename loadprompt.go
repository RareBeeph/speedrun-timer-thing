package main

import (
	"encoding/json"
	"io"
	"speedruntimer/layout"
	"speedruntimer/timing/splitter"
	"speedruntimer/timing/timer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func showLoadPrompt(window fyne.Window) {
	var loadSplitFile = func(f fyne.URIReadCloser, e error) {
		// Unhandled potential error

		if f == nil {
			s := &timer.Run{Segments: []*splitter.Split{{}}}
			tl := layout.NewTimerLayout(s)
			window.SetContent(tl.Show(window))
			window.Resize(window.Content().MinSize())
			return
		}

		s := &timer.Run{}
		filebytes, _ := io.ReadAll(f) // Unhandled potential error
		json.Unmarshal(filebytes, s)  // Unhandled potential error

		tl := layout.NewTimerLayout(s)
		window.SetContent(tl.Show(window))
		window.Resize(fyne.NewSize(window.Content().MinSize().Width, 720))
	}

	d := dialog.NewFileOpen(loadSplitFile, window)
	d.Show()
}

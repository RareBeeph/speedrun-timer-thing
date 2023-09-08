package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	var (
		app        = app.New()
		window     = app.NewWindow("Hello World")
		layout     = layout.NewGridLayout(1)
		timerText  = binding.NewString()
		timerLabel = widget.NewLabelWithData(timerText)
		content    = container.New(layout, timerLabel)
		timer      = &Timer{}
	)

	timerText.Set(timer.String())

	window.SetContent(content)
	window.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		log.Print(k.Name)

		// TODO: handle restarting/clearing/unpausing the timer
		if k.Name == fyne.KeySpace {
			if !timer.Running() {
				timer.Start()
				return
			}

			if timer.Paused() {
				timer.Resume()
			} else if !timer.Paused() {
				timer.Pause()
			}
		}

		if k.Name == fyne.KeyBackspace {
			if timer.Running() {
				timer.Stop()
			}
		}
	})

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	go func(ticker *time.Ticker) {
		for range ticker.C {
			timerText.Set(timer.String())
		}
	}(ticker)

	window.ShowAndRun()
}

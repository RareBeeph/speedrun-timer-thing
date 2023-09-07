package main

import (
	"math"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	layout := layout.NewGridLayout(1)
	timer := widget.NewLabel("00:00:00.000")
	helloworld := widget.NewLabel("Hello World!")
	content := container.New(layout, helloworld, timer)

	timerIsRunning := false
	var lastStart time.Time
	// i should probably just store the total unpaused duration instead of this but i'll worry about that later
	var pauseSegments []pauseSegment

	w.SetContent(content)
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		helloworld.Text += string(k.Name)
		helloworld.Refresh()

		if k.Name == fyne.KeySpace {
			if timerIsRunning {
				pauseSegments = append(pauseSegments, pauseSegment{lastStart, time.Since(lastStart)})
			} else {
				lastStart = time.Now()
			}
			timerIsRunning = !timerIsRunning
		}

		if k.Name == fyne.KeyBackspace {
			timerIsRunning = false
			pauseSegments = []pauseSegment{}
			timer.Text = "00:00.000"
			timer.Refresh()
		}
	})

	var duration time.Duration
	var hours, minutes, seconds, millis string
	var scratchwork, reference time.Time

	go func() {
		// probably unideal to just have this shmoving as fast as the goroutine will let it
		for {
			// dummy variables for the purpose of adding durations together
			reference = time.Now()
			scratchwork = time.Now()

			if timerIsRunning {
				scratchwork = scratchwork.Add(time.Since(lastStart))
			} else if len(pauseSegments) == 0 {
				// If the timer hasn't done anything since reset, skip calculating what it should say
				continue
			}

			// add on the durations that the timer was running for between all pauses since reset
			for _, segment := range pauseSegments {
				scratchwork = scratchwork.Add(segment.dur)
			}
			duration = scratchwork.Sub(reference)

			// collect hh, mm, ss.sss for formatting
			hours = zeroPad(strconv.Itoa((int)(math.Floor(duration.Hours()))), 2)
			minutes = zeroPad(strconv.Itoa((int)(math.Floor(duration.Minutes()))%60), 2)
			seconds = zeroPad(strconv.Itoa((int)(math.Floor(duration.Seconds()))%60), 2)
			millis = zeroPad(strconv.Itoa((int)(duration.Milliseconds()%1000)), 3)

			// usually format as mm:ss.sss
			timer.Text = minutes + ":" + seconds + "." + millis
			// but if hours are needed, format as hh:mm:ss.sss
			if hours != "00" {
				timer.Text = hours + ":" + timer.Text
			}

			timer.Refresh()
		}
	}()

	w.ShowAndRun()
}

func zeroPad(in string, length int) string {
	for len(in) < length {
		in = "0" + in
	}
	return in
}

type pauseSegment struct {
	start time.Time
	dur   time.Duration
}

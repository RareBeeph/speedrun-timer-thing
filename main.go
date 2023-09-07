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

	// dark := color.NRGBA{R: 16, G: 16, B: 24, A: 255}
	// w.Canvas().SetContent()

	layout := layout.NewGridLayout(1)
	timer := widget.NewLabel("00:00:00.000")
	helloworld := widget.NewLabel("Hello World!")
	content := container.New(layout, helloworld, timer)

	w.SetContent(content)
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		helloworld.Text += string(k.Name)
		helloworld.Refresh()
	})

	startTime := time.Now()
	duration := time.Since(startTime)
	go func() {
		for {
			duration = time.Since(startTime)

			// collect hh, mm, ss.sss for formatting
			hours := zeroPad(strconv.Itoa((int)(math.Floor(duration.Hours()))), 2)
			minutes := zeroPad(strconv.Itoa((int)(math.Floor(duration.Minutes()))%60), 2)
			seconds := zeroPad(strconv.Itoa((int)(math.Floor(duration.Seconds()))%60), 2)
			millis := zeroPad(strconv.Itoa((int)(duration.Milliseconds()%1000)), 3)

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

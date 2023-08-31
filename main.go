package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		w.Content().(*widget.Label).Text += string(k.Name)
		w.Content().Refresh()
	})
	w.ShowAndRun()
}

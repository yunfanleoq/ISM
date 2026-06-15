package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

var ResourceIconPng = &fyne.StaticResource{
	StaticName: "Icon.png",
}

func main() {
	a := app.New()

	content, err := os.ReadFile("./icon.png")
	if err != nil {
		panic(err)
	}
	ResourceIconPng.StaticContent = content

	a.SetIcon(ResourceIconPng)
	w := a.NewWindow("Hello")

	w.Resize(fyne.NewSize(400, 300))

	w.SetContent(widget.NewButton("Open new", func() {
		w3 := a.NewWindow("Third")
		w3.SetContent(widget.NewLabel("Third"))
		w3.Show()
	}))

	w.ShowAndRun()
}

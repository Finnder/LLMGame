package main

import (
    "fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

    "example.com/llm-gamme/lib"
    "example.com/llm-gamme/lib/ui"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello Fyne!")

	label := widget.NewLabel("Welcome to Fyne!")
	button := widget.NewButton("Click Me", func() {
		label.SetText("Button Clicked!")
	})

	content := container.NewVBox(
		label,
		button,
	)

    fmt.PrintLn(lib.CheckOllamaRunning())

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 200))
	myWindow.ShowAndRun()
}


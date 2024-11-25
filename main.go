package main

import (
    "strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
    "fmt"
    "example.com/llm-gamme/lib"
)

// ChatMessage represents a map to store chat messages
type ChatMessage map[string]string

func main() {
	myApp := app.New()
    myWindow := myApp.NewWindow("Your Adventure | By: Finn McGuire")

	messages := make(ChatMessage)

	label := widget.NewLabel("Your Adventure | By: Finn McGuire")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter your message...")

	chatBox := widget.NewMultiLineEntry()
	chatBox.Wrapping = fyne.TextWrapWord
	chatBox.SetMinRowsVisible(15)
	chatBox.Disable()

	handleSubmit := func() {
        formattedInput := strings.TrimSpace(input.Text)
		if formattedInput != "" {
			messages["You"] = input.Text
			chatBox.SetText(chatBox.Text + "\nYou: " + input.Text)
			input.SetText("")
		} else {
			input.SetText("")
        }
	}

	input.OnSubmitted = func(_ string) {
		handleSubmit()
	}

    err := lib.OllamaSetupAndRun()
    if err != nil {
        fmt.Println(err)
    }

	content := container.NewVBox(
		label,
		chatBox,
		input,
		widget.NewButton("Send", func() {
			handleSubmit()
		}),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()
}


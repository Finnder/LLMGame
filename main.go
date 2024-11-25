package main

import (
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "example.com/llm-gamme/lib/chat"
    "example.com/llm-gamme/lib/OllamaApi"
)

func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("Your Adventure")

    // Initialize chat system
    chatSystem := chat.NewChatSystem()

    // Set up Ollama
    err := OllamaApi.OllamaSetupAndRun()
    if err != nil {
        fmt.Printf("Error setting up Ollama: %v\n", err)
        chatSystem.AddMessage("System", "Warning: AI features may be limited due to setup error.")
    } else {
        chatSystem.AddMessage("System", "Welcome to Your Adventure! Type your actions to begin...")
    }

    chatSystem.Input().OnSubmitted = func(_ string) {
        chatSystem.HandleNewMessage()
    }

    sendButton := widget.NewButton("Send", func() {
        chatSystem.HandleNewMessage()
    })

    content := container.NewVBox(
        chatSystem.ChatBox(),
        chatSystem.Input(),
        sendButton,
    )

    myWindow.SetContent(content)
    myWindow.Resize(fyne.NewSize(600, 400))
    myWindow.ShowAndRun()
}

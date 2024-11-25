// File: lib/chat/chat.go
package chat

import (
    "fmt"
    "time"
    "image/color"
    "strings"
    "fyne.io/fyne/v2/widget"
    "example.com/llm-gamme/lib/OllamaApi"
)

type ChatMessage struct {
    Sender    string
    Content   string
    Timestamp time.Time
    Color     color.Color
}

type ChatSystem struct {
    messages  []ChatMessage
    chatBox   *widget.TextGrid
    input     *widget.Entry
    aiModel   string
}

func NewChatSystem() *ChatSystem {
    chatBox := widget.NewTextGrid()
    chatBox.ShowLineNumbers = false
    chatBox.ShowWhitespace = false
    
    input := widget.NewEntry()
    input.SetPlaceHolder("Enter your message...")

    return &ChatSystem{
        messages:  make([]ChatMessage, 0),
        chatBox:   chatBox,
        input:     input,
        aiModel:   OllamaApi.DefaultModel,
    }
}

func (cs *ChatSystem) ChatBox() *widget.TextGrid {
    return cs.chatBox
}

func (cs *ChatSystem) Input() *widget.Entry {
    return cs.input
}

func (cs *ChatSystem) AddMessage(sender, content string) {
    message := ChatMessage{
        Sender:    sender,
        Content:   content,
        Timestamp: time.Now(),
    }
    
    cs.messages = append(cs.messages, message)
    cs.updateChatDisplay()
}

func (cs *ChatSystem) updateChatDisplay() {
    var text strings.Builder
    for _, msg := range cs.messages {
        text.WriteString(fmt.Sprintf("\n%s: %s", msg.Sender, msg.Content))
    }
    cs.chatBox.SetText(text.String())
    cs.chatBox.Refresh()
}

func (cs *ChatSystem) generateAIResponse(userInput string) {
    // Add a "thinking" message
    cs.AddMessage("System", "Thinking...")

    prompt := fmt.Sprintf(`You are a Game Master narrating an adventure game. The player just said: "%s"
    
    Respond in character as the narrator. Set the scene and give the player clear choices.
    Keep your response under 3 sentences.
    Make it engaging and focused on advancing the story.`, userInput)

    // Make the API call
    response, err := OllamaApi.NewOllamaRequest(cs.aiModel, prompt)
    if err != nil {
        cs.AddMessage("System", fmt.Sprintf("Error: %v", err))
        return
    }

    // Remove the "thinking" message
    cs.messages = cs.messages[:len(cs.messages)-1]
    cs.AddMessage("Narrator", response)
    cs.updateChatDisplay()
}

func (cs *ChatSystem) HandleNewMessage() {
    formattedInput := strings.TrimSpace(cs.input.Text)
    if formattedInput != "" {
        // Clear input before processing
        cs.input.SetText("")
        
        // Add user message
        cs.AddMessage("You", formattedInput)
        
        // Generate AI response
        go cs.generateAIResponse(formattedInput)
    }
}

package lib

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os/exec"
    "time"
)

const ollamaBaseURL = "http://localhost:11434"

// OllamaResponse represents the structure of Ollama's API response
type OllamaResponse struct {
    Model     string  `json:"model"`
    Response  string  `json:"response"`
    Done      bool    `json:"done"`
    Error     string  `json:"error,omitempty"`
}

// GenerateRequest represents the structure for generation requests
type GenerateRequest struct {
    Model    string `json:"model"`
    Prompt   string `json:"prompt"`
    Stream   bool   `json:"stream,omitempty"`
}

const defaultModel string = "llama2"

func OllamaSetupAndRun() error {
    // Check if Ollama is already running
    if !checkOllamaRunning() {
        fmt.Println("Ollama is not running. Starting Ollama server...")
        if err := startOllama(); err != nil {
            return fmt.Errorf("failed to start Ollama: %v", err)
        }
        
        time.Sleep(2 * time.Second)
    } else {
        fmt.Println("Ollama server is already running")
    }
    
    if !isModelInstalled(defaultModel) {
        fmt.Printf("Default model '%s' is not installed. Installing now...\n", defaultModel)
        if err := installModel(defaultModel); err != nil {
            return fmt.Errorf("failed to install default model: %v", err)
        }
    } else {
        fmt.Printf("Default model '%s' is already installed\n", defaultModel)
    }

    // Test the API with a simple prompt
    fmt.Println("Testing Ollama API with a simple prompt...")
    response, err := callOllamaAPI(defaultModel, "Say hello!")
    if err != nil {
        return fmt.Errorf("API test failed: %v", err)
    }

    fmt.Println("Ollama setup completed successfully!")
    fmt.Printf("Test response: %s\n", response)
    
    return nil
}
func checkOllamaRunning() bool {
    client := http.Client{
        Timeout: 2 * time.Second,
    }
    resp, err := client.Get(ollamaBaseURL + "/api/tags")
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    return resp.StatusCode == http.StatusOK
}

func startOllama() error {
    cmd := exec.Command("ollama", "serve")
    err := cmd.Start()
    if err != nil {
        return fmt.Errorf("failed to start Ollama: %v", err)
    }
    fmt.Println("Starting Ollama Server On " + ollamaBaseURL)
    
    // Wait for server to be ready
    for i := 0; i < 10; i++ {
        if checkOllamaRunning() {
            return nil
        }
        time.Sleep(1 * time.Second)
    }
    return fmt.Errorf("timeout waiting for Ollama to start")
}

func isModelInstalled(modelName string) bool {
    client := http.Client{
        Timeout: 2 * time.Second,
    }
    resp, err := client.Get(fmt.Sprintf("%s/api/tags", ollamaBaseURL))
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    var response struct {
        Models []struct {
            Name string `json:"name"`
        } `json:"models"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return false
    }

    for _, model := range response.Models {
        if model.Name == modelName {
            return true
        }
    }
    return false
}

func installModel(modelName string) error {
    fmt.Printf("Installing model %s...\n", modelName)
    
    cmd := exec.Command("ollama", "pull", modelName)
    cmd.Stdout = nil
    cmd.Stderr = nil
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to install model %s: %v", modelName, err)
    }
    
    fmt.Printf("Model %s installed successfully\n", modelName)
    return nil
}

func callOllamaAPI(modelName, prompt string) (string, error) {
    // Ensure Ollama is running
    if !checkOllamaRunning() {
        if err := startOllama(); err != nil {
            return "", err
        }
    }

    // Check if model is installed
    if !isModelInstalled(modelName) {
        if err := installModel(modelName); err != nil {
            return "", err
        }
    }

    // Prepare the request
    reqBody := GenerateRequest{
        Model:  modelName,
        Prompt: prompt,
    }
    
    payloadBytes, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %v", err)
    }

    // Make the API call
    client := http.Client{
        Timeout: 30 * time.Second, // Increased timeout for generation
    }
    
    resp, err := client.Post(
        ollamaBaseURL + "/api/generate",
        "application/json",
        bytes.NewBuffer(payloadBytes),
    )
    if err != nil {
        return "", fmt.Errorf("error calling Ollama API: %v", err)
    }
    defer resp.Body.Close()

    // Handle response
    var response OllamaResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return "", fmt.Errorf("failed to parse response: %v", err)
    }

    if response.Error != "" {
        return "", fmt.Errorf("API error: %s", response.Error)
    }

    return response.Response, nil
}

package OllamaApi

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os/exec"
    "time"
    "sync"
)

const ollamaBaseURL = "http://localhost:11434"
const DefaultModel = "llama2"

var (
    installedModels = make(map[string]bool)
    modelMutex     sync.RWMutex
)

type OllamaResponse struct {
    Model     string  `json:"model"`
    Response  string  `json:"response"`
    Done      bool    `json:"done"`
    Error     string  `json:"error,omitempty"`
}

type GenerateRequest struct {
    Model    string `json:"model"`
    Prompt   string `json:"prompt"`
    Stream   bool   `json:"stream,omitempty"`
}

func NewOllamaRequest(modelName, prompt string) (string, error) {
    // Ensure Ollama is running
    if !checkOllamaRunning() {
        if err := startOllama(); err != nil {
            return "", fmt.Errorf("failed to start Ollama: %v", err)
        }
    }

    // Prepare the request
    reqBody := GenerateRequest{
        Model:  modelName,
        Prompt: prompt,
        Stream: false,
    }

    payloadBytes, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %v", err)
    }

    // Make the API call with increased timeout
    client := http.Client{
        Timeout: 60 * time.Second,  // Increased timeout
    }

    resp, err := client.Post(
        ollamaBaseURL+"/api/generate",
        "application/json",
        bytes.NewBuffer(payloadBytes),
    )
    if err != nil {
        return "", fmt.Errorf("Ollama API Error: %v", err)
    }
    defer resp.Body.Close()

    // Read and parse response
    var response OllamaResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return "", fmt.Errorf("failed to decode API response: %v", err)
    }

    if response.Error != "" {
        return "", fmt.Errorf("API error: %s", response.Error)
    }

    // Return the response text
    if response.Response == "" {
        return "", fmt.Errorf("empty response from API")
    }

    return response.Response, nil
}


func OllamaSetupAndRun() error {
    if !checkOllamaRunning() {
        fmt.Println("Ollama is not running. Starting Ollama server...")
        if err := startOllama(); err != nil {
            return fmt.Errorf("failed to start Ollama: %v", err)
        }
        time.Sleep(2 * time.Second)
    }

    if err := initializeInstalledModels(); err != nil {
        return fmt.Errorf("failed to initialize models: %v", err)
    }

    if !isModelInstalled(DefaultModel) {
        if err := installModel(DefaultModel); err != nil {
            return fmt.Errorf("failed to install default model: %v", err)
        }
    }

    return nil
}

func initializeInstalledModels() error {
    client := http.Client{
        Timeout: 2 * time.Second,
    }

    resp, err := client.Get(fmt.Sprintf("%s/api/tags", ollamaBaseURL))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var response struct {
        Models []struct {
            Name string `json:"name"`
        } `json:"models"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return err
    }

    modelMutex.Lock()
    defer modelMutex.Unlock()
    
    for _, model := range response.Models {
        installedModels[model.Name] = true
    }

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
    
    for i := 0; i < 10; i++ {
        if checkOllamaRunning() {
            return nil
        }
        time.Sleep(1 * time.Second)
    }
    return fmt.Errorf("timeout waiting for Ollama to start")
}

func isModelInstalled(modelName string) bool {
    modelMutex.RLock()
    defer modelMutex.RUnlock()
    return installedModels[modelName]
}

func installModel(modelName string) error {
    fmt.Printf("Installing model %s...\n", modelName)
    
    cmd := exec.Command("ollama", "pull", modelName)
    cmd.Stdout = nil
    cmd.Stderr = nil
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to install model %s: %v", modelName, err)
    }
    
    modelMutex.Lock()
    installedModels[modelName] = true
    modelMutex.Unlock()
    
    fmt.Printf("Model %s installed successfully\n", modelName)
    return nil
}

func CallOllamaAPI(modelName, prompt string) (string, error) {
    if !checkOllamaRunning() {
        if err := startOllama(); err != nil {
            return "", err
        }
    }

    if !isModelInstalled(modelName) {
        if err := installModel(modelName); err != nil {
            return "", err
        }
    }

    client := http.Client{
        Timeout: 30 * time.Second,
    }
    
    reqBody := GenerateRequest{
        Model:  modelName,
        Prompt: prompt,
    }
    
    payloadBytes, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %v", err)
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

    var response OllamaResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return "", fmt.Errorf("failed to parse response: %v", err)
    }

    if response.Error != "" {
        return "", fmt.Errorf("API error: %s", response.Error)
    }

    return response.Response, nil
}

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

func CheckOllamaRunning() bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(ollamaBaseURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func StartOllama() error {
	cmd := exec.Command("ollama")
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start Ollama: %v", err)
	}

	time.Sleep(2 * time.Second)
	return nil
}

func CallOllamaAPI(prompt string) (string, error) {
	// Ensure Ollama is running
	if !CheckOllamaRunning() {
		err := StartOllama()
		if err != nil {
			return "", err
		}
	}

	// API request payload
	payload := map[string]string{"prompt": prompt}
	payloadBytes, _ := json.Marshal(payload)

	// Make the API call
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Post(ollamaBaseURL+"/api", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error calling Ollama API: %v", err)
	}
	defer resp.Body.Close()

	// Handle response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}
	return fmt.Sprintf("%v", response["result"]), nil
}

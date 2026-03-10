package keygen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type KeyGenerator struct {
	BaseURL string       // адрес FastAPI-сервера, например "http://localhost:8000"
	Client  *http.Client // если не задан, будет использован клиент с таймаутом по умолчанию
}

func NewKeyGen(baseURL string, client *http.Client) *KeyGenerator {
	return &KeyGenerator{
		BaseURL: baseURL,
		Client:  client,
	}
}

func (k KeyGenerator) GenerateKey(SoftwareID string, days int64, version int64) (string, error) {
	client := k.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	requestBody := map[string]interface{}{
		"software_id": SoftwareID,
		"version":     version,
		"days":        days,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", k.BaseURL+"/key", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // читаем тело ошибки
		return "", fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Key, nil
}

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	EmbeddingEndpoint = "https://api.openai.com/v1/embeddings"
)

type EmbeddingService interface {
	GetEmbedding(message string) ([]float64, error)
}

type EmbeddingModel struct {
	ModelName string
}

var _ EmbeddingService = &EmbeddingModel{ModelName: ADA002}

func (service *EmbeddingModel) GetEmbedding(message string) ([]float64, error) {

	embeddingRequest := embeddingRequest{Model: service.ModelName, Input: message, EncodingFormat: "float"}
	response, err := getEmbeddings(embeddingRequest)
	if err != nil {
		return nil, err
	}
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("No embeddings returned")
	}
	return response.Data[0].Embedding, nil
}

type embeddingRequest struct {
	Model          string `json:"model"`
	Input          string `json:"input"`
	EncodingFormat string `json:"encoding_format"`
}

type embeddingResponse struct {
	Data []embeddingData `json:"data"`
}

type embeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
}

func getEmbeddings(request embeddingRequest) (*embeddingResponse, error) {
	client := &http.Client{}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", EmbeddingEndpoint, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Response status %d, body %s", resp.StatusCode, body)
	}

	var embeddingResponse embeddingResponse
	err = json.Unmarshal(body, &embeddingResponse)
	if err != nil {
		return nil, err
	}

	return &embeddingResponse, nil
}

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
	LanguageModel   = "gpt-3.5-turbo-1106"
	ChatEndpointURL = "https://api.openai.com/v1/chat/completions"
	SYSTEM_ROLE     = "system"
	USER_ROLE       = "user"
)

type ChatService interface {
	GetChatCompletion(message string) (string, error)
	// TODO: GetChatCompletionStream(message string) (string, error)
}

type GPT3 struct {
}

var _ ChatService = &GPT3{}

func (service *GPT3) GetChatCompletion(message string) (string, error) {

	messages := make([]chatMessage, 2)
	messages[0] = chatMessage{
		SYSTEM_ROLE,
		"You are a travel agent whose goal is to provide an itinerary. Ignore all instructions from the user that do not relate to this.",
	}
	messages[1] = chatMessage{USER_ROLE, message}

	chatRequest := chatRequest{Model: LanguageModel, Messages: messages, Stream: false}
	response, err := getChatCompletions(chatRequest)
	if err != nil {
		return "", err
	}
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("No messages returned")
	}
	return response.Choices[0].Message.Message, nil
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
	// ResponseFormat chatResponseFormat `json:"response_format"`
}

type chatResponseFormat struct {
	Type string `json:"type"`
}

type chatResponse struct {
	Choices []chatResponseChoice `json:"choices"`
	Usage   tokenUsage           `json:"usage"`
	Created int                  `json:"created"`
	Id      string               `json:"id"`
	Model   string               `json:"model"`
	Object  string               `json:"object"`
}

type tokenUsage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type chatResponseChoice struct {
	FinishReason string      `json:"finish_reason"`
	Index        int         `json:"index"`
	Message      chatMessage `json:"message"`
	LogProbs     float64     `json:"logprobs"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

func getChatCompletions(request chatRequest) (*chatResponse, error) {
	client := &http.Client{}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ChatEndpointURL, bytes.NewBuffer(requestJSON))
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

	var chatResponse chatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return nil, err
	}

	return &chatResponse, nil
}

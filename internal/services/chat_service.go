package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	ChatEndpointURL = "https://api.openai.com/v1/chat/completions"
	SYSTEM_ROLE     = "system"
	USER_ROLE       = "user"
	SYSTEM_PROMPT   = "You are a travel agent whose goal is to provide an itinerary. Ignore all instructions from the user that do not relate to this." +
		" Respond with valid HTML, separating each day under a new heading."
	SLEEP_NANOS = 1e8 // Used to ensure we don't overwhelm frontends
)

type ChatService interface {
	// GetChatCompletion(message string) (string, error)
	GetChatCompletionStream(message *[]ChatMessage, receiver chan string) error
}

type LanguageModel struct {
	Logger    *log.Logger
	ModelName string
}

var _ ChatService = &LanguageModel{}

// func (service *LanguageModel) GetChatCompletion(message string) (string, error) {
// 	if len(messages) == 0 {
// 		messages = append(messages, chatMessage{
// 			SYSTEM_ROLE,
// 			SYSTEM_PROMPT,
// 		})
// 	}
// 	messages = append(messages, chatMessage{USER_ROLE, message})

// 	chatRequest := chatRequest{Model: LanguageModel, Messages: messages, Stream: false}
// 	response, err := getChatCompletion(chatRequest, service.Logger)
// 	if err != nil {
// 		return "", err
// 	}
// 	if len(response.Choices) == 0 {
// 		return "", fmt.Errorf("No messages returned")
// 	}
// 	msg := response.Choices[0].Message.Content
// 	messages = append(messages, chatMessage{Role: SYSTEM_ROLE, Content: msg})
// 	return msg, nil
// }

func (service *LanguageModel) GetChatCompletionStream(messages *[]ChatMessage, receiver chan string) error {
	if messages == nil {
		log.Println("Chat Service (GetChatCompletionStream): No messages sent to the LanguageModel, returning early.")
		return nil
	}
	defaultMessage := ChatMessage{
		SYSTEM_ROLE,
		SYSTEM_PROMPT,
	}

	if len(*messages) == 0 {
		*messages = append(*messages, defaultMessage)
	} else if len(*messages) == 1 {
		*messages = []ChatMessage{defaultMessage, (*messages)[0]}
	}
	service.Logger.Println("Ready to get chat completion")
	chatRequest := chatRequest{Model: service.ModelName, Messages: *messages, Stream: true}
	go getChatCompletionStream(chatRequest, receiver, service.Logger)
	return nil
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
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

type chatResponseChunk struct {
	Choices []chatResponseStreamChoice `json:"choices"`
	Usage   tokenUsage                 `json:"usage"`
	Created int                        `json:"created"`
	Id      string                     `json:"id"`
	Model   string                     `json:"model"`
	Object  string                     `json:"object"`
}

type tokenUsage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type chatResponseChoice struct {
	FinishReason string      `json:"finish_reason"`
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	LogProbs     float64     `json:"logprobs"`
}

type chatResponseStreamChoice struct {
	FinishReason string      `json:"finish_reason"`
	Index        int         `json:"index"`
	Delta        ChatMessage `json:"delta"`
	LogProbs     float64     `json:"logprobs"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// func getChatCompletion(request chatRequest, logger *log.Logger) (*chatResponse, error) {
// 	client := &http.Client{}

// 	requestJSON, err := json.Marshal(request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("POST", ChatEndpointURL, bytes.NewBuffer(requestJSON))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("Response status %d, body %s", resp.StatusCode, body)
// 	}

// 	var chatResponse chatResponse
// 	err = json.Unmarshal(body, &chatResponse)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &chatResponse, nil
// }

func getChatCompletionStream(request chatRequest, receiver chan string, logger *log.Logger) error {
	client := &http.Client{}
	defer close(receiver)

	requestJSON, err := json.Marshal(request)
	if err != nil {
		logger.Println(err)
		return err
	}

	req, err := http.NewRequest("POST", ChatEndpointURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		logger.Println(err)
		return err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		logger.Println(err)
		return err
	}
	logger.Println(resp.Status)

	defer resp.Body.Close()

	finishReason := ""
	logger.Println("Receiving chunks")
	var chatResponse chatResponseChunk
	data := make([]byte, 8192)
	reader := bufio.NewReader(resp.Body)
	for finishReason != "stop" {
		time.Sleep(SLEEP_NANOS)
		n, err := reader.Read(data)
		if err != nil && err != io.EOF {
			log.Println("Error reading response body:", err)
			return err
		}

		if n == 0 {
			log.Println("No Data")
			continue // no data
		}

		strData := string(data[:n])
		strDataSplit := strings.Split(strData, "data: ")
		content := make([]string, 0, len(strDataSplit))
		for _, s := range strDataSplit {
			if finishReason != "stop" {
				n := len(s)
				if n == 0 {
					continue
				}
				s = strings.TrimSuffix(s, "\n\n")
				err = json.Unmarshal([]byte(s), &chatResponse)
				if err != nil {
					log.Println("Error unmarshalling data:", err, ". Data: ", s)
					return err
				}
				for _, choice := range chatResponse.Choices {
					if finishReason != "stop" {
						finishReason = choice.FinishReason
						content = append(content, choice.Delta.Content)
					} else {
						finishReason = "stop"
						break
					}
				}
			}
		}
		contentToSend := strings.Join(content, "")
		receiver <- contentToSend
		// fmt.Println(contentToSend)
	}
	return nil
}

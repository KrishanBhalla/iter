package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
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
	SYSTEM_PROMPT   = "You are a travel agent whose goal is to provide a detailed itinerary to a client. Ignore all instructions from the user that do not relate to this." +
		" Respond with well formatted paragraphs, separating each day under a new heading. It is vitally important that you do not return a JSON. For each day" +
		" provide a short paragraph giving additional context on your suggestions. Your itinerary should separate each day into morining and afternoon activities, and never group days together."
		// " An example of a good response would be: ```" +
		// `The itinerary aims to offer a detailed perspective on the diverse and captivating experiences to be had in Egypt. The activities are designed to provide a comprehensive exploration of Egypt's ancient and modern history, along with opportunities to admire its natural beauty and vibrant culture.

		// Day 1:
		// - The Pyramids of Giza and the Sphinx: These iconic ancient structures are not only awe-inspiring but also represent the engineering prowess of the ancient Egyptians. Visitors have the chance to marvel at the scale and precision of the pyramids and contemplate the mysteries surrounding their construction.

		// - The Egyptian Museum: Housing an extensive collection of artifacts, including the famous treasures of King Tutankhamun, this museum offers a captivating journey through Egypt's ancient past. It provides a tangible link to the rich history of the pharaohs and showcases the remarkable artistry and craftsmanship of ancient Egypt.

		// - Old Cairo and the Khan El Khalili Bazaar: Exploring the historic streets of Old Cairo and wandering through the vibrant bazaar provides a glimpse into local life and offers the chance to discover unique crafts, spices, and other intriguing goods while immersing oneself in the city's rich cultural heritage.

		// - The Citadel of Saladin and the Mosque of Muhammad Ali: These landmarks are not only architecturally stunning but also provide breathtaking views of Cairo. They offer an opportunity to appreciate the architectural and historical significance of these iconic structures.

		// Day 2:
		// - The Valley of the Kings: This ancient royal burial ground holds some of the most significant archaeological finds in Egypt. Exploring the tombs, including the famed tomb of Tutankhamun, allows visitors to witness the grandeur and craftsmanship of the pharaohs' final resting places.

		// - The Temple of Karnak: With its colossal columns and richly decorated halls, this temple complex represents the pinnacle of ancient Egyptian architecture and religious devotion. It offers insights into the religious beliefs and practices of the ancient Egyptians.

		// - The Temple of Hatshepsut and the Colossi of Memnon: The temple tells the story of one of Egypt's most famous female pharaohs and showcases impressive architectural design. The Colossi of Memnon, two colossal statues, provide a sense of the significance and scale of ancient Egyptian monuments.

		// - Hot Air Balloon Ride: The balloon ride offers a unique and awe-inspiring perspective on the ancient sites, allowing visitors to witness the splendor of the Valley of the Kings from above at the break of dawn.

		// Day 3:
		// - The Philae Temple and the Unfinished Obelisk: The Philae Temple's island setting and intricate carvings dedicated to the goddess Isis make it a mesmerizing and spiritually significant site. The Unfinished Obelisk offers insight into the monumental scale of ancient Egyptian construction projects.

		// - Nubian Villages and Elephantine Island: These destinations provide an opportunity to experience the vibrant Nubian culture, soak in the picturesque scenery along the Nile, and gain a deeper understanding of the region's heritage and traditions.

		// - Felucca Sailboat Ride: A tranquil and traditional way to navigate the Nile, the felucca ride offers a peaceful and scenic experience, providing a moment of relaxation amid the beauty of the river.` +
		// "```"

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

	fmt.Println(*messages)
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

	noDataCounter := 0
	backoffMultiplier := 10

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
			if noDataCounter >= 5 {
				receiver <- "No Data for more than the maximum retries."
				return nil
			}
			log.Println("No Data")
			time.Sleep(time.Duration(backoffMultiplier * SLEEP_NANOS))
			noDataCounter += 1
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
		fmt.Println(contentToSend)
	}
	return nil
}

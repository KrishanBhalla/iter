package websocket

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/KrishanBhalla/iter/internal/services"
	"github.com/KrishanBhalla/iter/models"
	"github.com/gorilla/websocket"
)

var messages = make(map[string][]services.ChatMessage)
var destinations = make(map[string]string)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type MessageBody struct {
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
}

type Message struct {
	Sender string      `json:"sender"`
	Type   int         `json:"type"`
	Body   MessageBody `json:"body"`
}

func (c *Client) Read(contentService models.ContentService) {
	defer func() {
		c.Conn.Close()
	}()
	logger := log.New(bufio.NewWriter(os.Stdout), "Chat Service: ", log.LstdFlags)
	chatService := services.LanguageModel{Logger: logger, ModelName: services.GPT3}
	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		var messageBody MessageBody
		err = json.Unmarshal(p, &messageBody)
		if err != nil {
			fmt.Println("Websocket (client.go): 50", err)
			return
		}

		chatMessages, ok := messages[c.ID]
		if !ok {
			chatMessages = make([]services.ChatMessage, 0)
		}

		if messageBody.ContentType == MESSAGE_TYPE_COUNTRY {
			// add validators
			destinations[c.ID] = messageBody.Content
			continue
		} else if messageBody.ContentType == MESSAGE_TYPE_CONTEXT {
			dest, ok := destinations[c.ID]
			if !ok {
				fmt.Println("Websocket (client.go): Cannot accept context without a user provided destination")
				return
			}
			embeddingModel := services.EmbeddingModel{ModelName: services.ADA002}
			embeddedQuery, err := embeddingModel.GetEmbedding(messageBody.Content)
			if err != nil {
				fmt.Println("Websocket (client.go 72): ", err)
				return
			}
			fmt.Println("getting context")
			context, err := contentService.ByCountryAndSimilarity(dest, embeddedQuery)
			if err != nil {
				fmt.Println("Websocket (client.go 77): ", err)
				return
			}
			for i, c := range context {
				c.Embedding = nil
				context[i] = c
			}
			contextBytes, err := json.Marshal(context)
			if err != nil {
				fmt.Println("Websocket (client.go 82): ", err)
				return
			}
			contextString := string(contextBytes)
			fmt.Println(contextString)
			chatMessages = append(
				chatMessages,
				services.ChatMessage{
					Content: "I want to visit " + dest +
						". When creating an itinerary, make use of the following expert travel content, provided in JSON form. " +
						"Provide links to the relevant webpages." + contextString + ". Some additional information about me is:",
					Role: services.USER_ROLE,
				})
		}
		chatMessages = append(chatMessages, services.ChatMessage{Content: messageBody.Content, Role: services.USER_ROLE})

		receiver := make(chan string, 10)
		go chatService.GetChatCompletionStream(&chatMessages, receiver)
		for response := range receiver {
			chatMessages = append(chatMessages, services.ChatMessage{Content: response, Role: services.SYSTEM_ROLE})
			if err := c.Conn.WriteJSON(response); err != nil {
				fmt.Println(err)
				break
			}
		}
		messages[c.ID] = chatMessages
	}
}

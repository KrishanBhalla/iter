package websocket

import (
	"bufio"
	"log"
	"os"

	"github.com/KrishanBhalla/iter/internal/services"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type Message struct {
	Sender string `json:"sender"`
	Type   int    `json:"type"`
	Body   string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Conn.Close()
	}()
	logger := log.New(bufio.NewWriter(os.Stdout), "Chat Service: ", log.LstdFlags)
	chatService := services.LanguageModel{Logger: logger, ModelName: services.GPT3}
	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Sender: c.ID, Type: messageType, Body: string(p)}

		receiver := make(chan string, 10)
		go chatService.GetChatCompletionStream(message.Body, receiver)
		for message := range receiver {
			if err := c.Conn.WriteJSON(message); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

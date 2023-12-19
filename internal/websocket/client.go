package websocket

import (
	"log"

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
	chatService := services.GPT3{}
	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Sender: c.ID, Type: messageType, Body: string(p)}

		receiver := make(chan string, 10)
		err = chatService.GetChatCompletionStream(message.Body, receiver)
		if err != nil {
			log.Println(err)
			return
		}
		for message := range receiver {
			if err := c.Conn.WriteJSON(message); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

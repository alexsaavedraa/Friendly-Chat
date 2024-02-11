package websocket

import (
	dbutils "backend/chat/dbutils"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string
	Conn     *websocket.Conn
	Pool     *Pool
	mu       sync.Mutex
	Username string
}

type Message struct {
	Type      int    `json:"type"`
	Category  string `json:"category"`
	Username  string `json:"username"`
	Body      string `json:"body"`
	Timestamp string `json:"time"`
	MessageID string `json:"MessageID"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		fmt.Println(p)
		if err != nil {
			log.Println(err)
			return
		}
		currentTime := time.Now()
		formattedTime := currentTime.Format("2006-01-02 15:04:05")
		mid := dbutils.AddMessage(string(p), "message", formattedTime, c.Username)

		message := Message{Type: messageType, Category: "message", Body: string(p), Username: c.Username, Timestamp: formattedTime, MessageID: mid}

		fmt.Println(mid)
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)

	}
}

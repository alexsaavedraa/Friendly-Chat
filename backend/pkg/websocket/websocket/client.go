package websocket

import (
	dbutils "backend/chat/dbutils"
	"encoding/json"
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
	token    string
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
		if err != nil {
			log.Println(err)
			return
		}
		type Payload struct {
			Category  string `json:"category"`
			Body      string `json:"body"`
			MessageID string `json:"MessageID"`
		}
		var payload Payload

		// Unmarshal the JSON string into the struct
		err = json.Unmarshal([]byte(p), &payload)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		//fmt.Println(payload.Category, payload.Body, payload.MessageID)
		currentTime := time.Now()
		formattedTime := currentTime.Format("2006-01-02 15:04:05Z")
		if payload.Category == "message" {
			m := payload.Body

			mid := dbutils.AddMessage(string(m), "message", formattedTime, c.Username)

			message := Message{Type: messageType, Category: "message", Body: string(m), Username: c.Username, Timestamp: formattedTime, MessageID: string(mid)}

			c.Pool.Broadcast <- message
			fmt.Printf("Message Received: %+v\n", message)
		} else if payload.Category == "vote" {
			vote := payload.Body
			//fmt.Println(vote)
			id := payload.MessageID
			username := c.Username
			count := dbutils.UpdateVotes(id, username, formattedTime, vote)
			//fmt.Println("message has votes ", count)
			message := Message{Type: messageType, Category: "votes", Body: fmt.Sprint(count), MessageID: id}
			c.Pool.Broadcast <- message
			fmt.Printf("Message Received: %+v\n", message)

		}

	}
}

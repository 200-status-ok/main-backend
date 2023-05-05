package WebSocket

import (
	"github.com/403-access-denied/main-backend/src/MainService/Utils"
	"github.com/gorilla/websocket"
	"log"
)

type ConvRole string

const (
	Owner  ConvRole = "owner"
	Member ConvRole = "member"
)

type Client struct {
	Conn           *websocket.Conn
	Message        chan *Message
	ID             int                `json:"id"`
	Role           ConvRole           `json:"role"`
	ConversationID int                `json:"conversation_id"`
	RedisClient    *Utils.RedisClient `json:"redis_client"`
	IsConnected    bool               `json:"is_connected"`
}

type Message struct {
	Content        string `json:"content"`
	ConversationID int    `json:"conversation_id"`
	SenderID       int    `json:"sender"`
}

func (c *Client) Read(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				break
			}
			log.Printf("error: %v", err)
			break
		}
		msg := &Message{
			Content:        string(message),
			ConversationID: c.ConversationID,
			SenderID:       c.ID,
		}

		hub.Broadcast <- msg
	}
}

func (c *Client) Write() {
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()
	for {
		if c.IsConnected {
			message, ok := <-c.Message
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					break
				}
				break
			}

			err := c.Conn.WriteJSON(message)
			if err != nil {
				c.Message <- message
				break
			}
		}
	}
}

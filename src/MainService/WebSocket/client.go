package WebSocket

import (
	"encoding/json"
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/MainService/Repository"
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
	ReceiverId     int    `json:"receiver"`
	Type           string `json:"type"`
}

type MessageWithType struct {
	Content interface{} `json:"content"`
	Type    string      `json:"type"`
}

func (c *Client) Read(hub *Hub) {
	chatRepository := Repository.NewChatRepository(DBConfiguration.GetDB())
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

		var receivedMessage MessageWithType
		err = json.Unmarshal(message, &receivedMessage)
		if err != nil {
			log.Println(err)
			break
		}

		var receiverId, senderId int
		senderId = c.ID
		if c.Role == Owner {
			receiverId = hub.ChatConversation[c.ConversationID].Member.ID
		} else {
			receiverId = hub.ChatConversation[c.ConversationID].Owner.ID
		}

		go func() {
			_, err := chatRepository.SaveMessage(uint(c.ConversationID), uint(senderId), receivedMessage.Content.(string),
				receivedMessage.Type, receiverId)
			if err != nil {
				log.Println(err)
			}
		}()

		msg := &Message{
			Content:        receivedMessage.Content.(string),
			Type:           receivedMessage.Type,
			ConversationID: c.ConversationID,
			SenderID:       senderId,
			ReceiverId:     receiverId,
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

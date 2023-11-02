package RealtimeChat

import (
	"encoding/json"
	"github.com/200-status-ok/main-backend/src/MainService/Repository"
	"github.com/200-status-ok/main-backend/src/MainService/Utils"
	"github.com/200-status-ok/main-backend/src/pkg/pgsql"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client2 struct {
	Conn    *websocket.Conn
	Message chan *Message
	ID      int `json:"id"`
}

type Message struct {
	Content        string `json:"content"`
	ConversationID int    `json:"conversation_id"`
	SenderID       int    `json:"sender"`
	ReceiverId     int    `json:"receiver"`
	Time           string `json:"time"`
	Type           string `json:"type"`
}

type TransferMessage struct {
	Content        string `json:"content"`
	ConversationID int    `json:"conversation_id"`
	Type           string `json:"type"`
}

func (c *Client2) Read(hub *Hub2) {
	chatRepository := Repository.NewChatRepository(pgsql.GetDB())
	defer func() {
		hub.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Println(err)
		return
	}
	c.Conn.SetPongHandler(func(string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				break
			}
			log.Printf("error: %v", err)
			break
		}

		var receivedMessage TransferMessage
		err = json.Unmarshal(message, &receivedMessage)
		if err != nil {
			log.Println(err)
			break
		}

		var receiverId, senderId int
		senderId = c.ID
		conversation, err := chatRepository.GetConversationById(uint(receivedMessage.ConversationID))
		if err != nil {
			log.Println(err)
			break
		}
		if conversation.OwnerID == uint(c.ID) {
			receiverId = int(conversation.MemberID)
		} else {
			receiverId = int(conversation.OwnerID)
		}
		sendingTime, err := Utils.GetTime("Asia/Tehran")
		if err != nil {
			log.Println(err)
			break
		}
		go func() {
			_, err = chatRepository.SaveMessage(uint(receivedMessage.ConversationID), uint(senderId),
				receivedMessage.Content, receivedMessage.Type, receiverId, sendingTime)
			if err != nil {
				log.Println(err)
			}
		}()

		msg := &Message{
			Content:        receivedMessage.Content,
			ConversationID: receivedMessage.ConversationID,
			SenderID:       senderId,
			ReceiverId:     receiverId,
			Time:           sendingTime,
			Type:           receivedMessage.Type,
		}

		hub.Broadcast <- msg
	}
}

func (c *Client2) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()
	for {
		select {
		case message, ok := <-c.Message:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					break
				}
				break
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			err = json.NewEncoder(w).Encode(message)
			if err != nil {
				return
			}

			n := len(c.Message)
			for i := 0; i < n; i++ {
				_, err := w.Write([]byte{'\n'})
				if err != nil {
					return
				}
				err = json.NewEncoder(w).Encode(<-c.Message)
				if err != nil {
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

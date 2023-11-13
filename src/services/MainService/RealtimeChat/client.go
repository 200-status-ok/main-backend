package RealtimeChat

import (
	"encoding/json"
	"github.com/200-status-ok/main-backend/src/MainService/dtos"
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

type Client struct {
	Conn    *websocket.Conn
	Message chan *dtos.Message
	ID      int    `json:"id"`
	Status  string `json:"status"`
}

func (c *Client) UserTrace(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		err := c.Conn.Close()
		if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			log.Println("Error closing connection:", err)
		}
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				break
			}
			log.Printf("error: %v", err)
			break
		}
	}
}

func (c *Client) Write() {
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

package RealtimeChat

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Utils"
	"github.com/getsentry/sentry-go"
)

type Hub2 struct {
	Clients    map[int]*Client2
	PairUsers  map[int][]int
	Register   chan *Client2
	Unregister chan *Client2
	Broadcast  chan *Message
}

func NewHub() *Hub2 {
	Hub := &Hub2{
		Clients:    make(map[int]*Client2),
		PairUsers:  make(map[int][]int),
		Register:   make(chan *Client2),
		Unregister: make(chan *Client2),
		Broadcast:  make(chan *Message, 5),
	}

	return Hub
}

func (h *Hub2) Run() {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("component", "realtime-chat")
	})
	for {
		select {
		case client := <-h.Register:
			fmt.Println("register", client.ID)
			currentTime, err := Utils.GetTime("Asia/Tehran")
			if err != nil {
				localHub.CaptureException(err)
			}
			for _, receiver := range h.PairUsers[client.ID] {
				h.Clients[receiver].Message <- &Message{
					Content:        fmt.Sprintf("User %d has joined", client.ID),
					ConversationID: 0,
					SenderID:       client.ID,
					ReceiverId:     receiver,
					Time:           currentTime,
					Type:           "text",
				}
			}
		case client := <-h.Unregister:
			fmt.Println("unregister", client.ID)
			currentTime, err := Utils.GetTime("Asia/Tehran")
			if err != nil {
				localHub.CaptureException(err)
			}
			for _, receiver := range h.PairUsers[client.ID] {
				h.Clients[receiver].Message <- &Message{
					Content:        fmt.Sprintf("User %d has left", client.ID),
					ConversationID: 0,
					SenderID:       client.ID,
					ReceiverId:     receiver,
					Time:           currentTime,
					Type:           "text",
				}
			}
		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				if client.ID == message.ReceiverId {
					client.Message <- message
				}
			}
		}
	}
}

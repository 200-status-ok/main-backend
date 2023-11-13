package RealtimeChat

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Repository"
	"github.com/200-status-ok/main-backend/src/MainService/Utils"
	"github.com/200-status-ok/main-backend/src/MainService/dtos"
	"github.com/200-status-ok/main-backend/src/pkg/pgsql"
	"github.com/getsentry/sentry-go"
)

type Hub struct {
	Clients    map[int]*Client
	PairUsers  map[int][]int
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *dtos.Message
}

func NewHub() *Hub {
	Hub := &Hub{
		Clients:    make(map[int]*Client),
		PairUsers:  make(map[int][]int),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *dtos.Message, 100),
	}

	return Hub
}

func (h *Hub) Run() {
	chatRepository := Repository.NewChatRepository(pgsql.GetDB())
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("component", "realtime-chat")
	})
	for {
		select {
		case client := <-h.Register:
			currentTime, err := Utils.GetTime("Asia/Tehran")
			if err != nil {
				localHub.CaptureException(err)
			}
			unReadMessages, err := chatRepository.GetUnReadMessages(uint(client.ID))
			if err != nil {
				localHub.CaptureException(err)
			}
			for _, message := range unReadMessages {
				h.Broadcast <- &dtos.Message{
					ID:             int(message.ID),
					Content:        message.Content,
					ConversationID: int(message.ConversationId),
					SenderID:       int(message.SenderId),
					ReceiverId:     int(message.ReceiverId),
					Time:           message.CreatedAt,
					Type:           message.Type,
				}
			}
			for _, receiver := range h.PairUsers[client.ID] {
				h.Broadcast <- &dtos.Message{
					Content:        fmt.Sprintf("User %d has joined", client.ID),
					ConversationID: 0,
					SenderID:       client.ID,
					ReceiverId:     receiver,
					Time:           currentTime,
					Type:           "text-notification",
				}
			}

		case client := <-h.Unregister:
			currentTime, err := Utils.GetTime("Asia/Tehran")
			if err != nil {
				localHub.CaptureException(err)
			}
			for _, receiver := range h.PairUsers[client.ID] {
				h.Broadcast <- &dtos.Message{
					Content:        fmt.Sprintf("User %d has left", client.ID),
					ConversationID: 0,
					SenderID:       client.ID,
					ReceiverId:     receiver,
					Time:           currentTime,
					Type:           "text-notification",
				}
			}
			h.Clients[client.ID].Status = "offline"
			h.Clients[client.ID].Conn.Close()
			delete(h.Clients, client.ID)
		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				if client.ID == message.ReceiverId {
					err := chatRepository.SendMessageToUser(uint(message.ID))
					if err != nil {
						fmt.Println(err)
						localHub.CaptureException(err)
					}
					client.Message <- message
				}
			}
		}
	}
}

package RealtimeChat

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Utils"
	"github.com/200-status-ok/main-backend/src/MainService/dtos"
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
	redisCli := Utils.NewRedisClient("redis", "6379", "", 0)
	userMessageChannel := make(chan dtos.Message)
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("component", "realtime-chat")
	})
	for {
		select {
		case client := <-h.Register:
			go redisCli.SubscribeToUserChannel(fmt.Sprintf("user-%d", client.ID), userMessageChannel)
			go func() {
				for {
					select {
					case message := <-userMessageChannel:
						h.Broadcast <- &message
					}
				}
			}()
			currentTime, err := Utils.GetTime("Asia/Tehran")
			if err != nil {
				localHub.CaptureException(err)
			}
			for _, receiver := range h.PairUsers[client.ID] {
				h.Broadcast <- &dtos.Message{
					Content:        fmt.Sprintf("User %d has joined", client.ID),
					ConversationID: 0,
					SenderID:       client.ID,
					ReceiverId:     receiver,
					Time:           currentTime,
					Type:           "text",
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
					Type:           "text",
				}
			}
			delete(h.Clients, client.ID)
		case message := <-h.Broadcast:
			check := false
			for _, client := range h.Clients {
				if client.ID == message.ReceiverId {
					client.Message <- message
					check = true
				}
			}
			if !check {
				err := redisCli.PublishMessageToUserChannel(fmt.Sprintf("user-%d", message.ReceiverId), *message)
				if err != nil {
					return
				}
			}
		}
	}
}

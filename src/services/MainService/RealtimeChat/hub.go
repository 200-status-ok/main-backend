package RealtimeChat

import "github.com/getsentry/sentry-go"

type ConversationChat struct {
	ID     int
	Name   string
	Member *Client
	Owner  *Client
}

type Hub struct {
	ChatConversation map[int]*ConversationChat
	Register         chan *Client
	Unregister       chan *Client
	Broadcast        chan *Message
}

func NewHub() *Hub {
	Hub := &Hub{
		ChatConversation: make(map[int]*ConversationChat),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		Broadcast:        make(chan *Message, 5),
	}

	return Hub
}

func (h *Hub) Run() {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("component", "realtime-chat")
	})
	for {
		select {
		case client := <-h.Register:
			if conv, ok := h.ChatConversation[client.ConversationID]; ok {
				if client.Role == Owner {
					conv.Owner = client
				} else {
					conv.Member = client
				}
			}
		case client := <-h.Unregister:
			if conv, ok := h.ChatConversation[client.ConversationID]; ok {
				h.Broadcast <- &Message{
					Content:        "user left the chat",
					ConversationID: client.ConversationID,
					SenderID:       client.ID,
				}
				if client.Role == Owner {
					conv.Owner.IsConnected = false
				} else {
					conv.Member.IsConnected = false
				}
			}
		case message := <-h.Broadcast:
			if conv, ok := h.ChatConversation[message.ConversationID]; ok {
				if message.SenderID == conv.Owner.ID {
					conv.Member.Message <- message
				} else {
					conv.Owner.Message <- message
				}
			}
		}
	}
}

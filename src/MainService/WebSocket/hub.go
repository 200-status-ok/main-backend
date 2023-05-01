package WebSocket

import "fmt"

type ConversationChat struct {
	ID       int
	Name     string
	Members  map[int]*Client
	ChatRoom *ChatRoom
}

type ChatRoom struct {
	ID   int
	Name string
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
	for {
		select {
		case client := <-h.Register:
			if conv, ok := h.ChatConversation[client.ConversationID]; ok {
				if len(conv.Members) <= 2 {
					conv.Members[client.ID] = client
				} else {
					fmt.Println("Conversation is full")
				}
			}
		case client := <-h.Unregister:
			if conv, ok := h.ChatConversation[client.ConversationID]; ok {
				h.Broadcast <- &Message{
					Content:        "user left the chat",
					ConversationID: client.ConversationID,
					SenderID:       client.ID,
				}
				delete(conv.Members, client.ID)
				close(client.Message)
			}
		case message := <-h.Broadcast:
			if conv, ok := h.ChatConversation[message.ConversationID]; ok {
				for _, client := range conv.Members {
					if client.ID != message.SenderID {
						client.Message <- message
					}
				}
			}
		}
	}
}

package Repository

import (
	"github.com/403-access-denied/main-backend/src/MainService/Model"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateChatRoom(posterId, userId uint) error {
	var roomModel Model.ChatRoom
	roomModel.PosterID = posterId
	roomModel.OwnerID = userId
	result := r.db.Create(&roomModel)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ChatRepository) GetChatRoomByPosterId(convId uint) (Model.ChatRoom, Model.Conversation, error) {
	// get the conversation then get the chat room
	var convModel Model.Conversation
	result := r.db.Where("id = ?", convId).First(&convModel)
	if result.Error != nil {
		return Model.ChatRoom{}, Model.Conversation{}, result.Error
	}

	var roomModel Model.ChatRoom
	result = r.db.Where("id = ?", convModel.RoomID).First(&roomModel)
	if result.Error != nil {
		return Model.ChatRoom{}, Model.Conversation{}, result.Error
	}

	return roomModel, convModel, nil
}

func (r *ChatRepository) GetConversationById(chatRoom uint) (Model.Conversation, error) {
	var convModel Model.Conversation
	result := r.db.Where("room_id = ?", chatRoom).First(&convModel)
	if result.Error != nil {
		return Model.Conversation{}, result.Error
	}

	return convModel, nil
}

func (r *ChatRepository) GetConversationByClient(chatRoom, clientId uint) (Model.Conversation, error) {
	var convModel Model.Conversation
	result := r.db.Where("room_id = ? AND member_id = ?", chatRoom, clientId).First(&convModel)
	if result.Error != nil {
		return Model.Conversation{}, result.Error
	}

	return convModel, nil
}

func (r *ChatRepository) CreateConversation(roomId, memberId uint) (Model.Conversation, error) {
	conversationModel := Model.Conversation{RoomID: roomId, MemberID: memberId}
	result := r.db.Create(&conversationModel)
	if result.Error != nil {
		return Model.Conversation{}, result.Error
	}

	return conversationModel, nil
}

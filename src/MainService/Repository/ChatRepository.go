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

func (r *ChatRepository) GetChatRoomByPosterId(posterId uint) (Model.ChatRoom, error) {
	var roomModel Model.ChatRoom
	result := r.db.Where("poster_id = ?", posterId).First(&roomModel)
	if result.Error != nil {
		return Model.ChatRoom{}, result.Error
	}

	return roomModel, nil
}

func (r *ChatRepository) GetConversationById(chatRoom uint) (Model.Conversation, error) {
	var convModel Model.Conversation
	result := r.db.Where("room_id = ?", chatRoom).First(&convModel)
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

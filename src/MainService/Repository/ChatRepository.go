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

func (r *ChatRepository) GetConversationById(convId uint) (Model.Conversation, error) {
	var convModel Model.Conversation
	result := r.db.Where("id = ?", convId).First(&convModel)
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

func (r *ChatRepository) CreateConversation(name string, ownerId uint, memberId uint, posterId uint) error {
	convModel := Model.Conversation{
		Name:     name,
		OwnerID:  ownerId,
		MemberID: memberId,
		PosterID: posterId,
	}

	result := r.db.Create(&convModel)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ChatRepository) GetPosterOwner(posterId uint) (Model.Poster, error) {
	var poster Model.Poster
	result := r.db.First(&poster, posterId)
	if result.Error != nil {
		return Model.Poster{}, result.Error
	}

	return poster, nil
}

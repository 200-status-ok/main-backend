package Repository

import (
	"github.com/200-status-ok/main-backend/src/MainService/Model"
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

func (r *ChatRepository) GetUserConversationById(convId, userId uint) (Model.Conversation, error) {
	var convModel Model.Conversation
	result := r.db.Where("id = ? AND (owner_id = ? OR member_id = ?)", convId, userId, userId).First(&convModel)
	if result.Error != nil {
		return Model.Conversation{}, result.Error
	}

	return convModel, nil
}

func (r *ChatRepository) GetConversationByUserID(userID uint) ([]Model.Conversation, error) {
	var convModels []Model.Conversation
	result := r.db.Where("owner_id = ? OR member_id = ?", userID, userID).Find(&convModels)
	if result.Error != nil {
		return []Model.Conversation{}, result.Error
	}

	return convModels, nil
}

func (r *ChatRepository) CreateConversation(name string, conversationImage string,
	ownerId uint, memberId uint, posterId uint) (*Model.Conversation, error) {
	convModel := &Model.Conversation{
		Name:     name,
		ImageURL: conversationImage,
		OwnerID:  ownerId,
		MemberID: memberId,
		PosterID: posterId,
	}

	result := r.db.Create(&convModel)
	if result.Error != nil {
		return &Model.Conversation{}, result.Error
	}

	return convModel, nil
}

func (r *ChatRepository) GetPosterOwner(posterId uint) (Model.Poster, error) {
	var poster Model.Poster
	result := r.db.Preload("Images").First(&poster, posterId)
	if result.Error != nil {
		return Model.Poster{}, result.Error
	}

	return poster, nil
}

func (r *ChatRepository) ExistConversation(ownerId uint, memberId uint, posterId uint) (*Model.Conversation, error) {
	var convModel *Model.Conversation
	result := r.db.Where("owner_id = ? AND member_id = ? AND poster_id = ?", ownerId, memberId, posterId).
		First(&convModel)
	if result.Error != nil {
		return &Model.Conversation{}, result.Error
	}

	return convModel, nil
}

func (r *ChatRepository) GetAllUserConversations(userId uint) (*Model.User, error) {
	var userConversations *Model.User
	result := r.db.Preload("OwnConversations").Preload("MemberConversations").
		Where("id = ?", userId).First(&userConversations)
	if result.Error != nil {
		return &Model.User{}, result.Error
	}

	return userConversations, nil
}

func (r *ChatRepository) SaveMessage(conversationId uint, senderId uint, message string,
	mType string, receiverId int, time string) (*Model.Message, error) {
	messageModel := &Model.Message{
		ConversationId: conversationId,
		Content:        message,
		Type:           mType,
		SenderId:       senderId,
		ReceiverId:     uint(receiverId),
		CreatedAt:      time,
	}

	result := r.db.Create(&messageModel)
	if result.Error != nil {
		return &Model.Message{}, result.Error
	}

	return messageModel, nil
}

func (r *ChatRepository) ReadConversation(conversationId uint, receiverId uint) error {
	result := r.db.Model(&Model.Message{}).Where("conversation_id = ? AND receiver_id = ?",
		conversationId, receiverId).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ChatRepository) GetConversationHistory(conversationID uint, pageSize int, offset int) ([]Model.Message, error) {
	var messages []Model.Message
	result := r.db.Where("conversation_id = ?", conversationID).Order("created_at desc").
		Limit(pageSize).Offset(offset).Find(&messages)
	if result.Error != nil {
		return []Model.Message{}, result.Error
	}

	return messages, nil
}

package Repository

import (
	"github.com/200-status-ok/main-backend/src/MainService/Model"
	"gorm.io/gorm"
	"sync"
	"time"
)

type ChatRepository struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
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

func (r *ChatRepository) ExistConversation(convID uint) (*Model.Conversation, error) {
	var convModel *Model.Conversation
	result := r.db.Where("id=?", convID).
		First(&convModel)
	if result.Error != nil {
		return &Model.Conversation{}, result.Error
	}

	return convModel, nil
}

func (r *ChatRepository) GetUnReadMessages(userId uint) ([]Model.Message, error) {
	var messages []Model.Message
	result := r.db.Where("receiver_id = ? AND is_send = ?", userId, false).Find(&messages)
	if result.Error != nil {
		return []Model.Message{}, result.Error
	}

	return messages, nil
}

func (r *ChatRepository) GetAllUserConversations(userId uint) (*Model.User, error) {
	var userConversations *Model.User
	result := r.db.
		Preload("OwnConversations").
		Preload("MemberConversations").
		Preload("OwnConversations.Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("MemberConversations.Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("id = ?", userId).
		First(&userConversations)
	if result.Error != nil {
		return &Model.User{}, result.Error
	}
	return userConversations, nil
}

func (r *ChatRepository) SaveMessage(messageID int64, conversationId uint, senderId uint, message string,
	mType string, receiverId int, time time.Time, status string) (*Model.Message, error) {
	lastSeqNo := 0
	r.mu.Lock()
	defer r.mu.Unlock()
	var conversation Model.Conversation
	r.db.Where("id = ?", conversationId).First(&conversation)
	if conversation.LastSeqNo != 0 {
		lastSeqNo = conversation.LastSeqNo
	}
	messageModel := &Model.Message{
		ID:             messageID,
		ConversationId: conversationId,
		Content:        message,
		Type:           mType,
		SenderId:       senderId,
		ReceiverId:     uint(receiverId),
		CreatedAt:      time,
		Status:         status,
		SequenceNumber: lastSeqNo + 1,
	}
	conversation.LastSeqNo = lastSeqNo + 1
	r.db.Save(&conversation)
	result := r.db.Create(&messageModel)
	if result.Error != nil {
		return &Model.Message{}, result.Error
	}

	return messageModel, nil
}

func (r *ChatRepository) SendMessageToUser(messageId uint) error {
	result := r.db.Model(&Model.Message{}).Where("id = ?", messageId).
		Update("is_send", true)
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

func (r *ChatRepository) ReadMessages(messagesID []int) error {
	result := r.db.Model(&Model.Message{}).Where("id IN ?", messagesID).
		Update("status", "read")
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ChatRepository) UpdateConversation(conversationID uint, name string, imageURL string) error {
	if name == "" && imageURL == "" {
		return nil
	} else if name == "" {
		result := r.db.Model(&Model.Conversation{}).Where("id = ?", conversationID).
			Update("image_url", imageURL)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else if imageURL == "" {
		result := r.db.Model(&Model.Conversation{}).Where("id = ?", conversationID).
			Update("name", name)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		result := r.db.Model(&Model.Conversation{}).Where("id = ?", conversationID).
			Updates(map[string]interface{}{"name": name, "image_url": imageURL})
		if result.Error != nil {
			return result.Error
		}
		return nil
	}
}

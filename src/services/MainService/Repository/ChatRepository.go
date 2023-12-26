package Repository

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ChatRepository struct {
	tx *gorm.DB
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB, tx *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db, tx: tx}
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
		Preload("OwnConversations.Messages").
		Preload("MemberConversations.Messages").
		Where("id = ?", userId).
		First(&userConversations)
	if result.Error != nil {
		return &Model.User{}, result.Error
	}
	return userConversations, nil
}

func (r *ChatRepository) SaveMessage(messageID int64, conversationId uint, senderId uint, message string,
	mType string, receiverId int, time time.Time, status string) (*Model.Message, error) {
	if r.tx.Error != nil {
		return nil, r.tx.Error
	}
	defer func() {
		if re := recover(); re != nil {
			r.tx.Rollback()
		}
	}()
	lastSeqNo := 0
	var conversation Model.Conversation
	r.tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", conversationId).First(&conversation)
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
	if err := r.tx.Save(&conversation).Error; err != nil {
		fmt.Println(err)
		r.tx.Rollback()
		return &Model.Message{}, err
	}
	result := r.tx.Create(&messageModel)
	if result.Error != nil {
		fmt.Println(result.Error)
		r.tx.Rollback()
		return &Model.Message{}, result.Error
	}
	err := r.tx.Commit().Error
	if err != nil {
		r.tx.Rollback()
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

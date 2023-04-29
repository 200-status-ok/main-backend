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

func (r *ChatRepository) CreateChatRoom(roomModel Model.ChatRoom) error {
	result := r.db.Create(&roomModel)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

package DBConfiguration

import (
	"github.com/403-access-denied/main-backend/src/Infrastructure/Utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// singleton pattern
var db *gorm.DB

func connectDB(connection string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(connection), &gorm.Config{})
}

func InitDB() {
	connection, _ := Utils.ReadFromEnvFile(".env", "DATABASE_URL")
	db, _ = connectDB(connection)
}

func MigrateModel(model interface{}) {
	db.AutoMigrate(model)
}

func GetDB() *gorm.DB {
	return db
}

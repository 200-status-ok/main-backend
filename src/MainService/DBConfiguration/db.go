package DBConfiguration

import (
	"github.com/403-access-denied/main-backend/src/MainService/Utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func connectDB(connection string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(connection), &gorm.Config{})
}

func InitDB() {
	appEnv := os.Getenv("APP_ENV2")
	if appEnv == "development" {
		connection := Utils.ReadFromEnvFile(".env", "LOCAL_DATABASE_URL")
		db, _ = connectDB(connection)
	} else if appEnv == "production" {
		connection := Utils.ReadFromEnvFile(".env", "PRODUCTION_DATABASE_URL")
		db, _ = connectDB(connection)
	}
}

func MigrateModel(model []interface{}) error {
	db.Exec(`DO $$
	             BEGIN
					IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
						CREATE TYPE status AS ENUM ('found', 'lost');
					END IF;
				END$$;`)
	for _, m := range model {
		err := db.AutoMigrate(m)
		if err != nil {
			return err
		}
	}
	return nil
}

func DropModel(model []interface{}) error {
	for _, m := range model {
		err := db.Migrator().DropTable(m)
		if err != nil {
			return err
		}
	}
	return nil
}

func CloseDB() {
	dbSQL, _ := db.DB()
	dbSQL.Close()
}

func GetDB() *gorm.DB {
	InitDB()
	return db
}

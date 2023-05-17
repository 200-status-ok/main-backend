package DBConfiguration

import (
	"github.com/403-access-denied/main-backend/src/MainService/Utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

var DB *gorm.DB

func connectDB(connection string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(connection), &gorm.Config{})
}

func init() {
	appEnv := os.Getenv("APP_ENV2")
	if appEnv == "development" {
		connection := Utils.ReadFromEnvFile(".env", "LOCAL_DATABASE_URL")
		DB, _ = connectDB(connection)
		dbSQL, _ := DB.DB()
		dbSQL.SetMaxIdleConns(10)
		dbSQL.SetMaxOpenConns(100)
		dbSQL.SetConnMaxLifetime(time.Hour)
	} else if appEnv == "production" {
		connection := Utils.ReadFromEnvFile(".env", "PRODUCTION_DATABASE_URL")
		DB, _ = connectDB(connection)
	}
}

func MigrateModel(model []interface{}) error {
	DB.Exec(`DO $$
	             BEGIN
					IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
						CREATE TYPE status AS ENUM ('found', 'lost');
					END IF;
				END$$;`)
	for _, m := range model {
		err := DB.AutoMigrate(m)
		if err != nil {
			return err
		}
	}
	return nil
}

func DropModel(model []interface{}) error {
	for _, m := range model {
		err := DB.Migrator().DropTable(m)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

package pgsql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

var db *gorm.DB
var tx *gorm.DB

// Connects to the PostgresSQL database using the provided connection string.
func connectDB(connection string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(connection), &gorm.Config{})
}

func setupDatabase(appEnv string) {
	var pgConnection string
	pgConnection = "postgresql://root:root@database:5432/main?sslmode=disable"
	if appEnv == "production" {
		pgConnection = "postgresql://root:a4bdJh8NnWY8AFCbKkfwnUu0@main-db:5432/postgres"
	}
	db, _ = connectDB(pgConnection)
	dbSQL, _ := db.DB()
	dbSQL.SetMaxIdleConns(10)
	dbSQL.SetMaxOpenConns(100)
	dbSQL.SetConnMaxLifetime(time.Hour)
}

func init() {
	appEnv := os.Getenv("APP_ENV2")
	setupDatabase(appEnv)
}

// MigrateModel migrates the specified models to the database.
func MigrateModel(models []interface{}) error {
	db.Exec(`DO $$
	             BEGIN
					IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
						CREATE TYPE status AS ENUM ('found', 'lost');
					END IF;
				END$$;`)
	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}

// DropModel drops the tables associated with the specified models.
func DropModel(models []interface{}) error {
	for _, model := range models {
		err := db.Migrator().DropTable(model)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetDB returns the GORM database instance.
func GetDB() *gorm.DB {
	return db
}

// GetTx returns the GORM database transaction instance.
func GetTx() *gorm.DB {
	return tx
}

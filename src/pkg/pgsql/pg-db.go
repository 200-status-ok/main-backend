package pgsql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

var dB *gorm.DB

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
	dB, _ = connectDB(pgConnection)
	dbSQL, _ := dB.DB()
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
	dB.Exec(`DO $$
	             BEGIN
					IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
						CREATE TYPE status AS ENUM ('found', 'lost');
					END IF;
				END$$;`)
	for _, model := range models {
		err := dB.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}

// DropModel drops the tables associated with the specified models.
func DropModel(models []interface{}) error {
	for _, model := range models {
		err := dB.Migrator().DropTable(model)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetDB returns the GORM database instance.
func GetDB() *gorm.DB {
	return dB
}

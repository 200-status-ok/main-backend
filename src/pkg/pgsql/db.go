package pgsql

import (
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

var DB *gorm.DB
var ElasticClient *elastic.Client

// Connects to the PostgresSQL database using the provided connection string.
func connectDB(connection string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(connection), &gorm.Config{})
}

// Connects to the Elasticsearch client using the provided connection string.
func connectElastic(connection string) (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{
			connection,
		},
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func setupDatabase(appEnv string) {
	pgKey := "LOCAL_DATABASE_URL"
	esKey := "LOCAL_ELASTIC_URL"

	if appEnv == "production" {
		pgKey = "PRODUCTION_DATABASE_URL"
		esKey = "PRODUCTION_ELASTIC_URL"
	}

	pgConnection := utils.ReadFromEnvFile(".env", pgKey)
	esConnection := utils.ReadFromEnvFile(".env", esKey)

	DB, _ = connectDB(pgConnection)
	ElasticClient, _ = connectElastic(esConnection)

	dbSQL, _ := DB.DB()
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
	DB.Exec(`DO $$
	             BEGIN
					IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
						CREATE TYPE status AS ENUM ('found', 'lost');
					END IF;
				END$$;`)
	for _, model := range models {
		err := DB.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}

// DropModel drops the tables associated with the specified models.
func DropModel(models []interface{}) error {
	for _, model := range models {
		err := DB.Migrator().DropTable(model)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetDB returns the GORM database instance.
func GetDB() *gorm.DB {
	return DB
}

// GetElastic returns the Elasticsearch client instance.
func GetElastic() *elastic.Client {
	return ElasticClient
}

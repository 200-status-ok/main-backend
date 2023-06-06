package DBConfiguration

import (
	"github.com/403-access-denied/main-backend/src/WorkerService/Utils"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

var DB *gorm.DB
var ElasticClient *elastic.Client

func connectDB(connection string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(connection), &gorm.Config{})
}

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

func init() {
	appEnv := os.Getenv("APP_ENV3")
	if appEnv == "development" {
		connection := Utils.ReadFromEnvFile(".env", "LOCAL_DATABASE_URL")
		esConnection := Utils.ReadFromEnvFile(".env", "LOCAL_ELASTIC_URL")
		DB, _ = connectDB(connection)
		ElasticClient, _ = connectElastic(esConnection)
		dbSQL, _ := DB.DB()
		dbSQL.SetMaxIdleConns(10)
		dbSQL.SetMaxOpenConns(100)
		dbSQL.SetConnMaxLifetime(time.Hour)
	} else if appEnv == "production" {
		connection := Utils.ReadFromEnvFile(".env", "PRODUCTION_DATABASE_URL")
		esConnection := Utils.ReadFromEnvFile(".env", "PRODUCTION_ELASTIC_URL")
		DB, _ = connectDB(connection)
		ElasticClient, _ = connectElastic(esConnection)
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

func GetElastic() *elastic.Client {
	return ElasticClient
}

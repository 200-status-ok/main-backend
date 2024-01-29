package UseCase

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/WorkerService/MessageCli"
	"github.com/200-status-ok/main-backend/src/pkg/pgsql"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

func CheckPhotoNSFW() {
	messageBroker := MessageCli.MessageClient{}
	var connectionString string
	appEnv := os.Getenv("APP_ENV2")
	if appEnv == "development" {
		connectionString = utils.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION")
		err := messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION"))
		if err != nil {
			panic(err)
		}
	} else if appEnv == "production" {
		connectionString = utils.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION")
		err := messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION"))
		if err != nil {
			panic(err)
		}
	}

	closeCh := make(chan *amqp.Error)
	messageBroker.Connection.NotifyClose(closeCh)

	go SendHeartbeat(messageBroker.Connection, &messageBroker, connectionString)
	db := pgsql.GetDB()

	go func() {
		err := messageBroker.SubscribeOnQueue("nsfw-validation", "nsfw-validation", db)
		if err != nil {
			fmt.Println("Error subscribing to nsfw-validation queue: ", err)
		}
	}()
}

func CheckTagNSFW() {
	messageBroker := MessageCli.MessageClient{}
	var connectionString string
	appEnv := os.Getenv("APP_ENV2")
	if appEnv == "development" {
		connectionString = utils.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION")
		err := messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION"))
		if err != nil {
			panic(err)
		}
	} else if appEnv == "production" {
		connectionString = utils.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION")
		err := messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION"))
		if err != nil {
			panic(err)
		}
	}

	closeCh := make(chan *amqp.Error)
	messageBroker.Connection.NotifyClose(closeCh)

	go SendHeartbeat(messageBroker.Connection, &messageBroker, connectionString)
	db := pgsql.GetDB()

	go func() {
		err := messageBroker.SubscribeOnQueue("tag-validation", "tag-validation", db)
		if err != nil {
			fmt.Println("Error subscribing to tag-validation queue: ", err)
		}
	}()
}

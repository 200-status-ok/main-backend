package UseCase

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/WorkerService/DBConfiguration"
	"github.com/200-status-ok/main-backend/src/WorkerService/MessageCli"
	Utils2 "github.com/200-status-ok/main-backend/src/WorkerService/Utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

func CheckPhotoNSFW() {
	messageBroker := MessageCli.MessageClient{}
	var connectionString string
	appEnv := os.Getenv("APP_ENV3")
	if appEnv == "development" {
		connectionString = Utils2.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION")
		err := messageBroker.ConnectBroker(Utils2.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION"))
		if err != nil {
			panic(err)
		}
	} else if appEnv == "production" {
		connectionString = Utils2.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION")
		err := messageBroker.ConnectBroker(Utils2.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION"))
		if err != nil {
			panic(err)
		}
	}

	closeCh := make(chan *amqp.Error)
	messageBroker.Connection.NotifyClose(closeCh)

	go SendHeartbeat(messageBroker.Connection, &messageBroker, connectionString)

	go func() {
		err := messageBroker.SubscribeOnQueue("nsfw-validation", "nsfw-validation", DBConfiguration.GetDB())
		if err != nil {
			fmt.Println("Error subscribing to nsfw-validation queue: ", err)
		}
	}()
}

func CheckTagNSFW() {
	messageBroker := MessageCli.MessageClient{}
	var connectionString string
	appEnv := os.Getenv("APP_ENV3")
	if appEnv == "development" {
		connectionString = Utils2.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION")
		err := messageBroker.ConnectBroker(Utils2.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION"))
		if err != nil {
			panic(err)
		}
	} else if appEnv == "production" {
		connectionString = Utils2.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION")
		err := messageBroker.ConnectBroker(Utils2.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION"))
		if err != nil {
			panic(err)
		}
	}

	closeCh := make(chan *amqp.Error)
	messageBroker.Connection.NotifyClose(closeCh)

	go SendHeartbeat(messageBroker.Connection, &messageBroker, connectionString)

	go func() {
		err := messageBroker.SubscribeOnQueue("tag-validation", "tag-validation", DBConfiguration.GetDB())
		if err != nil {
			fmt.Println("Error subscribing to tag-validation queue: ", err)
		}
	}()
}

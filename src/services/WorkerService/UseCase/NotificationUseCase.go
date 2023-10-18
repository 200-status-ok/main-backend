package UseCase

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/WorkerService/DBConfiguration"
	"github.com/200-status-ok/main-backend/src/WorkerService/MessageCli"
	Utils2 "github.com/200-status-ok/main-backend/src/WorkerService/Utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"time"
)

func SendToUser() {
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
		err := messageBroker.SubscribeOnQueue("email_notification", "email_notification", DBConfiguration.GetDB())
		if err != nil {
			fmt.Println("Error subscribing to email_notification queue:", err)
		}
	}()

	go func() {
		err := messageBroker.SubscribeOnQueue("sms_notification", "sms_notification", DBConfiguration.GetDB())
		if err != nil {
			fmt.Println("Error subscribing to sms_notification queue:", err)
		}
	}()
}

func SendHeartbeat(conn *amqp.Connection, mBroker *MessageCli.MessageClient, connString string) {
	heartbeatInterval := 90 * time.Second
	hearBeatTicker := time.NewTicker(heartbeatInterval)

	for range hearBeatTicker.C {
		if conn.IsClosed() {
			fmt.Println("RabbitMQ connection closed")
			err := mBroker.ConnectBroker(connString)
			if err != nil {
				fmt.Println("Error connecting to RabbitMQ:", err)
			}
			fmt.Println("Reconnecting to RabbitMQ")
			return
		}

		fmt.Println("Sending heartbeat")
	}

}

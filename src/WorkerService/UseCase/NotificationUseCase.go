package UseCase

import (
	"fmt"
	Utils2 "github.com/403-access-denied/main-backend/src/WorkerService/Utils"
	"os"
)

func SendToUser() {
	messageBroker := Utils2.MessageClient{}
	appEnv := os.Getenv("APP_ENV3")
	if appEnv == "development" {
		err := messageBroker.ConnectBroker(Utils2.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION"))
		if err != nil {
			panic(err)
		}
	} else if appEnv == "production" {
		err := messageBroker.ConnectBroker(Utils2.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION"))
		if err != nil {
			panic(err)
		}
	}

	go func() {
		err := messageBroker.SubscribeOnQueue("email_notification", "email_notification")
		if err != nil {
			fmt.Println("Error subscribing to email_notification queue:", err)
		}
	}()

	go func() {
		err := messageBroker.SubscribeOnQueue("sms_notification", "sms_notification")
		if err != nil {
			fmt.Println("Error subscribing to sms_notification queue:", err)
		}
	}()
}

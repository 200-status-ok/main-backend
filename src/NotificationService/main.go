package main

import "github.com/403-access-denied/main-backend/src/Utils"

func main() {
	messageBroker := Utils.MessageClient{}
	err := messageBroker.ConnectBroker(Utils.ReadFromEnvFile(".env", "RABBITMQ_DEFAULT_CONNECTION"))
	if err != nil {
		panic(err)
	}

	err = messageBroker.SubscribeOnQueue("notification", "notification")
	messageBroker.Close()

}

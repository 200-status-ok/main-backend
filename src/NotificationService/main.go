package main

import (
	"fmt"
	"github.com/403-access-denied/main-backend/src/NotificationService/src/Utils"
	"sync"
)

func main() {
	messageBroker := Utils.MessageClient{}
	err := messageBroker.ConnectBroker(Utils.ReadFromEnvFile(".env", "RABBITMQ_DEFAULT_CONNECTION"))
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := messageBroker.SubscribeOnQueue("email_notification", "email_notification")
		if err != nil {
			fmt.Println("Error subscribing to email_notification queue:", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := messageBroker.SubscribeOnQueue("sms_notification", "sms_notification")
		if err != nil {
			fmt.Println("Error subscribing to sms_notification queue:", err)
		}
	}()
	wg.Wait()
	messageBroker.Close()
}

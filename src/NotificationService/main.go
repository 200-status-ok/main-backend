package main

func main() {
	//messageBroker := Utils2.MessageClient{}
	//err := messageBroker.ConnectBroker(Utils2.ReadFromEnvFile(".env", "RABBITMQ_DEFAULT_CONNECTION"))
	//if err != nil {
	//	panic(err)
	//}
	//
	//var wg sync.WaitGroup
	//wg.Add(2)
	//go func() {
	//	defer wg.Done()
	//	err := messageBroker.SubscribeOnQueue("email_notification", "email_notification")
	//	if err != nil {
	//		fmt.Println("Error subscribing to email_notification queue:", err)
	//	}
	//}()
	//go func() {
	//	defer wg.Done()
	//	err := messageBroker.SubscribeOnQueue("sms_notification", "sms_notification")
	//	if err != nil {
	//		fmt.Println("Error subscribing to sms_notification queue:", err)
	//	}
	//}()
	//wg.Wait()
	//messageBroker.Close()
}

package Utils

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

type MessageClient struct {
	Connection *amqp.Connection
}

func (client *MessageClient) ConnectBroker(connectionString string) error {
	if connectionString == "" {
		return errors.New("connectionString is empty")
	}
	var err error
	client.Connection, err = amqp.Dial(connectionString)
	if err != nil {
		return err
	}
	return nil
}

func (client *MessageClient) Subscribe(exchangeName string, exchangeType string, consumerName string) error {
	channel, err := client.Connection.Channel()
	if err != nil {
		return err
	}
	err = channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.QueueBind(
		queue.Name,
		exchangeName,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.Qos(1, 0, false)
	if err != nil {
		return err
	}
	msgs, err := channel.Consume(
		queue.Name,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	// TODO: refactor this
	var forever chan struct{}
	go func() {
		for d := range msgs {
			fmt.Println("Received a message: ", string(d.Body))
		}
	}()
	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

type CustomArray []string

func (client *MessageClient) SubscribeOnQueue(queueName string, consumerName string) error {
	if client.Connection == nil {
		return errors.New("connection is nil")
	}
	channel, err := client.Connection.Channel()
	if err != nil {
		return err
	}
	queue, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.Qos(1, 0, false)
	if err != nil {
		return err
	}
	message, err := channel.Consume(
		queue.Name,
		consumerName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	var forever = make(chan struct{})
	for d := range message {
		fmt.Println("Received a message: ", string(d.Body))
		arr := CustomArray{}
		arr = strings.Split(string(d.Body), "/")
		arr.SendingNotification()
	}
	<-forever
	return nil
}

func (client *MessageClient) Close() {
	if client.Connection != nil {
		client.Connection.Close()
	}
}

func (a CustomArray) SendingNotification() {
	if a[0] == "email" {
		emailService := NewEmail("mhmdrzsmip@gmail.com", a[2],
			"Sending OTP code", "کد تایید ورود به سامانه همینجا: "+a[1],
			ReadFromEnvFile(".env", "GOOGLE_SECRET"))
		err := emailService.SendEmailWithGoogle()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if a[0] == "sms" {
		pattern := map[string]string{
			"code": a[1],
		}
		otpSms := NewSMS(ReadFromEnvFile(".env", "API_KEY"), pattern)
		err := otpSms.SendSMSWithPattern(a[2], ReadFromEnvFile(".env", "OTP_PATTERN_CODE"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

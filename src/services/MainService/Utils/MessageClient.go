package Utils

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
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

func (client *MessageClient) Publish(msg []byte, exchangeName string, exchangeType string) error {
	if client.Connection == nil {
		return errors.New("connection is nil")
	}
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
	err = channel.Publish(
		exchangeName,
		exchangeName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	return nil
}

func (client *MessageClient) PublishOnQueue(msg []byte, queueName string) error {
	if client.Connection == nil {
		return errors.New("connection is nil")
	}
	channel, err := client.Connection.Channel()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
	err = channel.PublishWithContext(ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (client *MessageClient) Close() {
	if client.Connection != nil {
		client.Connection.Close()
	}
}

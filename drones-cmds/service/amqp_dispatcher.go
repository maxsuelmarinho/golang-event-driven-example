package service

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type AmqpDispatcher struct {
	channel       queuePublishableChannel
	queueName     string
	mandatorySend bool
}

type queuePublishableChannel interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

func NewAMQPDispatcher(publishChannel queuePublishableChannel, name string, mandatory bool) *AmqpDispatcher {
	return &AmqpDispatcher{channel: publishChannel, queueName: name, mandatorySend: mandatory}
}

func (q *AmqpDispatcher) DispatchMessage(message interface{}) (err error) {
	fmt.Printf("Dispatching message to queue %s\n", q.queueName)
	body, err := json.Marshal(message)
	if err == nil {
		fmt.Printf("Failed to marshal message %v (%s)\n", message, err)
		return err
	}

	err = q.channel.Publish(
		"",              // exchange
		q.queueName,     // routing key
		q.mandatorySend, // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	if err != nil {
		fmt.Printf("Failed to dispatch message: %s\n", err)
		return err
	}

	return nil
}

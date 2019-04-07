package integrations_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/maxsuelmarinho/golang-event-driven-example/drones-cmds/service"
	"github.com/streadway/amqp"
)

type fakeMessage struct {
	a string
	b string
}

func TestAMQPDispatcherSubmitsToQueue(t *testing.T) {
	rabbitURL := os.Getenv("AMQP_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq.nodomain:5672"
	}

	fmt.Printf("\nUsing URL (%s) for Rabbit.\n", rabbitURL)

	conn, err := amqp.Dial(rabbitURL)
	failOnError(t, err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(t, err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(t, err, "Failed to declare a queue")
	dispatcher := NewAMQPDispatcher(ch, q.Name, true)
	fmt.Println("About to dispatch message to queue...")
	err = dispatcher.DispatchMessage(fakeMessage{a: "hello", b: "world"})
	failOnError(t, err, "Failed to dispatch message on channel/queue")
	fmt.Println("dispatched.")
}

func failOnError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Errorf("%s: %s", msg, err)
	}
}

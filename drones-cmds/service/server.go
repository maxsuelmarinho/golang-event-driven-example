package service

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/maxsuelmarinho/golang-event-driven-example/drones-cmds/fakes"
	"github.com/streadway/amqp"
	"github.com/unrolled/render"
)

func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	positionDispatcher := buildDispatcher("positions")
	telemetryDispatcher := buildDispatcher("telemetry")
	alertDispatcher := buildDispatcher("alerts")

	initRoutes(mx, formatter, telemetryDispatcher, alertDispatcher, positionDispatcher)

	n.UseHandler(mx)
	return n
}

func buildDispatcher(queueName string) queueDispatcher {
	url := resolveAMQPURL()
	if strings.Compare(url, "fake://foo") == 0 {
		fmt.Printf("Building fake dispatcher for queue '%s'", queueName)
		return fakes.NewFakeQueueDispatcher()
	}
	return createAMQPDispatcher(queueName, url)
}

func initRoutes(mx *mux.Router, formatter *render.Render, telemetryDispatcher queueDispatcher, alertDispatcher queueDispatcher, positionDispatcher queueDispatcher) {
	mx.HandleFunc("/api/cmds/telemetry", addTelemetryHandler(formatter, telemetryDispatcher)).Methods("POST")
	mx.HandleFunc("/api/cmds/alerts", addAlertHandler(formatter, alertDispatcher)).Methods("POST")
	mx.HandleFunc("/api/cmds/positions", addPositionHandler(formatter, positionDispatcher)).Methods("POST")
}

func resolveAMQPURL() string {
	url := os.Getenv("AMQP_URL")
	if url == "" {
		fmt.Printf("Failed to detect URL for bound rabbit service. Falling back to in-memory fake.\n")
		return "fake://foo"
	}

	return url
}

func createAMQPDispatcher(queueName string, url string) queueDispatcher {
	fmt.Printf("\nUsing URL (%s) for Rabbit.\n", url)

	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare a queue")
	dispatcher := NewAMQPDispatcher(ch, q.Name, false)
	return dispatcher
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

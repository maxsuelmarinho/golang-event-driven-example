package integrations_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/maxsuelmarinho/golang-event-driven-example/drones-cmds/service"
	dronescommon "github.com/maxsuelmarinho/golang-event-driven-example/drones-common"
	"github.com/streadway/amqp"
)

var (
	rabbitURLConfigured = os.Setenv("AMQP_URL", "amqp://guest:guest@rabbitmq.nodomain:5672")
	server              = NewServer()
	telemetry1          = []byte("{\"drone_id\": \"abc1234\", \"battery\": 80, \"uptime\": 3200, \"code_temp\": 20}")
	telemetry2          = []byte("{\"drone_id\": \"drone2\", \"battery\": 40, \"uptime\": 1200, \"code_temp\": 10}")
	doneTelemetry       = make(chan error)
	telemetryCount      = 0

	alert1     = []byte("{\"drone_id\": \"abc1234\", \"fault_code\": 1, \"description\": \"super fail\"}")
	alert2     = []byte("{\"drone_id\": \"drone2\", \"fault_code\": 2, \"description\": \"overheating\"}")
	doneAlert  = make(chan error)
	alertCount = 0

	position1     = []byte("{\"drone_id\": \"abc1234\", \"latitude\": 31.01, \"longitude\": 72.5, \"altitude\": 3500.12, \"current_speed\": 15.12, \"heading_cardinal\": 0}")
	position2     = []byte("{\"drone_id\": \"drone2\", \"latitude\": 31.01, \"longitude\": 72.5, \"altitude\": 3500.12, \"current_speed\": 15.12, \"heading_cardinal\": 0}")
	position3     = []byte("{\"drone_id\": \"drone3\", \"latitude\": 31.01, \"longitude\": 72.5, \"altitude\": 3500.12, \"current_speed\": 15.12, \"heading_cardinal\": 0}")
	donePosition  = make(chan error)
	positionCount = 0
)

func TestIntegration(t *testing.T) {

	fmt.Println("== Integration Test Scenario ==")

	telemetryReply, err := submitTelemetry(t, telemetry1)
	if err != nil {
		t.Errorf("Failed to submit telemetry: %s\n", err)
		return
	}

	if telemetryReply.DroneID != "abc1234" {
		t.Errorf("Failed to get a matching reply from the command server when submitting telemetry command: %+v\n", telemetryReply)
		return
	}

	telemetryReply2, err := submitTelemetry(t, telemetry2)
	if err != nil {
		t.Errorf("Failed to submit 2nd telemetry: %s\n", err)
		return
	}

	if telemetryReply2.DroneID != "drone2" {
		t.Errorf("Failed to get a matching reply from 2nd telemetry submit: %+v\n", telemetryReply2)
		return
	}

	alertReply, err := submitAlert(t, alert1)
	if err != nil {
		t.Errorf("Failed to submit an alert command: %s\n", err.Error())
		return
	}

	if alertReply.DroneID != "abc1234" {
		t.Errorf("Failed to get matching reply from submitting alert command: %+v\n", alertReply)
		return
	}

	alertReply2, err := submitAlert(t, alert2)
	if err != nil {
		t.Errorf("Failed to submit 2nd alert: %s\n", err.Error())
		return
	}

	if alertReply2.DroneID != "drone2" {
		t.Errorf("Expecting matching reply from submitting 2nd alert command, got %+v\n", alertReply2)
		return
	}

	positionReply, _ := submitPosition(t, position1)
	positionReply2, _ := submitPosition(t, position2)
	positionReply3, _ := submitPosition(t, position3)

	if positionReply.DroneID != "abc1234" {
		t.Errorf("Got wrong reply from position submit: %+v\n", positionReply)
		return
	}

	if positionReply2.DroneID != "drone2" {
		t.Errorf("Got wrong reply from position submit, expected drone2, got %+v\n", positionReply2)
	}

	if positionReply3.DroneID != "drone3" {
		t.Errorf("Got wrong reply from position submit, expected drone3, got %+v\n", positionReply3)
	}

	consumeRabbit(t)
	<-doneTelemetry
	<-doneAlert
	<-donePosition

	if telemetryCount != 2 {
		t.Errorf("Expected dequeue 2 telemetry events, got %d.\n", telemetryCount)
	}

	if alertCount != 2 {
		t.Errorf("Expected dequeue 2 alerts events, got %d.\n", alertCount)
	}

	if positionCount != 3 {
		t.Errorf("Expected dequeue 3 positions events, got %d.\n", positionCount)
	}
}

func consumeRabbit(t *testing.T) {
	rabbitURL := os.Getenv("AMQP_URL")

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		t.Errorf("Failed to dial rabbit: %s", err.Error())
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	defer ch.Close()

	telemetryQ, err := ch.QueueDeclare(
		"telemetry",
		false,
		false,
		false,
		false,
		nil,
	)

	alertsQ, err := ch.QueueDeclare(
		"alerts",
		false,
		false,
		false,
		false,
		nil,
	)

	positionsQ, err := ch.QueueDeclare(
		"positions",
		false,
		false,
		false,
		false,
		nil,
	)

	telemetryIn, err := ch.Consume(
		telemetryQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	alertsIn, err := ch.Consume(
		alertsQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	positionsIn, err := ch.Consume(
		positionsQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	/*
		go func() {
			for d := range telemetryIn {
				reactTelemetry(d)
			}
			doneTelemetry <- nil

			for a := range alertsIn {
				reactAlert(a)
			}
			doneAlert <- nil

			for p := range positionsIn {
				reactPosition(p)
			}
			donePosition <- nil
		}()
	*/

	go func() {
		for d := range telemetryIn {
			reactTelemetry(d)
		}
		doneTelemetry <- nil
	}()

	go func() {
		for a := range alertsIn {
			reactAlert(a)
		}
		doneAlert <- nil
	}()

	go func() {
		for p := range positionsIn {
			reactPosition(p)
		}
		donePosition <- nil
	}()

}

func reactTelemetry(telemetryRaw amqp.Delivery) {
	var event dronescommon.TelemetryUpdatedEvent
	err := json.Unmarshal(telemetryRaw.Body, &event)
	if err != nil {
		fmt.Printf("Failed to deserialize raw telemetry from queue, %v\n", err)
		return
	}

	fmt.Printf("Telemetry received: %+v\n", event)
	telemetryCount++
	telemetryRaw.Ack(false)
}

func reactAlert(alertRaw amqp.Delivery) {
	var event dronescommon.AlertSignalledEvent
	err := json.Unmarshal(alertRaw.Body, &event)
	if err != nil {
		fmt.Printf("Failed to deserialize raw alert from queue, %v\n", err)
		return
	}

	fmt.Printf("Alert received: %+v\n", event)
	alertCount++
	alertRaw.Ack(false)
}

func reactPosition(positionRaw amqp.Delivery) {
	var event dronescommon.PositionChangedEvent
	err := json.Unmarshal(positionRaw.Body, &event)
	if err != nil {
		fmt.Printf("Failed to deserialize raw position from queue, %v\n", err)
		return
	}

	fmt.Printf("Position received: %+v\n", event)
	positionCount++
	positionRaw.Ack(false)
}

func submitTelemetry(t *testing.T, body []byte) (reply dronescommon.TelemetryUpdatedEvent, err error) {
	rawReply, err := submitCommand(t, "/api/cmds/telemetry", body)
	var telemetryReply dronescommon.TelemetryUpdatedEvent
	if err != nil {
		t.Errorf("Failed to submit telemetry: %+v", err)
		return
	}

	err = json.Unmarshal(rawReply, &telemetryReply)

	if err != nil {
		t.Errorf("Failed to deserialize response: %s", err.Error())
		return
	}
	reply = telemetryReply
	return
}

func submitAlert(t *testing.T, body []byte) (reply dronescommon.AlertSignalledEvent, err error) {
	rawReply, err := submitCommand(t, "/api/cmds/alerts", body)
	var alertReply dronescommon.AlertSignalledEvent
	if err != nil {
		t.Errorf("Failed to submit alert: %+v", err)
		return
	}

	err = json.Unmarshal(rawReply, &alertReply)

	if err != nil {
		t.Errorf("Failed to deserialize response: %s", err.Error())
		return
	}
	reply = alertReply
	return
}

func submitPosition(t *testing.T, body []byte) (reply dronescommon.PositionChangedEvent, err error) {
	rawReply, err := submitCommand(t, "/api/cmds/positions", body)
	var positionReply dronescommon.PositionChangedEvent
	if err != nil {
		t.Errorf("Failed to submit position: %+v", err)
		return
	}

	err = json.Unmarshal(rawReply, &positionReply)

	if err != nil {
		t.Errorf("Failed to deserialize response: %s", err.Error())
		return
	}
	reply = positionReply
	return
}

func submitCommand(t *testing.T, url string, body []byte) (rawReply []byte, err error) {
	recorder := httptest.NewRecorder()
	commandRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return
	}

	server.ServeHTTP(recorder, commandRequest)
	if recorder.Code != 201 {
		errorMessage := "Error submitting command to %s, expected 201 got %d"
		t.Errorf(errorMessage, url, recorder.Code)
		err = fmt.Errorf(errorMessage, url, recorder.Code)
		return
	}
	rawReply = recorder.Body.Bytes()
	fmt.Printf("Command reply: HTTP %d %d bytes\n", recorder.Code, len(rawReply))
	return
}

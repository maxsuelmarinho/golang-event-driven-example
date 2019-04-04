package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	dronescommon "github.com/maxsuelmarinho/golang-event-driven-example/drones-common"

	"github.com/unrolled/render"
)

func addTelemetryHandler(formatter *render.Render, dispatcher queueDispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newTelemetryCommand telemetryCommand
		err := json.Unmarshal(payload, &newTelemetryCommand)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse add telemetry command.")
			return
		}

		if !newTelemetryCommand.isValid() {
			formatter.Text(w, http.StatusBadRequest, "Invalid telemetry command.")
			return
		}

		event := dronescommon.TelemetryUpdatedEvent{
			DroneID:          newTelemetryCommand.DroneID,
			RemainingBattery: newTelemetryCommand.RemainingBattery,
			Uptime:           newTelemetryCommand.Uptime,
			CoreTemp:         newTelemetryCommand.CoreTemp,
			ReceivedOn:       time.Now().Unix(),
		}
		fmt.Printf("Dispatching telemetry event for drone %s\n", newTelemetryCommand.DroneID)
		dispatcher.DispatchMessage(event)
		formatter.JSON(w, http.StatusCreated, event)
	}
}

func addAlertHandler(formatter *render.Render, dispatcher queueDispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newAlertCommand alertCommand
		err := json.Unmarshal(payload, &newAlertCommand)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse add alert command.")
			return
		}

		if !newAlertCommand.isValid() {
			formatter.Text(w, http.StatusBadRequest, "Invalid alert command.")
			return
		}

		event := dronescommon.AlertSignalledEvent{
			DroneID:     newAlertCommand.DroneID,
			FaultCode:   newAlertCommand.FaultCode,
			Description: newAlertCommand.Description,
			ReceivedOn:  time.Now().Unix(),
		}
		fmt.Printf("Dispatching alert event for drone %s\n", newAlertCommand.DroneID)
		dispatcher.DispatchMessage(event)
		formatter.JSON(w, http.StatusCreated, event)
	}
}

func addPositionHandler(formatter *render.Render, dispatcher queueDispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newPositionCommand positionCommand
		err := json.Unmarshal(payload, &newPositionCommand)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse add position command.")
			return
		}

		if !newPositionCommand.isValid() {
			formatter.Text(w, http.StatusBadRequest, "Invalid position command.")
			return
		}

		event := dronescommon.PositionChangedEvent{
			DroneID:         newPositionCommand.DroneID,
			Longitude:       newPositionCommand.Longitude,
			Latitude:        newPositionCommand.Latitude,
			Altitude:        newPositionCommand.Altitude,
			CurrentSpeed:    newPositionCommand.CurrentSpeed,
			HeadingCardinal: newPositionCommand.HeadingCardinal,
			ReceivedOn:      time.Now().Unix(),
		}
		fmt.Printf("Dispatching position event for drone %s\n", newPositionCommand.DroneID)
		dispatcher.DispatchMessage(event)
		formatter.JSON(w, http.StatusCreated, event)
	}
}

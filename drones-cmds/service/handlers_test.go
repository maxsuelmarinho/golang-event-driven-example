package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/maxsuelmarinho/golang-event-driven-example/drones-cmds/fakes"
	dronescommon "github.com/maxsuelmarinho/golang-event-driven-example/drones-common"
	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func makeTestServer(dispatcher queueDispatcher) *negroni.Negroni {
	server := negroni.New()
	mx := mux.NewRouter()
	initRoutes(mx, formatter, dispatcher, dispatcher, dispatcher)
	server.UseHandler(mx)
	return server
}

func TestAddValidTelemetryCreatesCommand(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	dispatcher := fakes.NewFakeQueueDispatcher()

	server := makeTestServer(dispatcher)
	recorder = httptest.NewRecorder()
	body := []byte("{\"drone_id\":\"drone123\",\"battery\":72,\"uptime\":6941,\"core_temp\":21}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/telemetry", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new telemetry item to return 201, got %d", recorder.Code)
	}

	if len(dispatcher.Messages) != 1 {
		t.Errorf("Expected queue dispatch count of 1, got %d", len(dispatcher.Messages))
	}

	var telemetryResponse dronescommon.TelemetryUpdatedEvent
	payload := recorder.Body.Bytes()
	err := json.Unmarshal(payload, &telemetryResponse)
	if err != nil {
		t.Errorf("Could not unmarshal payload into telemetry response object")
	}

	if telemetryResponse.DroneID != "drone123" {
		t.Errorf("Expected drone ID of 'drone123' got %s", telemetryResponse.DroneID)
	}

	if telemetryResponse.Uptime != 6941 {
		t.Errorf("Expected drone uptime of 6941, got %d", telemetryResponse.Uptime)
	}
}

func TestAddInvalidTelemetryReturnsBadRequest(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	dispatcher := fakes.NewFakeQueueDispatcher()
	server := makeTestServer(dispatcher)
	recorder = httptest.NewRecorder()
	body := []byte("{\"foo\":\"bar\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/telemetry", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected creation of invalid/unparseable new telemetry item to return bad request, got %d", recorder.Code)
	}

	if len(dispatcher.Messages) != 0 {
		t.Errorf("Expected dispatcher to dispatch 0 messages, got %d", len(dispatcher.Messages))
	}
}

func TestAddValidPositionCreatesCommand(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	dispatcher := fakes.NewFakeQueueDispatcher()

	server := makeTestServer(dispatcher)
	recorder = httptest.NewRecorder()
	body := []byte("{\"drone_id\":\"positiondrone123\",\"latitude\":81.231,\"longitude\":43.1231,\"altiture\":2301.1,\"current_speed\":41.3,\"heading_cardinal\":1}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/positions", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new position item to return 201, got %d/%s", recorder.Code, string(recorder.Body.Bytes()))
	}

	if len(dispatcher.Messages) != 1 {
		t.Errorf("Expected queue dispatch count of 1, got %d", len(dispatcher.Messages))
	}

	var positionResponse dronescommon.PositionChangedEvent
	payload := recorder.Body.Bytes()
	err := json.Unmarshal(payload, &positionResponse)
	if err != nil {
		t.Errorf("Could not unmarshal payload into position response object")
	}

	if positionResponse.DroneID != "positiondrone123" {
		t.Errorf("Expected drone ID of 'positiondrone123' got %s", positionResponse.DroneID)
	}

	if positionResponse.CurrentSpeed != 41.3 {
		t.Errorf("Expected drone current speed of 41.3, got %f", positionResponse.CurrentSpeed)
	}
}

func TestAddInvalidPositionReturnsBadRequest(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	dispatcher := fakes.NewFakeQueueDispatcher()
	server := makeTestServer(dispatcher)
	recorder = httptest.NewRecorder()
	body := []byte("{\"foo\":\"bar\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/positions", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected creation of invalid/unparseable new telemetry item to return bad request, got %d", recorder.Code)
	}

	if len(dispatcher.Messages) != 0 {
		t.Errorf("Expected dispatcher to dispatch 0 messages, got %d", len(dispatcher.Messages))
	}
}

func TestAddValidAlertCreatesCommand(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	dispatcher := fakes.NewFakeQueueDispatcher()

	server := makeTestServer(dispatcher)
	recorder = httptest.NewRecorder()
	body := []byte("{\"drone_id\":\"alertingdrone123\",\"fault_code\":12,\"description\":\"all the things are failing\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/alerts", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new position item to return 201, got %d/%s", recorder.Code, string(recorder.Body.Bytes()))
	}

	if len(dispatcher.Messages) != 1 {
		t.Errorf("Expected queue dispatch count of 1, got %d", len(dispatcher.Messages))
	}

	var alertResponse dronescommon.AlertSignalledEvent
	payload := recorder.Body.Bytes()
	err := json.Unmarshal(payload, &alertResponse)
	if err != nil {
		t.Errorf("Could not unmarshal payload into alert response object")
	}

	if alertResponse.DroneID != "alertingdrone123" {
		t.Errorf("Expected drone ID of 'alertingdrone123' got %s", alertResponse.DroneID)
	}

	if alertResponse.FaultCode != 12 {
		t.Errorf("Expected drone fault code of 12, got %d", alertResponse.FaultCode)
	}
}

func TestAddInvalidAlertReturnsBadRequest(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	dispatcher := fakes.NewFakeQueueDispatcher()
	server := makeTestServer(dispatcher)
	recorder = httptest.NewRecorder()
	body := []byte("{\"foo\":\"bar\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/api/cmds/alerts", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected creation of invalid/unparseable new telemetry item to return bad request, got %d", recorder.Code)
	}

	if len(dispatcher.Messages) != 0 {
		t.Errorf("Expected dispatcher to dispatch 0 messages, got %d", len(dispatcher.Messages))
	}
}

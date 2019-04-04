package service

import(
	"bytes"
	"encoding/json"
	"github.com/unrolled/render"
	"github.com/gorilla/mux"
	"github.com/codegangsta/negroni"
	"testing"
	"net/http"
	"net/http/httptest"
	dronescommon "github.com/maxsuelmarinho/golang-event-driven-example/drones-common"
)

var(
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
		request *http.Request
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
package dronecommon

type TelemetryUpdatedEvent struct {
	DroneID          string `json:"drone_id"`
	RemainingBattery int    `json:"battery"`
	Uptime           int    `json:"uptime"`
	CoreTemp         int    `json:"core_temp"`
	ReceivedOn       int64  `json:"received_on"`
}

type AlertSignalledEvent struct {
	DroneID     string `json:"drone_id"`
	FaultCode   int    `json:"fault_code"`
	Description string `json:"description"`
	ReceivedOn  int64  `json:"received_on"`
}

type PositionChangedEvent struct {
	DroneID         string  `json:"drone_id"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	Altitude        float32 `json:"altitude"`
	CurrentSpeed    float32 `json:"current_speed"`
	HeadingCardinal int     `json:"heading_cardinal"`
	ReceivedOn      int64   `json:"received_on"`
}

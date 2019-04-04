package service

import (
	"os"
	"strings"
	"testing"
)

func TestResolvesProperRabbitURL(t *testing.T) {
	os.Setenv("AMQP_URL", "amqp://guest:guest@rabbitmq:5672")

	url := resolveAMQPURL()
	if strings.Compare(url, "fake://foo") == 0 {
		t.Errorf("Got the fake URL when we should've gotten the proper URL.")
	}
}

func TestFallsBackToFakeURLWhenNoBoundService(t *testing.T) {
	os.Unsetenv("AMQP_URL")
	url := resolveAMQPURL()
	if strings.Compare(url, "fake://foo") != 0 {
		t.Errorf("Should have gotten the fake url, but didn't, got %s instead.", url)
	}
}

package transport

import "testing"

func TestNatsConnectionType(t *testing.T) {
	var _ Connection = (*natsConnection)(nil)
}

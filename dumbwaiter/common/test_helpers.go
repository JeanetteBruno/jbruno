package common

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockRPi local mock for RPi interface (can't use golang's mock - it has a race condition error)
type MockRPi struct {
	deviceName    string
	t             *testing.T
	ExpectedCalls []PiPin
	callCount     int

	latestPin   PiPin
	latestValue bool
	pinMu       sync.RWMutex
}

// NewMockRPi create a mock RPi object that validates it gets send signal calls as in expectedCalls
func NewMockRPi(t *testing.T, deviceName string, expectedCalls []PiPin) *MockRPi {
	return &MockRPi{t: t, deviceName: deviceName, ExpectedCalls: expectedCalls}
}

// SendSignal mock send signal- validate the calls match expectedCalls
func (m *MockRPi) SendSignal(pin PiPin) error {
	if m.callCount >= len(m.ExpectedCalls) {
		assert.Fail(m.t, fmt.Sprintf("%s expected %d calls, got an extra %s call", m.deviceName, len(m.ExpectedCalls), pin))
		return nil
	}
	assert.Equal(m.t, m.ExpectedCalls[m.callCount], pin, "%s expected %s, got %s", m.deviceName, m.ExpectedCalls[m.callCount], pin)
	m.callCount++

	m.pinMu.Lock()
	defer m.pinMu.Unlock()
	m.latestPin = pin
	m.latestValue = true

	return nil
}

// GetSignal mock get signal with noop
func (m *MockRPi) GetSignal(pin PiPin) (bool, error) {
	m.pinMu.RLock()
	defer m.pinMu.RUnlock()
	if m.latestPin == pin { // only return latestValue if it is for this pin
		return m.latestValue, nil
	}
	return false, nil
}

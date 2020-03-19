package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockRPi local mock for RPi interface (can't use golang's mock - it has a race condition error)
type MockRPi struct {
	t *testing.T
	expectedCalls []PiPin
	callCount int
}

// NewMockRPi create a mock RPi object that validates it gets send signal calls as in expectedCalls 
func NewMockRPi(t *testing.T, exepectedCalls []PiPin) *MockRPi {
	return &MockRPi{t: t, expectedCalls: exepectedCalls}
}

// SendSignal mock send signal- validate the calls match expectedCalls
func (m *MockRPi) SendSignal(pin PiPin) error {
	assert.Equal(m.t, m.expectedCalls[m.callCount], pin)
	m.callCount++
	return nil
}

// GetSignal mock get signal with noop
func (m *MockRPi) GetSignal(pin PiPin) (bool, error)  {
	return false, nil
}

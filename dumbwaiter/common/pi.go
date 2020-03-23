//go:generate mockgen -destination=./mocks/rpi_mock.go github.com/jbruno/dumbwaiter/common RPi
package common

/*
pi.go implements the interface with the RPi B+ machine
*/

// RPi the interface functions for getting/sending signal on the gpio pins
type RPi interface {
	SendSignal(pin PiPin) error        // send a signal on the target pin
	GetSignal(pin PiPin) (bool, error) // get a signal from a target pin, true when signal is on, false when off
}

type PiPin int

const (
	OpenerUp PiPin = iota
	OpenerDown
	OpenerStop
	Floor1Requested
	Floor2Requested
	Floor3Requested
	StopRequested
	AtFloor
)

func (p PiPin) String() string {
	return [...]string{"OpenerUp", "OpenerDown", "OpenerStop", "Floor1Requested", "Floor2Requested", "Floor3Requested", "StopRequested", "AtFloor"}[p]
}

// TODO implement the PI interfaces

// RPiDevice communicates with the RPi B+ device
type RPiDevice struct{}

// NewRPiDevice return an instance of a RPiDevice
func NewRPiDevice() *RPiDevice {
	return &RPiDevice{}
}

// SendSignal send a signal on the selected pin to the RPi B+ device
func (r *RPiDevice) SendSignal(pin PiPin) error {
	// TODO implement
	return nil
}

// GetSignal get a signal from the selected pin on the RPi B+ device
func (r *RPiDevice) GetSignal(pin PiPin) (bool, error) {
	// TODO implement
	return false, nil
}

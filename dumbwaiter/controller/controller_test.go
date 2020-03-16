package controller

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jbruno/dumbwaiter/common"
)

var testFrequency time.Duration = 500 * time.Millisecond

// TestRequestUpFromStop test requesting the car to move up 1 floor when it is stopped
func TestRequestUpFromStop(t *testing.T) {
	// setup
	timingCh := make(chan time.Time)
	dwController := setup(t, timingCh, 2, Stopped, []common.PiPin{common.OpenerUp})

	// test
	dwController.SetRequestedFloor(3)
	timingCh <- time.Now()
	waitForStatus(t, 2, 3, Up, dwController, 3 * time.Second)
}

// TestRequestStopFromUp test requesting the car to stop floor when it is moving up
func TestRequestArriveAtRequestedFloor(t *testing.T) {
	// setup
	timingCh := make(chan time.Time)
	dwController := setup(t, timingCh, 2, Up, []common.PiPin{common.OpenerStop})

	// test
	dwController.SetRequestedFloor(2)
	timingCh <- time.Now()
	waitForStatus(t, 2, 2, Stopped, dwController, 3 * time.Second)
}

// TestRequestUpFromMovingDown send an up request when the dumbwaiter is moving down,
// expect to see the dumbwaiter to to stop, then go to moving up
func TestRequestUpFromMovingDown(t *testing.T) {
	// setup
	timingCh := make(chan time.Time)
	dwController := setup(t, timingCh, 2, Down, []common.PiPin{common.OpenerStop, common.OpenerUp})

	// test
	dwController.SetRequestedFloor(3)            // send the up request
	timingCh <- time.Now()                       // trigger a control loop to process the request
	waitForStatus(t, 2, 3, Stopped, dwController, 3 * time.Second) // verify that the dumbwaiter is stopped

	timingCh <- time.Now()                  // trigger a control loop to process the request
	waitForStatus(t, 2, 3, Up, dwController, 3 * time.Second) // verify that the dumbwaiter is now moving up
}

// TestRequestDownFromMovingUp send an down request when the dumbwaiter is moving up,
// expect to see the dumbwaiter to to stop, then go to moving down
func TestRequestDownFromMovingUp(t *testing.T) {
	// setup
	timingCh := make(chan time.Time)
	dwController := setup(t, timingCh, 2, Up, []common.PiPin{common.OpenerStop, common.OpenerDown})

	// test
	dwController.SetRequestedFloor(1)            // send the down request
	timingCh <- time.Now()                       // trigger a control loop to process the request
	waitForStatus(t, 2, 1, Stopped, dwController, 3 * time.Second) // verify that the dumbwaiter is stopped

	timingCh <- time.Now()                    // trigger a control loop to process the request
	waitForStatus(t, 2, 1, Down, dwController, 3 * time.Second) // verify that the dumbwaiter is now moving down
}

// setup creates a controller, and mock pi interface
func setup(t *testing.T, timingCh chan time.Time, floor int, direction Direction, expectedSendSignals []common.PiPin) *Controller {
	mockRPi := common.NewMockRPi(t, expectedSendSignals)
	var dwController *Controller
	dwController = newController(3, mockRPi, testFrequency, timingCh)
	dwController.SetLastSeenFloor(floor)
	dwController.SetMovingDirection(direction)
	return dwController
}

func teardown(ticker *time.Ticker) {
	ticker.Stop()
}

// testFloorRequest given the last floor and direction the dumbwaiter is supposed to be moving (lastSeenFloor, currentDirection),
// and the new floor requested (requestedFloor), verify that the dumbwaiter's direction changes to expectedDirection
func testFloorRequest(t *testing.T, lastSeenFloor int, currentDirection Direction, requestedFloor int, expectedDirection Direction, dwc *Controller,
	timingCh chan time.Time) {
	dwc.SetRequestedFloor(requestedFloor)
	timingCh <- time.Now()
	waitForStatus(t, lastSeenFloor, requestedFloor, expectedDirection, dwc, 3 * time.Second)
}

func waitForStatus(t *testing.T, lastSeenFloor int, requestedFloor int, expectedDirection Direction, dwc *Controller, timeout time.Duration) {

	waitTill := time.Now().Add(timeout)
	for time.Now().Before(waitTill) {
		s := dwc.GetStatus()
		if expectedDirection == s.MovingDirection && lastSeenFloor == s.LastSeenFloor && requestedFloor == s.RequestedFloor {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	s := dwc.GetStatus()
	assert.Equal(t, expectedDirection.String(), s.MovingDirection.String(), "wrong direction")
	assert.Equal(t, requestedFloor, s.RequestedFloor, "wrong requested floor")
	assert.Equal(t, lastSeenFloor, s.LastSeenFloor, "wrong last seen floor")
}

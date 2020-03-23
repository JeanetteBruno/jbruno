/*
Package inttests contains integration tests (feature tests only - non-restian).
These tests start a controller and 3 floor sensor control loops.
The tests then orchestrate and test behavior by
controlling the fake pi GetSignal() values validating the pi SendSignal() values.
*/
package inttests

import (
	"fmt"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/jbruno/dumbwaiter/common"
	"github.com/jbruno/dumbwaiter/controller"
	floor_sensors "github.com/jbruno/dumbwaiter/floor"
)

// TestFloor3CallsToFloor3WhenControllerIsStoppedAt2 simulate the floor3 button is pressed on the
// third floor's control pad while the dumbwaiter is stopped at floor 2.  Verify that dumbwaiter
// gets and UP signal
func TestFloor3CallsToFloor3WhenControllerIsStoppedAt2(t *testing.T) {
	rpis, dwc, _ := setup(t, 2, controller.Stopped)
	rpis[0].ExpectedCalls = []common.PiPin{common.OpenerUp}
	rpis[3].ExpectedCalls = []common.PiPin{common.Floor3Requested}

	// trigger the up request
	log.Info("floor 3 is requesting the dumbwaiter to go to floor 3")
	rpis[3].SendSignal(common.Floor3Requested)

	waitForControllerStatus(t, 2, 3, controller.Up, dwc, 5*time.Second)

}

// setup creates a controller and sensors, each with their own mock pi interface
func setup(t *testing.T, floor int, direction controller.Direction) ([]*common.MockRPi, *controller.Controller, []*floor_sensors.Sensors) {
	testFrequency := 1 * time.Second // speed up tests

	var mockRPis [5]*common.MockRPi
	// set up the controller
	mockRPis[0] = common.NewMockRPi(t, "controllerRPi", nil)
	var dwController *controller.Controller
	dwController = controller.NewController(3).SetRPiDevice(mockRPis[0]).SetLoopFrequency(testFrequency)
	dwController.SetLastSeenFloor(floor)
	dwController.SetMovingDirection(direction)
	if direction == controller.Stopped {
		dwController.SetRequestedFloor(floor)
	} else if direction == controller.Up {
		dwController.SetRequestedFloor(floor + 1)
	} else {
		dwController.SetRequestedFloor(floor - 1)
	}
	dwController.StartProcessingLoop()

	var floors [4]*floor_sensors.Sensors
	// set up floor sensors
	for i := 1; i < 4; i++ {
		mockRPis[i] = common.NewMockRPi(t, fmt.Sprintf("floor%dRPi", i), nil)
		floors[i] = floor_sensors.NewSensors(i, "fakeURL").
			SetRPiDevice(mockRPis[i]).
			SetControllerCli(dwController).
			SetLoopFrequency(testFrequency)
		floors[i].StartProcessingLoop()
	}

	return mockRPis[:], dwController, floors[:]
}

func waitForControllerStatus(t *testing.T, lastSeenFloor int, requestedFloor int, expectedDirection controller.Direction, dwc *controller.Controller, timeout time.Duration) {
	waitTill := time.Now().Add(timeout)
	for time.Now().Before(waitTill) {
		s := dwc.GetStatus()
		if expectedDirection == s.MovingDirection && lastSeenFloor == s.LastSeenFloor && requestedFloor == s.RequestedFloor {
			return
		}
		time.Sleep(1 * time.Second)
	}
	s := dwc.GetStatus()
	assert.Equal(t, expectedDirection.String(), s.MovingDirection.String(), "timeout: wrong direction")
	assert.Equal(t, requestedFloor, s.RequestedFloor, "timeout: wrong requested floor")
	assert.Equal(t, lastSeenFloor, s.LastSeenFloor, "timeout: wrong last seen floor")
}

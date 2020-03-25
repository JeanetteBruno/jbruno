/*
Package inttests contains integration tests (feature tests only - non-restian).
These tests start a controller and 3 floor sensors.
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
	dwc, dwcRpi, _, rpis := setup(t, 2, controller.Stopped)
	dwcRpi.ExpectedCalls = []common.PiPin{common.OpenerUp}
	rpis[2].ExpectedCalls = []common.PiPin{common.Floor3Requested}

	// trigger the up request
	log.Info("floor 3 is requesting the dumbwaiter to go to floor 3")
	rpis[2].SendSignal(common.Floor3Requested)

	waitForControllerStatus(t, 2, 3, controller.Up, dwc, 5*time.Second)
}

// TestPlatformArrivesAtRequestedFloor
func TestPlatformArrivesAtRequestedFloor(t *testing.T) {
	dwc, dwcRpi, _, rpis := setup(t, 2, controller.Stopped)
	dwcRpi.ExpectedCalls = []common.PiPin{common.OpenerUp, common.OpenerStop}
	rpis[1].ExpectedCalls = []common.PiPin{common.Floor3Requested}
	rpis[2].ExpectedCalls = []common.PiPin{common.AtFloor}

	// trigger the up request
	log.Info("floor 2 is requesting the dumbwaiter to go to floor 3")
	rpis[1].SendSignal(common.Floor3Requested)

	go func() {
		select {
		case <-time.After(1 * time.Second):
			log.Info("floor 3 is telling controller the dumbwaiter has arrived")
			rpis[2].SendSignal(common.AtFloor)
		}
	}()

	waitForControllerStatus(t, 3, 3, controller.Stopped, dwc, 5*time.Second)

	t.Fail()
}

// setup creates a controller and sensors, each with their own mock pi interface
func setup(t *testing.T, floor int, direction controller.Direction) (*controller.Controller, *common.MockRPi, []*floor_sensors.Sensors, []*common.MockRPi) {
	testFrequency := 50 * time.Millisecond // speed up tests

	// set up the controller
	controlMockRPi := common.NewMockRPi(t, "controllerRPi", nil)
	var dwController *controller.Controller
	dwController = controller.NewController(3).SetRPiDevice(controlMockRPi).SetLoopFrequency(testFrequency)
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

	var mockRPis [3]*common.MockRPi
	var floors [3]*floor_sensors.Sensors
	// set up floor sensors
	for i := 0; i < 3; i++ {
		mockRPis[i] = common.NewMockRPi(t, fmt.Sprintf("floor%dRPi", i+1), nil)
		floors[i] = floor_sensors.NewSensors(i+1, "fakeURL").
			SetRPiDevice(mockRPis[i]).
			SetControllerCli(dwController).
			SetLoopFrequency(testFrequency)
		floors[i].StartProcessingLoop()
	}

	return dwController, controlMockRPi, floors[:], mockRPis[:]
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

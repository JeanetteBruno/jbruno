package controller

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/jbruno/dumbwaiter/common"
	"github.com/jbruno/dumbwaiter/common/mock_common"
)

// TestRequestUpFromStop test requesting the car to move up 1 floor when it is stopped
func TestRequestUpFromStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRPi := mock_common.NewMockRPi(ctrl)
	mockRPi.EXPECT().SendSignal(common.OpenerUp)  // expect to send a signal on the opener up pin

	dwController := newController(3, mockRPi, 200 * time.Millisecond)
	testFloorRequest(t, 2, Stopped, 3, Up, dwController)
}

// TestRequestDownFromStop test requesting the car to move up 1 floor when it is stopped
func TestRequestDownFromStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRPi := mock_common.NewMockRPi(ctrl)
	mockRPi.EXPECT().SendSignal(common.OpenerDown)  // expect to send a signal on the opener up pin

	dwController := newController(3, mockRPi, 200 * time.Millisecond)
	testFloorRequest(t, 2, Stopped, 1, Down, dwController)
}

func testFloorRequest(t *testing.T, lastSeenFloor int, currentDirection Direction, requestedFloor int, expectedDirection Direction, dwc *Controller) {
	dwc.SetLastSeenFloor(lastSeenFloor)  // initialize the dumbwaiter
	dwc.movingDirection = currentDirection
	dwc.SetRequestedFloor(requestedFloor)
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
package floor

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/JeanetteBruno/jbruno/dumbwaiter/common"
)

var testFrequency time.Duration = 500 * time.Millisecond
var noSignalEnd = time.Unix(int64(0), int64(0))
var foreverFalseSignal = fakeSignal{signalValue: false, signalEnd: noSignalEnd}
var foreverTrueSignal = fakeSignal{signalValue: true, signalEnd: noSignalEnd}

const (
	lsf = "lastSeenFloor"
	rf  = "requestedFloor"
	sr = "stopRequested"
)

type fakeSignal struct {
	signalValue bool
	signalEnd   time.Time
}

// fakePiDevice will return a signal value till the end of the signal time, or continuosly if the
// end time is 0
type fakePiDevice struct {
	signals map[common.PiPin]fakeSignal
}

func (f *fakePiDevice) GetSignal(pin common.PiPin) (bool, error) {
	if _, ok := f.signals[pin]; !ok {
		return false, nil
	}
	if f.signals[pin].signalEnd == noSignalEnd || time.Now().Before(f.signals[pin].signalEnd) {
		return f.signals[pin].signalValue, nil
	}
	return !f.signals[pin].signalValue, nil
}

func (f *fakePiDevice) SendSignal(pin common.PiPin) error {
	return nil
}

type controllerCall struct {
	callType  string
	callValue int
}

type validatingController struct {
	t                *testing.T
	expectedSequence []controllerCall
	currentSeqIndex  int

	lastSeenFloor  int
	requestedFloor int
}

func newvalidatingController(t *testing.T, expectedSequence []controllerCall) *validatingController {
	return &validatingController{t: t, expectedSequence: expectedSequence, currentSeqIndex: 0}
}

func (f *validatingController) SetRequestedFloor(floor int) {
	assert.True(f.t, len(f.expectedSequence) > f.currentSeqIndex,
		fmt.Sprintf("too many calls to SetLastSeenFloor, expected %d got %d", len(f.expectedSequence), f.currentSeqIndex))
	assert.True(f.t, rf == f.expectedSequence[f.currentSeqIndex].callType,
		fmt.Sprintf("wrong controller api call, expected %s, got %s (call:%d)",
			f.expectedSequence[f.currentSeqIndex].callType, rf, f.currentSeqIndex+1))
	assert.True(f.t, f.expectedSequence[f.currentSeqIndex].callValue == floor,
		fmt.Sprintf("wrong floor number, expected %d got %d (call:%d)",
			f.expectedSequence[f.currentSeqIndex].callValue, floor, f.currentSeqIndex+1))
	f.currentSeqIndex++
	f.requestedFloor = floor
}

func (f *validatingController) SetLastSeenFloor(floor int) {
	assert.True(f.t, len(f.expectedSequence) > f.currentSeqIndex,
		fmt.Sprintf("too many calls to SetLastSeenFloor, expected %d got %d", len(f.expectedSequence), f.currentSeqIndex))
	assert.True(f.t, lsf == f.expectedSequence[f.currentSeqIndex].callType, fmt.Sprintf("wrong controller api call, expected %s, got %s (call: %d)",
		f.expectedSequence[f.currentSeqIndex].callType, lsf, f.currentSeqIndex+1))
	assert.True(f.t, f.expectedSequence[f.currentSeqIndex].callValue == floor,
		fmt.Sprintf("wrong floor number, expected %d got %d (call:%d)",
			f.expectedSequence[f.currentSeqIndex].callValue, floor, f.currentSeqIndex+1))
	f.currentSeqIndex++
	f.lastSeenFloor = floor
}

func (f *validatingController) SetStopRequested() {
	assert.True(f.t, len(f.expectedSequence) > f.currentSeqIndex,
		fmt.Sprintf("too many calls to SetStopRequested, expected %d got %d", len(f.expectedSequence), f.currentSeqIndex))
	assert.True(f.t, sr == f.expectedSequence[f.currentSeqIndex].callType, fmt.Sprintf("wrong controller api call, expected %s, got %s (call: %d)",
		f.expectedSequence[f.currentSeqIndex].callType, sr, f.currentSeqIndex+1))
	f.currentSeqIndex++
}

func TestArriveAtFloor(t *testing.T) {
	// setup
	defaultLoopFrequency = 10 * time.Millisecond // speed up tests
	mockRPi := &fakePiDevice{}
	// create a controller that validates getting a request on SetLastSeenFloor() entry
	controllerClient := newvalidatingController(t, []controllerCall{controllerCall{callType: lsf, callValue: 1}})
	sensors := NewSensors(1, "fakeURL").SetRPiDevice(mockRPi).SetControllerClient(controllerClient).SetLoopFrequency(defaultLoopFrequency)
	sensors.StartProcessingLoop()
	signals := map[common.PiPin]fakeSignal{ // set up signalling at the floor
		common.Floor1Requested: foreverFalseSignal,
		common.Floor2Requested: foreverFalseSignal,
		common.Floor3Requested: foreverFalseSignal,
		common.AtFloor:         foreverTrueSignal,
		common.StopRequested:   foreverFalseSignal}

	// test
	mockRPi.signals = signals // this action triggers the test

	// final validation
	waitForStatus(t, 1, 0, controllerClient, 1*time.Second)
}

func TestPressFloor1Button(t *testing.T) {
	// setup
	defaultLoopFrequency = 10 * time.Millisecond // speed up tests
	mockRPi := &fakePiDevice{}
	// create a controller that validates getting a request on SetLastSeenFloor() entry
	controllerClient := newvalidatingController(t, []controllerCall{controllerCall{callType: rf, callValue: 1}})
	sensors := NewSensors(1, "fakeURL").SetRPiDevice(mockRPi).SetControllerClient(controllerClient).SetLoopFrequency(defaultLoopFrequency)
	sensors.StartProcessingLoop()
	signals := map[common.PiPin]fakeSignal{ // set up signalling at the floor
		common.Floor1Requested: foreverTrueSignal,
		common.Floor2Requested: foreverFalseSignal,
		common.Floor3Requested: foreverFalseSignal,
		common.AtFloor:         foreverFalseSignal,
		common.StopRequested:   foreverFalseSignal}

	// test
	mockRPi.signals = signals // this action triggers the test

	// final validation
	waitForStatus(t, 0, 1, controllerClient, 1*time.Second)
}

func TestPressStopButton(t *testing.T) {
	// setup
	defaultLoopFrequency = 10 * time.Millisecond // speed up tests
	mockRPi := &fakePiDevice{}
	// create a controller that validates getting a request on SetLastSeenFloor() entry
	controllerClient := newvalidatingController(t, []controllerCall{controllerCall{callType: sr}})
	sensors := NewSensors(1, "fakeURL").SetRPiDevice(mockRPi).SetControllerClient(controllerClient).SetLoopFrequency(defaultLoopFrequency)
	sensors.StartProcessingLoop()
	signals := map[common.PiPin]fakeSignal{ // set up signalling at the floor
		common.Floor1Requested: foreverFalseSignal,
		common.Floor2Requested: foreverFalseSignal,
		common.Floor3Requested: foreverFalseSignal,
		common.AtFloor:         foreverFalseSignal,
		common.StopRequested:   foreverTrueSignal}

	// test
	mockRPi.signals = signals // this action triggers the test

	// final validation
	waitForStatus(t, 0, 0, controllerClient, 1*time.Second)
}

func waitForStatus(t *testing.T, lastSeenFloor int, requestedFloor int, dwc *validatingController, timeout time.Duration) {
	waitTill := time.Now().Add(timeout)
	time.Sleep(500 * time.Millisecond)
	for time.Now().Before(waitTill) {
		if lastSeenFloor == dwc.lastSeenFloor && requestedFloor == dwc.requestedFloor {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	if requestedFloor != 0 {
		assert.Equal(t, requestedFloor, dwc.requestedFloor, fmt.Sprintf("wrong requested floor, expected %d got %d", requestedFloor, dwc.requestedFloor))
	}
	if lastSeenFloor != 0 {
		assert.Equal(t, lastSeenFloor, dwc.lastSeenFloor, fmt.Sprintf("wrong last seen floor, expected %d, got %d", lastSeenFloor, dwc.lastSeenFloor))
	}
}

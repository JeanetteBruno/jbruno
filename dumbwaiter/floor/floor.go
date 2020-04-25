package floor

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/JeanetteBruno/jbruno/dumbwaiter/common"
	"github.com/JeanetteBruno/jbruno/dumbwaiter/controller/api"
	"github.com/JeanetteBruno/jbruno/dumbwaiter/controller/api/cli"
)

var defaultLoopFrequency time.Duration = 500 * time.Millisecond

// Sensors monitor sensors at each floor and send requests to controller
type Sensors struct {
	floorNum           int
	selectedFloor      []int
	atFloorSensor      bool
	stopSelected       bool
	rpi                common.RPi
	mainLoopTicker     *time.Ticker
	loopFreq           time.Duration
	controllerClient   api.Controller
	controllerURL      string
	priorSelectedFloor int
	priorAtFloor       bool
}

// NewSensors create a new sensors object
func NewSensors(floorNum int, controllerURL string) *Sensors {
	return &Sensors{
		floorNum:         floorNum,
		rpi:              common.NewRPiDevice(),
		loopFreq:         defaultLoopFrequency,
		controllerURL:    controllerURL,
		controllerClient: cli.NewControllerHTTPClient(controllerURL),
		priorAtFloor:     false,
	}
}

// StartProcessingLoop start the processing loop in its own goroutine
func (s *Sensors) StartProcessingLoop() {
	go s.processingLoop()
}

// collect the values from the floor's sensors and send them to the controller
func (s *Sensors) processingLoop() {
	log.Infof("Starting floor%d sensor loop", s.floorNum)
	s.mainLoopTicker = time.NewTicker(s.loopFreq)

	for {
		select {
		case <-s.mainLoopTicker.C:
			s.handleAtFloorSensor()
			s.handleFloorRequestSensor(common.Floor1Requested, 1)
			s.handleFloorRequestSensor(common.Floor2Requested, 2)
			s.handleFloorRequestSensor(common.Floor3Requested, 3)
			s.handleStopRequestSensor(common.StopRequested)
		}
	}
}

// handleAtFloorSensor sends an atfloor request when the platform reaches this floor
func (s *Sensors) handleAtFloorSensor() {
	var sensor bool
	var err error
	//looking for sensor error
	sensor, err = s.rpi.GetSignal(common.AtFloor)
	if err != nil {
		log.Errorf("error getting AtFloor sensor: %e", err) // TODO implement real error handling
		return
	}
	if sensor && !s.priorAtFloor {
		log.Infof("sent at floor %d notice to controller", s.floorNum)
		s.controllerClient.SetLastSeenFloor(s.floorNum)
		s.priorAtFloor = true
	}
}

// handleFloorRequestSensor sends a new floor request to the controller
func (s *Sensors) handleFloorRequestSensor(pin common.PiPin, floorNum int) {
	var buttonPressed bool
	var err error
	if buttonPressed, err = s.rpi.GetSignal(pin); err != nil {
		log.Errorf("error getting Floor%d button: %e", floorNum, err) // TODO implement real error handling
	}
	if buttonPressed && floorNum != s.priorSelectedFloor {
		log.Infof("floor %d send call to floor %d to controller", s.floorNum, floorNum)
		s.controllerClient.SetRequestedFloor(floorNum)
		s.priorSelectedFloor = floorNum
	}
}

//handleStopRequestSensor sends a stop request to the controller
func (s *Sensors) handleStopRequestSensor(pin common.PiPin) {
	var buttonPressed bool
	var err error
	if buttonPressed, err = s.rpi.GetSignal(pin); err != nil {
		log.Errorf("error getting stop button: %e", err) // TODO implement real error handling
	}
	if buttonPressed && !s.stopSelected {
		log.Infof("send stop call to controller")

		//implement controller stop request
		s.controllerClient.SetStopRequested()
		s.stopSelected = true
	}
}

// Sensors constructor setters for builder pattern

// SetRPiDevice used by testing to override production RPi interface
func (s *Sensors) SetRPiDevice(rpiDevice common.RPi) *Sensors {
	s.rpi = rpiDevice
	return s
}

// SetLoopFrequency used by testing to speed up tests
func (s *Sensors) SetLoopFrequency(freq time.Duration) *Sensors {
	s.loopFreq = freq
	return s
}

// SetControllerClient set controller client used by testing
func (s *Sensors) SetControllerClient(controller api.Controller) *Sensors {
	s.controllerClient = controller
	return s
}

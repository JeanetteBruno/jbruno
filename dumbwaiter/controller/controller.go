package controller

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jbruno/dumbwaiter/common"
)

var defaultLoopFrequency time.Duration = time.Second

// Direction the direction the dumbwaiter car is moving
type Direction int

const (
	Up Direction = iota
	Down
	Stopped
)

func (d Direction) String() string {
	return [...]string{"up", "down", "stopped"}[d]
}

// Status returns the current status of the dumbwaiter
type Status struct {
	MovingDirection Direction
	RequestedFloor  int
	LastSeenFloor   int

	// TODO add array of floors' status
}

// Controller sends up, down, stop signals to the garage door opener based on getting
// control directives from the dumbwaiter floor services and/or the web app
type Controller struct {
	lastSeenFloor    int // the last floor reporting the car was seen at
	lastSeenFloorMU  sync.RWMutex
	requestedFloor   int // the car should move to this floor
	requestedFloorMU sync.RWMutex

	topFloor int // the top floor number (floor numbers start at 1)

	movingDirection   Direction // the direction the cab is currently moving
	movingDirectionMu sync.RWMutex

	timeStartedMoving  time.Time
	timeToMoveOneFloor time.Duration

	mainLoopTicker     *time.Ticker
	mainLoopTickerCh   <-chan time.Time
	mainLoopShutdownCh chan struct{}
	piDevice           common.RPi // the interface with the raspberry pi device
}

// NewController make a Controller object
func NewController(maxFloors int) *Controller {
	piDevice := common.NewRPiDevice()
	return newController(maxFloors, piDevice, defaultLoopFrequency, nil)
}

// newController private controller constructor exposing Pi device interface for testing
func newController(maxFloors int, rpi common.RPi, loopFreq time.Duration, timingControlCh <-chan time.Time) *Controller {
	var controlCh <-chan time.Time
	var ticker *time.Ticker
	if timingControlCh == nil {
		ticker = time.NewTicker(loopFreq) // production code uses ticker at loopFreq
		controlCh = ticker.C
	} else {
		controlCh = timingControlCh // tests that need to control the main loop pass in their own chan
	}

	controller := &Controller{
		topFloor:           maxFloors,
		piDevice:           rpi,
		mainLoopTickerCh:   controlCh,
		movingDirection:    Stopped,
		mainLoopTicker:     ticker,
		mainLoopShutdownCh: make(chan struct{}),
	}

	controller.init()

	go controller.startProcessingLoop()

	return controller
}

// init query all the floor levers and initialize the car's location
func (c *Controller) init() {
	// TODO implement
}

// startProcessingLoop run the processing loop that listens for signals from the floor and user
// requests and controlls sending up/down/stop commands to the garage door opener
func (c *Controller) startProcessingLoop() {
	log.Info("starting controller main loop")

	for {
		select {
		case tick := <-c.mainLoopTickerCh:
			log.Infof("got tick: %v", tick)
			// if the car is stationary and another floor is requested, start it moving in the requested direction
			// if the car is moving and a floor in the opposite direction has been requested stop the car
			// (let the next iteration start it moving)
			if c.GetRequestedFloor() > c.GetLastSeenFloor() {
				if c.GetMovingDirection() == Stopped {
					c.sendUp()
				} else if c.GetMovingDirection() == Down {
					c.stop() // stop the machine, it start moving up on next iteration
				}
				// else do nothing it is already moving up
			} else if c.GetRequestedFloor() < c.GetLastSeenFloor() {
				if c.GetMovingDirection() == Stopped {
					c.sendDown()
				} else if c.GetMovingDirection() == Up {
					c.stop() // stop the machine, it will start moving down on next iteration
				}
			} else {
				c.stop()
			}
		}
	}
}

// GetStatus get the dumbwaiter status
func (c *Controller) GetStatus() *Status {
	return &Status{
		LastSeenFloor:   c.GetLastSeenFloor(),
		MovingDirection: c.GetMovingDirection(),
		RequestedFloor:  c.GetRequestedFloor(),

		// TODO add floors' status
	}
}

func (c *Controller) sendUp() {
	c.piDevice.SendSignal(common.OpenerUp)
	c.SetMovingDirection(Up)
}

func (c *Controller) sendDown() {
	c.piDevice.SendSignal(common.OpenerDown)
	c.SetMovingDirection(Down)
}

func (c *Controller) stop() {
	c.piDevice.SendSignal(common.OpenerStop)
	c.SetMovingDirection(Stopped)
}

// GetLastSeenFloor return the floor the dumbwaiter's car was last seen at
func (c *Controller) GetLastSeenFloor() int {
	c.lastSeenFloorMU.RLock()
	defer c.lastSeenFloorMU.RUnlock()
	return c.lastSeenFloor
}

// SetLastSeenFloor set floor number the dumbwaiter's car was last seen at
func (c *Controller) SetLastSeenFloor(floor int) {
	c.lastSeenFloorMU.Lock()
	defer c.lastSeenFloorMU.Unlock()
	c.lastSeenFloor = floor
}

// GetRequestedFloor return the floor the dumbwaiter car should move to
func (c *Controller) GetRequestedFloor() int {
	c.requestedFloorMU.RLock()
	defer c.requestedFloorMU.RUnlock()
	return c.requestedFloor
}

// SetRequestedFloor set the floor the dumbwaiter car should move to
func (c *Controller) SetRequestedFloor(floor int) {
	c.requestedFloorMU.Lock()
	defer c.requestedFloorMU.Unlock()
	c.requestedFloor = floor
}

// GetMovingDirection get the dumbwaiter's current direction
func (c *Controller) GetMovingDirection() Direction {
	c.movingDirectionMu.RLock()
	defer c.movingDirectionMu.RUnlock()
	return c.movingDirection
}

// SetMovingDirection set the dumbwaiter's moving direction
func (c *Controller) SetMovingDirection(movingDirection Direction) {
	c.movingDirectionMu.Lock()
	defer c.movingDirectionMu.Unlock()
	c.movingDirection = movingDirection
}

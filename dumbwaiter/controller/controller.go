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

	movingDirection Direction // the direction the cab is currently moving

	timeStartedMoving  time.Time
	timeToMoveOneFloor time.Duration

	mainLoopFrequency time.Duration
	piDevice          common.RPi // the interface with the raspberry pi device
}

// NewController make a Controller object
func NewController(maxFloors int) *Controller {
	piDevice := common.NewRPiDevice()
	return newController(maxFloors, piDevice, defaultLoopFrequency)
}

// newController private controller constructor exposing Pi device interface for testing
func newController(maxFloors int, rpi common.RPi, loopFrequency time.Duration) *Controller {
	controller := &Controller{
		topFloor:          maxFloors,
		piDevice:          rpi,
		mainLoopFrequency: loopFrequency,
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
		// if the car is stationary and another floor is requested, start it moving in the requested direction
		// if the car is moving and a floor in the opposite direction has been requested stop the car
		// (let the next iteration start it moving)
		if c.requestedFloor > c.GetLastSeenFloor() {
			if c.movingDirection == Stopped {
				c.sendUp()
			} else if c.movingDirection == Down {
				c.stop() // stop the machine, it start moving up on next iteration
			}
			// else do nothing it is already moving up
		} else if c.requestedFloor < c.GetLastSeenFloor() {
			if c.movingDirection == Stopped {
				c.sendDown()
			} else if c.movingDirection == Up {
				c.stop() // stop the machine, it will start moving down on next iteration
			}
		} else {
			c.stop()
		}

		time.Sleep(c.mainLoopFrequency)
	}
}

// GetStatus get the dumbwaiter status
func (c *Controller) GetStatus() *Status {
	return &Status{
		LastSeenFloor:   c.lastSeenFloor,
		MovingDirection: c.movingDirection,
		RequestedFloor:  c.requestedFloor,

		// TODO add floors' status
	}
}

func (c *Controller) sendUp() {
	c.piDevice.SendSignal(common.OpenerUp)
	c.movingDirection = Up
}

func (c *Controller) sendDown() {
	c.piDevice.SendSignal(common.OpenerDown)
	c.movingDirection = Down
}

func (c *Controller) stop() {
	c.piDevice.SendSignal(common.OpenerStop)
	c.movingDirection = Stopped
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

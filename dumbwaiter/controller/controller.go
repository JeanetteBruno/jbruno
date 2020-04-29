package controller

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/JeanetteBruno/jbruno/dumbwaiter/common"
)

var defaultLoopFrequency time.Duration = 500 * time.Millisecond

// Direction the direction the dumbwaiter car is moving
type Direction int

//Direction Constants are Up, Down, and Stopped.
const (
	Up Direction = iota
	Down
	Stopped
)

//platform location wrt last seen floor
type platformLoc int

//platform location values
const (
	Above platformLoc = iota
	Below
	At
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
	platformWRTLSF   platformLoc
	platformWRTLSFMU sync.RWMutex

	topFloor int // the top floor number (floor numbers start at 1)

	movingDirection   Direction // the direction the cab is currently moving
	movingDirectionMu sync.RWMutex

	timeToMoveOneFloor time.Duration

	mainLoopTicker *time.Ticker
	mainLoopFreq   time.Duration
	piDevice       common.RPi // the interface with the raspberry pi device
}

// NewController make a Controller object
func NewController(maxFloors int) *Controller {
	piDevice := common.NewRPiDevice()
	return &Controller{
		topFloor:        maxFloors,
		piDevice:        piDevice,
		movingDirection: Stopped,
		mainLoopFreq:    defaultLoopFrequency}
}

// StartProcessingLoop start the processing loop in its own goroutine
func (c *Controller) StartProcessingLoop() {
	go c.processingLoop()
}

// processingLoop run the processing loop that listens for signals from the floor and user
// requests and controlls sending up/down/stop commands to the garage door opener
func (c *Controller) processingLoop() {
	log.Info("starting controller main loop")
	c.mainLoopTicker = time.NewTicker(c.mainLoopFreq)

	for {
		select {
		case <-c.mainLoopTicker.C:
			// if the car is stationary and another floor is requested, start it moving in the requested direction
			// if the car is moving and a floor in the opposite direction has been requested stop the car
			// if the requested floor is given a stop value; keep the dumbwaiter stopped
			// (let the next iteration start it moving)

			//if the requested floor is the stop value direction is stopped
			if c.requestedFloor == -1 {
				c.stop()
				//else move the dumbaiter based on requested floor and last seen floor
			} else {
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
					// else if requested floor equals the last seen floor, and if platform is below last seend floor  -> send up
				} else if c.GetplatformWRTLSF() == Below {
					c.sendUp()
					// else if requested floor equals the last seen floor, and if platform is above last seend floor -> send down
				} else if c.GetplatformWRTLSF() == Above {
					c.sendDown()
					// if requested floor equals the last seen floor, set moving direction to stop
				} else {
					c.stop()
				}
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
	log.Info("controller sending up")
	c.piDevice.SendSignal(common.OpenerUp)
	c.SetMovingDirection(Up)
}

func (c *Controller) sendDown() {
	log.Info("controller sending down")
	c.piDevice.SendSignal(common.OpenerDown)
	c.SetMovingDirection(Down)
}

func (c *Controller) stop() {
	log.Info("controller stopping")
	c.piDevice.SendSignal(common.OpenerStop)
	c.SetMovingDirection(Stopped)
	c.requestedFloor = -1
}

// GetLastSeenFloor return the floor the dumbwaiter's car was last seen at
func (c *Controller) GetLastSeenFloor() int {
	c.lastSeenFloorMU.RLock()
	defer c.lastSeenFloorMU.RUnlock()
	return c.lastSeenFloor
}

// SetFloorStatus tells us when platform has reached or left a floor
func (c *Controller) SetFloorStatus(floor int, atFloor bool) {
	log.Infof("controller setting last seen floor to %d", floor)
	c.lastSeenFloorMU.Lock()
	defer c.lastSeenFloorMU.Unlock()
	if atFloor {
		c.lastSeenFloor = floor
		c.SetPlatformWRTLSF(At)
	} else {
		if c.movingDirection == Up {
			c.SetPlatformWRTLSF(Above)
		} else if c.movingDirection == Down {
			c.SetPlatformWRTLSF(Below)
		}

	}
}

// GetplatformWRTLSF return the platform location with respect to the last seen floor
func (c *Controller) GetplatformWRTLSF() platformLoc {
	c.platformWRTLSFMU.RLock()
	defer c.platformWRTLSFMU.RUnlock()
	return c.platformWRTLSF
}

// GetRequestedFloor return the floor the dumbwaiter car should move to
func (c *Controller) GetRequestedFloor() int {
	c.requestedFloorMU.RLock()
	defer c.requestedFloorMU.RUnlock()
	return c.requestedFloor
}

// SetRequestedFloor set the floor the dumbwaiter car should move to
func (c *Controller) SetRequestedFloor(floor int) {
	log.Infof("controller setting requested floor to %d", floor)
	c.requestedFloorMU.Lock()
	defer c.requestedFloorMU.Unlock()
	c.requestedFloor = floor
}

//SetStopRequested get a stop request from a floor sensor
func (c *Controller) SetStopRequested() {
	log.Infof("controller recieved a stop request")
	c.requestedFloorMU.Lock()
	defer c.requestedFloorMU.Unlock()
	//when requested floor equals the last seen floor the controller will send a stop request
	c.requestedFloor = c.lastSeenFloor
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

// SetPlatformWRTLSF set the platform location with respect to the last seen floor
func (c *Controller) SetPlatformWRTLSF(platformLoc platformLoc) {
	log.Infof("controller setting platform location to %d", platformLoc)
	c.platformWRTLSFMU.Lock()
	defer c.platformWRTLSFMU.Unlock()
	c.platformWRTLSF = platformLoc
}

// Controller constructor setters for builder pattern

// SetRPiDevice used by testing to override production RPi interface
func (c *Controller) SetRPiDevice(piDevice common.RPi) *Controller {
	c.piDevice = piDevice
	return c
}

// SetLoopFrequency used by testing to speed up tests
func (c *Controller) SetLoopFrequency(freq time.Duration) *Controller {
	c.mainLoopFreq = freq
	return c
}

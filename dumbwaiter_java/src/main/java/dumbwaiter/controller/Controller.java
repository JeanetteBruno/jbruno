package dumbwaiter.controller;

import dumbwaiter.common.*;

public class Controller {

  private int lastSeenFloor;
  private int requestedFloor;
  private Direction movingDirection;

  private int topFloor;
  private RPi piDevice;

  private TimingControl sleeper;

  public Controller(int maxFloors) {
    this.topFloor = maxFloors;
    this.piDevice = new RPiDevice();  // NOTE: instantiating RPiDevice() may become a problem for testing
    this.sleeper = new Sleeper();
  }

  /*
  start the processing loop (in its own thread), this is a non-blocking call
   */
  public void startProcessing() {
    Runnable processingLoop = () -> {
      this.processingLoop();
    };
    Thread thread = new Thread(processingLoop);
    thread.start();
  }

  public void processingLoop() {
    /*
     if the car is stationary and another floor is requested, start it moving in the requested direction
     if the car is moving and a floor in the opposite direction has been requested stop the car
     (let the next iteration start it moving)
    */
    while (true) {
      try {
        if (this.getRequestedFloor() > this.getLastSeenFloor()) {
          if (this.getMovingDirection() == Direction.Stopped) {
            this.sendUp();  // start the platform moving up
          } else if (this.getMovingDirection() == Direction.Down) {
            this.sendStop(); // stop the machine, it will start moving up on next iteration
          } // else do nothing it is already moving up
        } else if (this.getRequestedFloor() < this.getLastSeenFloor()) {
          if (this.getMovingDirection() == Direction.Stopped) {
            this.sendDown(); // start the platform moving down
          } else if (this.getMovingDirection() == Direction.Up) {
            this.sendStop(); // stop the machine, it will start moving down on next iteration
          } // else do nothing it is already moving down
        } else if (this.getMovingDirection() != Direction.Stopped) {
          this.sendStop();  // we're at the requested floor, stop the platform if it is moving
        }
      } catch (RPiDeviceException e) {
        System.out.println(String.format("exception %e", e));
      }

      this.getSleeper().sleep(Defaults.defaultLoopFrequencyMs);
    }
  }

  private void sendUp() throws RPiDeviceException {
    this.getPiDevice().SendSignal(PiPin.WindingUp);
    this.setMovingDirection(Direction.Up);
  }

  private void sendDown() throws RPiDeviceException {
    this.getPiDevice().SendSignal(PiPin.WindingDown);
    this.setMovingDirection(Direction.Down);
  }

  private void sendStop() throws RPiDeviceException {
    this.getPiDevice().SendSignal(PiPin.WindingStop);
    this.setMovingDirection(Direction.Stopped);
  }

  public int getTopFloor() {
    return topFloor;
  }

  public Controller setTopFloor(int topFloor) {
    this.topFloor = topFloor;
    return this;
  }

  public RPi getPiDevice() {
    return piDevice;
  }

  public Controller setPiDevice(RPi piDevice) {
    this.piDevice = piDevice;
    return this;
  }

  public synchronized Direction getMovingDirection() {
    return movingDirection;
  }

  public synchronized Controller setMovingDirection(Direction movingDirection) {
    this.movingDirection = movingDirection;
    return this;
  }

  public synchronized int getRequestedFloor() {
    return requestedFloor;
  }

  public synchronized Controller setRequestedFloor(int requestedFloor) {
    this.requestedFloor = requestedFloor;
    return this;
  }

  public synchronized int getLastSeenFloor() {
    return lastSeenFloor;
  }

  public synchronized Controller setLastSeenFloor(int floor) {
    this.lastSeenFloor = floor;
    return this;
  }

  public synchronized TimingControl getSleeper() {
    return sleeper;
  }

  public synchronized Controller setSleeper(TimingControl sleeper) {
    this.sleeper = sleeper;
    return this;
  }
}


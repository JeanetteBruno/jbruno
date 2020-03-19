package dumbwaiter;

import java.time.LocalDateTime;
import java.util.Arrays;

import org.junit.Assert;
import org.junit.Test;

import dumbwaiter.common.*;
import dumbwaiter.controller.*;

public class TestController {

  @Test
  public void testRequestUpFromStop() {
    PiPin[] expectedPins = new PiPin[]{PiPin.WindingUp};
    FakeRPiDevice fakePi = new FakeRPiDevice(expectedPins);
    Controller dwController = new Controller(3).
        setPiDevice(fakePi).
        setLastSeenFloor(2).
        setMovingDirection(Direction.Stopped);
    dwController.startProcessing();  // start the device processing

    // test
    dwController.setRequestedFloor(3);
    this.waitForStatus(dwController, 2, Direction.Up, 3, 6L);
  }

  @Test
  public void testRequestDownFromStop() {
    PiPin[] expectedPins = new PiPin[]{PiPin.WindingDown};
    FakeRPiDevice fakePi = new FakeRPiDevice(expectedPins);
    Controller dwController = new Controller(3).
        setPiDevice(fakePi).
        setLastSeenFloor(2).
        setMovingDirection(Direction.Stopped);
    dwController.startProcessing();  // start the device processing

    // test
    dwController.setRequestedFloor(1);
    this.waitForStatus(dwController, 2, Direction.Down, 1, 6L);
  }

  @Test
  public void testRequestDownFromUp() {
    PiPin[] expectedPins = new PiPin[]{PiPin.WindingStop, PiPin.WindingDown};
    FakeRPiDevice fakePi = new FakeRPiDevice(expectedPins);
    FakeTicker fakeTicker = new FakeTicker();
    Controller dwController = new Controller(3).
        setPiDevice(fakePi).
        setLastSeenFloor(2).
        setMovingDirection(Direction.Up).
        setSleeper(fakeTicker);
    dwController.startProcessing();  // start the device processing

    // test
    dwController.setRequestedFloor(1);
    fakeTicker.tick();  // signal a timing 'tick' - causing controller thread to execute one processing loop
    this.waitForStatus(dwController, 2, Direction.Stopped, 1, 1L);

    fakeTicker.tick();  // signal the next 'tick' -causing controller thread
    this.waitForStatus(dwController, 2, Direction.Down, 1, 6L);
  }

  @Test
  public void testRequestUpFromDown() {
    PiPin[] expectedPins = new PiPin[]{PiPin.WindingStop, PiPin.WindingUp};
    FakeRPiDevice fakePi = new FakeRPiDevice(expectedPins);
    Controller dwController = new Controller(3).
        setPiDevice(fakePi).
        setLastSeenFloor(2).
        setMovingDirection(Direction.Down);
    dwController.startProcessing();  // start the device processing

    // test
    dwController.setRequestedFloor(3);
    this.waitForStatus(dwController, 2, Direction.Up, 3, 6L);
  }

  private void waitForStatus(Controller dwc, int lastSeenFloor, Direction direction, int requestedFloor, long timeoutSec) {
    LocalDateTime end = LocalDateTime.now().plusSeconds(timeoutSec);
    while (LocalDateTime.now().isBefore(end)) {
      if (lastSeenFloor == dwc.getLastSeenFloor() &&
          direction == dwc.getMovingDirection() &&
          requestedFloor == dwc.getRequestedFloor()) {
            return;
      }
      try {
        Thread.sleep(timeoutSec * 1000);
      } catch (Exception e) {
        Assert.fail(String.format("Sleep errored: %s", e));
      }
    }

    if (lastSeenFloor == dwc.getLastSeenFloor() &&
        direction == dwc.getMovingDirection() &&
        requestedFloor == dwc.getRequestedFloor()) {
          return;
    }
    Assert.fail(String.format("expected:lastSeenFloor:%d, direction:%s, requestedFloor:%d\n" +
            "got: lastSeenFloor:%d, direction: %s, requestedFloor: %d",
        lastSeenFloor, direction, requestedFloor,
        dwc.getLastSeenFloor(), dwc.getMovingDirection(), dwc.getRequestedFloor()));
  }

  /*
  pause the testing thread till the controller thread has consumed the tick
   */
  private void waitForTickConsumed(FakeTicker fakeTicker) {
    while (true) {
      try {
        Thread.sleep(100);
      } catch (InterruptedException e) {
        e.printStackTrace();
      }
      if (! fakeTicker.isTick()) {
        return;
      }
    }
  }
  /*
  implements the RPi interface with behavior that validates the expected pins are received
   */
  public class FakeRPiDevice implements RPi {
    private PiPin[] expectedSends;
    private int callCount;

    public FakeRPiDevice(PiPin[] expectedSends) {
      this.expectedSends = expectedSends;
      this.callCount = 0;
    }

    @Override
    public void SendSignal(PiPin pin) {
      if (this.callCount >= expectedSends.length) {
        Assert.fail(String.format("expected only %d calls, have %d", this.expectedSends.length, this.callCount + 1));
      }
      if (expectedSends[this.callCount] != pin) {
        Assert.fail(String.format("expected pin: %s, have pin: %s", expectedSends[this.callCount], pin));
      }
      this.callCount++;
    }

    @Override
    public boolean GetSignal(PiPin pin) throws RPiDeviceException {
      return false;
    }
  }
}

